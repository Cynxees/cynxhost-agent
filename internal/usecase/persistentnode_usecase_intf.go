package usecase

import (
	"context"
	"cynxhostagent/internal/model/request"
	"cynxhostagent/internal/model/response"
)

type PersistentNodeUseCase interface {
	RunPersistentNodeTemplateScript(ctx context.Context, req request.RunPersistentNodeTemplateScriptRequest, resp *response.APIResponse)
	SendCommand(ctx context.Context, req request.SendCommandRequest, resp *response.APIResponse)

	GetServerProperties(ctx context.Context, resp *response.APIResponse)
	SetServerProperties(ctx context.Context, req request.SetServerPropertiesRequest, resp *response.APIResponse)
}
