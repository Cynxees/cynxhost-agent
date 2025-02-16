package usecase

import (
	"context"
	"cynxhostagent/internal/model/request"
	"cynxhostagent/internal/model/response"
)

type PersistentNodeUseCase interface {

	// Dashboard
	RunPersistentNodeTemplateScript(ctx context.Context, req request.RunPersistentNodeTemplateScriptRequest, resp *response.APIResponse)
	GetNodeContainerStats(ctx context.Context, resp *response.APIResponse)

	// Console
	StreamLogs(ctx context.Context, req request.GetPersistentNodeRealTimeLogsRequest, channel chan string) error
	CreateSession(ctx context.Context, req request.StartSessionRequest, resp *response.APIResponse)
	SendCommand(ctx context.Context, req request.SendCommandRequest, resp *response.APIResponse)

	// Files
	SendSingleDockerCommand(ctx context.Context, req request.SendSingleDockerCommandRequest, resp *response.APIResponse)
	DownloadFile(ctx context.Context, req request.DownloadFileRequest, resp *response.APIResponse) (file []byte, err error)
	UploadFile(ctx context.Context, req request.UploadFileRequest, resp *response.APIResponse)
	RemoveFile(ctx context.Context, req request.RemoveFileRequest, resp *response.APIResponse)
	ListDirectory(ctx context.Context, req request.ListDirectoryRequest, resp *response.APIResponse)
	CreateDirectory(ctx context.Context, req request.CreateDirectoryRequest, resp *response.APIResponse)
	RemoveDirectory(ctx context.Context, req request.RemoveDirectoryRequest, resp *response.APIResponse)

	// Backup
	PushImage(ctx context.Context, resp *response.APIResponse)
	ListImages(ctx context.Context, resp *response.APIResponse)

	// Server Properties
	// GetServerProperties(ctx context.Context, resp *response.APIResponse)
	// SetServerProperties(ctx context.Context, req request.SetServerPropertiesRequest, resp *response.APIResponse)

	// TMUX Console
	// StreamTmuxLogs(ctx context.Context, resp *response.APIResponse)
	// SendTmuxCommand(ctx context.Context, req request.SendCommandRequest, resp *response.APIResponse)
}
