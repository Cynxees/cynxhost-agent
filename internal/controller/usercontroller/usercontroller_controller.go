package usercontroller

import (
	"context"
	"cynxhostagent/internal/helper"
	"cynxhostagent/internal/model/request"
	"cynxhostagent/internal/model/response"
	"cynxhostagent/internal/model/response/responsecode"
	"cynxhostagent/internal/usecase"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type UserController struct {
	userUsecase usecase.UserUseCase
	validator   *validator.Validate
}

func New(
	userUseCase usecase.UserUseCase,
	validate *validator.Validate,
) *UserController {
	return &UserController{
		userUsecase: userUseCase,
		validator:   validate,
	}
}

func (controller *UserController) BypassLoginUser(w http.ResponseWriter, r *http.Request) (context.Context, response.APIResponse) {
	var requestBody request.BypassLoginUserRequest
	var apiResponse response.APIResponse

	ctx := r.Context()
	requestBody.ClientIp = helper.GetClientIP(r)

	if err := helper.DecodeAndValidateRequest(r, &requestBody, controller.validator); err != nil {
		apiResponse.Code = responsecode.CodeValidationError
		apiResponse.Error = err.Error()
		return ctx, apiResponse
	}

	ctx = controller.userUsecase.BypassLoginUser(ctx, requestBody, &apiResponse)

	return ctx, apiResponse
}