package persistentnodeusecase

import (
	"context"
	"cynxhostagent/internal/model/request"
	"cynxhostagent/internal/model/response"
	"cynxhostagent/internal/model/response/responsecode"
	"cynxhostagent/internal/model/response/responsedata"
	"fmt"
	"math/rand"
	"strconv"
)

func (uc *PersistentNodeUseCaseImpl) StreamLogs(ctx context.Context, req request.GetPersistentNodeRealTimeLogsRequest, channel chan string) error {
	// Start streaming logs from a specific container
	go func() {
		// Start the interactive session
		err := uc.dockerManager.StreamOutput(req.SessionId, channel)
		if err != nil {
			fmt.Println("Error starting interactive session:", err)
			channel <- fmt.Sprintf("Error starting interactive session: %v", err)
			return
		}

		// Stream the logs and send them to the channel
		// Assuming the `StartInteractiveSession` is configured to stream logs to a WebSocket channel
		for logLine := range channel {
			// Sending logs to the channel for WebSocket transmission
			channel <- logLine
		}
	}()

	return nil
}

func (uc *PersistentNodeUseCaseImpl) CreateSession(ctx context.Context, req request.StartSessionRequest, resp *response.APIResponse) {
	// Create a new interactive session
	uniqueCode := strconv.Itoa(rand.Int())
	fmt.Println("Creating new session with code ", uniqueCode)
	err := uc.dockerManager.CreateNewSession(uniqueCode, uc.config.DockerConfig.ContainerName, req.Shell)
	if err != nil {
		resp.Code = responsecode.CodeDockerError
		resp.Error = fmt.Sprintf("Error creating new session: %v", err)
		return
	}

	resp.Code = responsecode.CodeSuccess
	resp.Data = responsedata.CreateSessionResponseData{
		SessionId: uniqueCode,
	}
}

func (uc *PersistentNodeUseCaseImpl) SendCommand(ctx context.Context, req request.SendCommandRequest, resp *response.APIResponse) {

	// Send the command to the Docker container
	err := uc.dockerManager.SendCommand(req.SessionId, req.Command, req.IsBase64Encoded)
	if err != nil {
		resp.Code = responsecode.CodeDockerError
		resp.Error = fmt.Sprintf("Error sending command to Docker container: %v", err)
		return
	}

	resp.Code = responsecode.CodeSuccess
}
