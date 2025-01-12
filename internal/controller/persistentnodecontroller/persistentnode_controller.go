package persistentnodecontroller

import (
	"context"
	"cynxhostagent/internal/dependencies"
	"cynxhostagent/internal/helper"
	"cynxhostagent/internal/model/request"
	"cynxhostagent/internal/model/response"
	"cynxhostagent/internal/model/response/responsecode"
	"cynxhostagent/internal/usecase"
	"io"
	"log"
	"net/http"
	"os"
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

	requestBody.SessionUser = sessionUser
	if err := helper.DecodeAndValidateRequest(r, &requestBody, controller.validator); err != nil {
		apiResponse.Code = responsecode.CodeValidationError
		apiResponse.Error = err.Error()
		return ctx, apiResponse
	}

	controller.persistentNodeUsecase.RunPersistentNodeTemplateScript(ctx, requestBody, &apiResponse)

	return ctx, apiResponse
}

func (controller *PersistentNodeController) GetPersistentNodeRealTimeLogs(conn *websocket.Conn) {

	// Path to the log file (update with your file's path)
	logFilePath := controller.config.Log.MinecraftLogFilePath

	// Open the log file for reading using os.Open, but treat it as an io.Reader
	file, err := os.Open(logFilePath)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Unable to open log file"))
		return
	}
	defer file.Close()

	// Seek to the end of the file
	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		log.Printf("Failed to seek to end of file: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Unable to seek log file"))
		return
	}

	// Use a buffer to read data
	buffer := make([]byte, 1024) // You can adjust the buffer size
	var lineBuffer []byte
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Read data from the file into the buffer
			n, err := file.Read(buffer)
			if err != nil && err != io.EOF {
				log.Printf("Error reading file: %v", err)
				conn.WriteMessage(websocket.TextMessage, []byte("Error: Unable to read log file"))
				return
			}

			if n > 0 {
				// Process the buffer to detect lines
				for i := 0; i < n; i++ {
					// Accumulate bytes until a newline character is found
					if buffer[i] == '\n' {
						// When a newline is detected, send the line over WebSocket
						line := string(lineBuffer)
						err := conn.WriteMessage(websocket.TextMessage, []byte(line))
						if err != nil {
							log.Printf("Failed to send message: %v", err)
							return
						}
						log.Println("Log sent:", line)

						// Reset the lineBuffer to start accumulating the next line
						lineBuffer = nil
					} else {
						// Otherwise, accumulate the byte in the lineBuffer
						lineBuffer = append(lineBuffer, buffer[i])
					}
				}
			}

			// Reset file pointer to continue reading from the current position
			file.Seek(0, io.SeekCurrent)
		}
	}
}
