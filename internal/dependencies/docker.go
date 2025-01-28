package dependencies

import (
	"encoding/base64"
	"fmt"
	"io"
	"sync"

	"github.com/docker/docker/client"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

var (
	mu          sync.Mutex // Mutex to ensure thread-safe access to the session map
	sshSessions map[string]*PersistentSession
)

type DockerManager struct {
	client *client.Client
}

type PersistentSession struct {
	client  *goph.Client
	session *ssh.Session
	stdin   io.WriteCloser
	stdout  io.Reader
}

func NewDockerManager() *DockerManager {
	sshSessions = make(map[string]*PersistentSession)
	return &DockerManager{}
}

// CreateNewSession creates a persistent SSH session.
func (m *DockerManager) CreateNewSession(sessionId string, host string, port uint, username string, password string) error {
	auth := goph.Password(password)

	fmt.Println("Creating new connection")
	client, err := goph.NewUnknown(username, host, port, auth)
	if err != nil {
		return fmt.Errorf("Failed to create ssh client: %v", err)
	}

	fmt.Println("Creating new session")
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("Failed to create ssh session: %v", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("Failed to get stdin pipe: %v", err)
	}

	// Set the output to os.Stdout
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Failed to set stdout: %v", err)
	}

	// Start an interactive shell
	err = session.Shell()
	if err != nil {
		return fmt.Errorf("Failed to start shell: %v", err)
	}

	fmt.Printf("Saving client for session: %s\n", sessionId)
	mu.Lock()
	sshSessions[sessionId] = &PersistentSession{
		client:  client,
		session: session,
		stdin:   stdin,
		stdout:  stdout,
	}
	mu.Unlock()

	fmt.Println("Session created successfully")
	return nil
}

// SendCommand sends a command to the persistent session.
func (m *DockerManager) SendCommand(sessionId string, command string, isBase64Encoded bool) error {
	mu.Lock()
	pSession, ok := sshSessions[sessionId]
	if !ok {
		mu.Unlock()
		return fmt.Errorf("Session not found")
	}
	mu.Unlock()

	commandByte := []byte(command + "\n")

	if isBase64Encoded {
		res, err := base64.StdEncoding.DecodeString(command)
		if err != nil {
			return fmt.Errorf("Failed to decode command: %v", err)
		}
		commandByte = res
	}

	// Write the command to the shell
	_, err := pSession.stdin.Write(commandByte)
	if err != nil {
		return fmt.Errorf("Failed to send command: %v", err)
	}

	return nil
}

// CloseSession closes the SSH session and cleans up resources.
func (m *DockerManager) CloseSession(sessionId string) error {
	mu.Lock()
	pSession, ok := sshSessions[sessionId]
	if ok {
		pSession.session.Close()
		pSession.client.Close()
		delete(sshSessions, sessionId)
	}
	mu.Unlock()

	fmt.Printf("Session %s closed\n", sessionId)
	return nil
}

func (m *DockerManager) StreamOutput(sessionId string, outChan chan string) error {
	mu.Lock()
	pSession, ok := sshSessions[sessionId]
	if !ok {
		mu.Unlock()
		return fmt.Errorf("Session not found")
	}
	mu.Unlock()

	buf := make([]byte, 1024) // Buffer to read output
	for {
		n, err := pSession.stdout.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("Error reading stdout: %v", err)
		}
		if n > 0 {
			outChan <- string(buf[:n]) // Send output to channel
		}
		if err == io.EOF {
			break // End of output stream
		}
	}

	return nil
}
