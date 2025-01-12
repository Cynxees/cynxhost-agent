package usecase

import (
	"context"
	"cynxhostagent/internal/model/request"
	"cynxhostagent/internal/model/response"
)

type PersistentNodeUseCase interface {
	RunPersistentNodeScript(ctx context.Context, req request.RunPersistentNodeScriptRequest, resp *response.APIResponse)
}
