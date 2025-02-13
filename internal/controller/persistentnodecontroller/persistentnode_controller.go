package persistentnodecontroller

import (
	"context"
	"cynxhostagent/internal/dependencies"
	"cynxhostagent/internal/helper"
	"cynxhostagent/internal/model/request"
	"cynxhostagent/internal/model/response"
	"cynxhostagent/internal/model/response/responsecode"
	"cynxhostagent/internal/usecase"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
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

	requestBody.PersistentNodeId = controller.config.App.PersistentNode.Id

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

// func (controller *PersistentNodeController) GetServerProperties(w http.ResponseWriter, r *http.Request) (context.Context, response.APIResponse) {
// 	var apiResponse response.APIResponse

// 	ctx := r.Context()

// 	controller.persistentNodeUsecase.GetServerProperties(ctx, &apiResponse)

// 	return ctx, apiResponse
// }

// func (controller *PersistentNodeController) SetServerProperties(w http.ResponseWriter, r *http.Request) (context.Context, response.APIResponse) {
// 	var requestBody request.SetServerPropertiesRequest
// 	var apiResponse response.APIResponse

// 	ctx := r.Context()

// 	if err := helper.DecodeAndValidateRequest(r, &requestBody, controller.validator); err != nil {
// 		apiResponse.Code = responsecode.CodeValidationError
// 		apiResponse.Error = err.Error()
// 		return ctx, apiResponse
// 	}

// 	controller.persistentNodeUsecase.SetServerProperties(ctx, requestBody, &apiResponse)

// 	return ctx, apiResponse
// }

func (controller *PersistentNodeController) GetNodeContainerStats(w http.ResponseWriter, r *http.Request) (context.Context, response.APIResponse) {
	var apiResponse response.APIResponse

	ctx := r.Context()

	controller.persistentNodeUsecase.GetNodeContainerStats(ctx, &apiResponse)

	return ctx, apiResponse
}

func (controller *PersistentNodeController) DownloadFile(w http.ResponseWriter, r *http.Request) {
	var requestBody request.DownloadFileRequest
	var apiResponse response.APIResponse

	// Decode and validate request
	if err := helper.DecodeAndValidateRequest(r, &requestBody, controller.validator); err != nil {
		apiResponse.Code = responsecode.CodeValidationError
		apiResponse.Error = err.Error()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the file content
	fileContent, err := controller.persistentNodeUsecase.DownloadFile(r.Context(), requestBody, &apiResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error downloading file: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Set headers for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", path.Base(requestBody.FilePath)))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileContent)))

	// Write file content to response
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(fileContent) // Ignore error since nothing can be done if writing fails
}

func (controller *PersistentNodeController) UploadFile(w http.ResponseWriter, r *http.Request) (context.Context, response.APIResponse) {
	var apiResponse response.APIResponse
	var requestBody request.UploadFileRequest

	ctx := r.Context()

	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
		apiResponse.Code = responsecode.CodeValidationError
		apiResponse.Error = "Failed to parse multipart form: " + err.Error()
		return ctx, apiResponse
	}

	// Retrieve file from form
	file, fileHeader, err := r.FormFile("file") // "file" should match the form field name
	if err != nil {
		apiResponse.Code = responsecode.CodeValidationError
		apiResponse.Error = "Failed to get file from request: " + err.Error()
		return ctx, apiResponse
	}
	defer file.Close()

	// Assign file details to request body
	requestBody.FileData = file
	requestBody.FileHeader = *fileHeader

	if err := helper.DecodeAndValidateRequest(r, &requestBody, controller.validator); err != nil {
		apiResponse.Code = responsecode.CodeValidationError
		apiResponse.Error = err.Error()
		return ctx, apiResponse
	}

	// Check if filename contains a path
	if strings.Contains(requestBody.FileName, "/") {
		apiResponse.Code = responsecode.CodeValidationError
		apiResponse.Error = "Invalid filename: " + requestBody.FileName
		return ctx, apiResponse
	}

	controller.persistentNodeUsecase.UploadFile(ctx, requestBody, &apiResponse)

	return ctx, apiResponse
}

func (controller *PersistentNodeController) RemoveFile(w http.ResponseWriter, r *http.Request) (context.Context, response.APIResponse) {
	var requestBody request.RemoveFileRequest
	var apiResponse response.APIResponse

	ctx := r.Context()

	if err := helper.DecodeAndValidateRequest(r, &requestBody, controller.validator); err != nil {
		apiResponse.Code = responsecode.CodeValidationError
		apiResponse.Error = err.Error()
		return ctx, apiResponse
	}

	controller.persistentNodeUsecase.RemoveFile(ctx, requestBody, &apiResponse)

	return ctx, apiResponse
}

func (controller *PersistentNodeController) ListDirectory(w http.ResponseWriter, r *http.Request) (context.Context, response.APIResponse) {
	var requestBody request.ListDirectoryRequest
	var apiResponse response.APIResponse

	ctx := r.Context()

	if err := helper.DecodeAndValidateRequest(r, &requestBody, controller.validator); err != nil {
		apiResponse.Code = responsecode.CodeValidationError
		apiResponse.Error = err.Error()
		return ctx, apiResponse
	}

	controller.persistentNodeUsecase.ListDirectory(ctx, requestBody, &apiResponse)

	return ctx, apiResponse
}
