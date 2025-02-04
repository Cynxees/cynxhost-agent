package dependencies

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var (
	mu          sync.Mutex // Mutex to ensure thread-safe access to the session map
	sshSessions map[string]*PersistentSession
)

type DockerManager struct {
	client *client.Client
}

type PersistentSession struct {
	pty    *os.File // The PTY file descriptor
	stdin  io.WriteCloser
	stdout io.Reader
}

func NewDockerManager() *DockerManager {
	sshSessions = make(map[string]*PersistentSession)

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	return &DockerManager{
		client: cli,
	}
}

// CreateNewSession creates a persistent SSH session with PTY to a Docker container.
func (m *DockerManager) CreateNewSession(sessionId string, containerName string, shell string) error {
	// Construct SSH command to access the Docker container
	cmd := exec.Command("docker", "exec", "-it", containerName, shell)
	ptyFile, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("Failed to create PTY: %v", err)
	}

	// Get stdin and stdout from the PTY
	stdin := ptyFile
	stdout := ptyFile

	// Save the session and PTY
	fmt.Printf("Saving client for session: %s\n", sessionId)
	mu.Lock()
	sshSessions[sessionId] = &PersistentSession{
		pty:    ptyFile,
		stdin:  stdin,
		stdout: stdout,
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

	var commandByte []byte
	if isBase64Encoded {
		res, err := base64.StdEncoding.DecodeString(command)
		if err != nil {
			return fmt.Errorf("Failed to decode base64 command: %v", err)
		}
		commandByte = res
	} else {
		commandByte = []byte(command)
	}

	// Write the command to the PTY
	fmt.Printf("Sending command to session %s: %s\n", sessionId, command)
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
		pSession.pty.Close() // Close the PTY
		delete(sshSessions, sessionId)
	}
	mu.Unlock()

	fmt.Printf("Session %s closed\n", sessionId)
	return nil
}

// StreamOutput streams output from the PTY to the provided channel.
func (m *DockerManager) StreamOutput(sessionId string, outChan chan string) error {
	mu.Lock()
	pSession, ok := sshSessions[sessionId]
	if !ok {
		mu.Unlock()
		return fmt.Errorf("Session not found")
	}
	mu.Unlock()

	buf := make([]byte, 1024) // Buffer to read output from PTY
	for {
		n, err := pSession.stdout.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("Error reading stdout: %v", err)
		}
		if n > 0 {
			encoded := base64.StdEncoding.EncodeToString(buf[:n])
			outChan <- encoded // Send Base64-encoded output to channel
		}
		if err == io.EOF {
			break // End of output stream
		}
	}

	return nil
}

func (m *DockerManager) GetContainerStats(containerNameOrId string) (*container.StatsResponse, error) {

	stats, err := m.client.ContainerStats(context.Background(), containerNameOrId, false)
	if err != nil {
		return nil, err
	}

	var data container.StatsResponse
	if err := json.NewDecoder(stats.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (m *DockerManager) SendSingleDockerCommand(containerNameOrId string, command string) (string, error) {

	cmd := exec.Command("docker", "exec", containerNameOrId, command)

	// Capture the output
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	// Run the command
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute docker command: %w", err)
	}

	return output.String(), nil
}
