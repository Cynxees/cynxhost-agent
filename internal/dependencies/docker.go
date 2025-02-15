package dependencies

import (
	"bytes"
	"context"
	"cynxhostagent/internal/model/response/responsedata"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
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

func (*DockerManager) WriteFile(filePath string, file multipart.File, header multipart.FileHeader, containerName string, fileName string) error {
	// Create a temporary file path with the original file name in the OS temporary directory
	tempDir := os.TempDir()
	tmpFilePath := path.Join(tempDir, fileName)

	// Create (or overwrite) the temporary file with the desired name
	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer tmpFile.Close()

	// Copy the file content from the uploaded file to the temporary file
	_, err = io.Copy(tmpFile, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	// Use the docker cp command to copy the file from the host to the container.
	// Make sure filePath in the container includes the desired file name.
	cmd := exec.Command("docker", "cp", tmpFilePath, fmt.Sprintf("%s:%s", containerName, filePath))
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to copy file to container: %v", err)
	}

	// Optionally, remove the temporary file after copying it into the container.
	_ = os.Remove(tmpFilePath)

	return nil
}

func (*DockerManager) GetFile(containerName, containerFilePath string) ([]byte, error) {
	// Create a temporary file to store the copied file from container
	tmpFile, err := os.CreateTemp("", "getfile-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up after function exits

	// Copy the file from container to temp file
	cmd := exec.Command("docker", "cp", fmt.Sprintf("%s:%s", containerName, containerFilePath), tmpFile.Name())
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to copy file from container: %v", err)
	}

	// Read the file content
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	return content, nil
}

func (*DockerManager) RemoveFile(containerName, containerFilePath string) error {
	// Use the docker exec command to remove the file from the container
	cmd := exec.Command("docker", "exec", containerName, "rm", containerFilePath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to remove file from container: %v", err)
	}

	return nil
}

func (dm *DockerManager) ListDirectory(containerName, containerDirPath string) ([]responsedata.File, error) {
	// List file names in the directory
	lsCmd := exec.Command("docker", "exec", containerName, "ls", containerDirPath)
	output, err := lsCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list directory: %v", err)
	}

	// Split the output into file names
	fileNames := bytes.Fields(output)
	if len(fileNames) == 0 {
		return []responsedata.File{}, nil // Return empty slice if directory is empty
	}

	var details []responsedata.File
	for _, f := range fileNames {
		filename := string(f)
		// Build full path of the file in the container
		fullPath := path.Join(containerDirPath, filename)

		// Use stat to get file details.
		// The format here is: filename|creation_time|modification_time|size
		// Note: %w returns the birth time if available, or a hyphen ("-")
		statCmd := exec.Command("docker", "exec", containerName, "stat", "--format=%n|%w|%y|%s", fullPath)
		statOutput, err := statCmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to stat file %s: %v", fullPath, err)
		}

		// Parse the output; expected format: <name>|<createdAt>|<updatedAt>|<size>
		parts := strings.Split(strings.TrimSpace(string(statOutput)), "|")
		if len(parts) != 4 {
			return nil, fmt.Errorf("unexpected stat output for %s: %s", fullPath, string(statOutput))
		}

		size, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse size for %s: %v", fullPath, err)
		}

		detail := responsedata.File{
			Filename:  parts[0],
			CreatedAt: parts[1],
			UpdatedAt: parts[2],
			Size:      size,
		}
		details = append(details, detail)
	}

	return details, nil
}

func (*DockerManager) UploadImageToAwsEcr(containerName, imageName, tag string, ecrConfig EcrConfig) error {

	ecrImage := fmt.Sprintf("%s:%s", ecrConfig.Registry, tag)

	// Save the image to a tar file
	fmt.Println("Saving image to tar file...")
	cmd := exec.Command("docker", "save", "-o", imageName+".tar", imageName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to save image to tar file: %v", err)
	}

	// Load the image back from the tar file
	fmt.Println("Loading image from tar file...")
	cmd = exec.Command("docker", "load", "-i", imageName+".tar")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to load image from tar file: %v", err)
	}

	// Tag the image for ECR
	fmt.Println("Tagging image for ECR...")
	cmd = exec.Command("docker", "tag", imageName, ecrImage)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to tag image: %v", err)
	}

	// Get ECR login password
	fmt.Println("Getting ECR login password...")
	cmd = exec.Command("aws", "ecr", "get-login-password", "--region", ecrConfig.Region)
	password, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get ECR login password: %v", err)
	}

	// Login to ECR
	fmt.Println("Logging into ECR...")
	cmd = exec.Command("docker", "login", "--username", ecrConfig.Username, "--password-stdin", ecrConfig.Registry)
	cmd.Stdin = strings.NewReader(string(password))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to login to ECR: %v", err)
	}

	// Push the image to ECR
	fmt.Println("Pushing image to ECR...")
	cmd = exec.Command("docker", "push", ecrImage)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push image to ECR: %v", err)
	}

	fmt.Println("Image successfully pushed to ECR:", ecrImage)
	return nil
}
