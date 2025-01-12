package usecase

import (
	"context"
	"cynxhostagent/internal/model/request"
	"cynxhostagent/internal/model/response"
)

type UserUseCase interface {
	BypassLoginUser(ctx context.Context, req request.BypassLoginUserRequest, resp *response.APIResponse) context.Context
}
