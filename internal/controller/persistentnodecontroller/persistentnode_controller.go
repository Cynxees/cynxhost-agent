package persistentnodecontroller

import (
	"context"
	"cynxhostagent/internal/dependencies"
	"cynxhostagent/internal/helper"
	"cynxhostagent/internal/model/request"
	"cynxhostagent/internal/model/response"
	"cynxhostagent/internal/model/response/responsecode"
	"cynxhostagent/internal/usecase"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

type PersistentNodeController struct {
	persistentNodeUsecase usecase.PersistentNodeUseCase
	validator             *validator.Validate
	config                *dependencies.Config
}

func New(
	persistentNodeUseCase usecase.PersistentNodeUseCase,
	validate *validator.Validate,
	config *dependencies.Config,
) *PersistentNodeController {
	return &PersistentNodeController{
		persistentNodeUsecase: persistentNodeUseCase,
		validator:             validate,
		config:                config,
	}
}

func (controller *PersistentNodeController) RunPersistentNodeTemplateScript(w http.ResponseWriter, r *http.Request) (context.Context, response.APIResponse) {
	var requestBody request.RunPersistentNodeTemplateScriptRequest
	var apiResponse response.APIResponse

	ctx := r.Context()

	sessionUser, ok := helper.GetUserFromContext(ctx)
	if !ok {
		apiResponse.Code = responsecode.CodeAuthenticationError
		apiResponse.Error = "User not found in context"
		return ctx, apiResponse
	}

	requestBody.PersistentNodeId = *controller.config.App.PersistentNodeId

	requestBody.SessionUser = sessionUser
	if err := helper.DecodeAndValidateRequest(r, &requestBody, controller.validator); err != nil {
		apiResponse.Code = responsecode.CodeValidationError
		apiResponse.Error = err.Error()
		return ctx, apiResponse
	}

	controller.persistentNodeUsecase.RunPersistentNodeTemplateScript(ctx, requestBody, &apiResponse)

	return ctx, apiResponse
}

func (controller *PersistentNodeController) GetPersistentNodeRealTimeLogs(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {
	// Channel to receive logs from Docker container
	logChannel := make(chan string)

	// Close the channel when the function returns
	defer close(logChannel)
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	id := r.URL.Query().Get("id")
	if id == "" {
		log.Println("ID parameter missing from the WebSocket URL")
		return
	}

	request := request.GetPersistentNodeRealTimeLogsRequest{
		SessionId: id,
		// TODO: Fill the fields
	}

	// Start streaming logs from the Docker container
	go func() {
		if err := controller.persistentNodeUsecase.StreamLogs(context.Background(), request, logChannel); err != nil {
			log.Printf("Error streaming Docker logs: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error: Unable to stream Docker logs"))
			close(logChannel) // Close the channel if there's an error
			return
		}
	}()

	// Process and send logs over WebSocket as they are received
	for logLine := range logChannel {
		err := conn.WriteMessage(websocket.TextMessage, []byte(logLine))
		if err != nil {
			log.Printf("Failed to send message: %v", err)
			break // Stop processing if there's an error sending the message
		}
		log.Println("Log sent:", logLine)
	}
}

func (controller *PersistentNodeController) SendCommand(w http.ResponseWriter, r *http.Request) (context.Context, response.APIResponse) {
	var requestBody request.SendCommandRequest
	var apiResponse response.APIResponse

	ctx := r.Context()

	if err := helper.DecodeAndValidateRequest(r, &requestBody, controller.validator); err != nil {
		apiResponse.Code = responsecode.CodeValidationError
		apiResponse.Error = err.Error()
		return ctx, apiResponse
	}

	controller.persistentNodeUsecase.SendCommand(ctx, requestBody, &apiResponse)

	return ctx, apiResponse
}


func (controller *PersistentNodeController) SendSingleDockerCommand(w http.ResponseWriter, r *http.Request) (context.Context, response.APIResponse) {
	var requestBody request.SendSingleDockerCommandRequest
	var apiResponse response.APIResponse

	ctx := r.Context()

	if err := helper.DecodeAndValidateRequest(r, &requestBody, controller.validator); err != nil {
		apiResponse.Code = responsecode.CodeValidationError
		apiResponse.Error = err.Error()
		return ctx, apiResponse
	}

	controller.persistentNodeUsecase.SendSingleDockerCommand(ctx, requestBody, &apiResponse)

	return ctx, apiResponse
}

func (controller *PersistentNodeController) CreateSession(w http.ResponseWriter, r *http.Request) (context.Context, response.APIResponse) {
	var apiResponse response.APIResponse
	var requestBody request.StartSessionRequest

	ctx := r.Context()

	if err := helper.DecodeAndValidateRequest(r, &requestBody, controller.validator); err != nil {
		apiResponse.Code = responsecode.CodeValidationError
		apiResponse.Error = err.Error()
		return ctx, apiResponse
	}

	controller.persistentNodeUsecase.CreateSession(ctx, requestBody, &apiResponse)

	return ctx, apiResponse
}

func (controller *PersistentNodeController) GetServerProperties(w http.ResponseWriter, r *http.Request) (context.Context, response.APIResponse) {
	var apiResponse response.APIResponse

	ctx := r.Context()

	controller.persistentNodeUsecase.GetServerProperties(ctx, &apiResponse)

	return ctx, apiResponse
}

func (controller *PersistentNodeController) SetServerProperties(w http.ResponseWriter, r *http.Request) (context.Context, response.APIResponse) {
	var requestBody request.SetServerPropertiesRequest
	var apiResponse response.APIResponse

	ctx := r.Context()

	if err := helper.DecodeAndValidateRequest(r, &requestBody, controller.validator); err != nil {
		apiResponse.Code = responsecode.CodeValidationError
		apiResponse.Error = err.Error()
		return ctx, apiResponse
	}

	controller.persistentNodeUsecase.SetServerProperties(ctx, requestBody, &apiResponse)

	return ctx, apiResponse
}

func (controller *PersistentNodeController) GetNodeContainerStats(w http.ResponseWriter, r *http.Request) (context.Context, response.APIResponse) {
	var apiResponse response.APIResponse

	ctx := r.Context()

	controller.persistentNodeUsecase.GetNodeContainerStats(ctx, &apiResponse)

	return ctx, apiResponse
}
