package userusecase

import (
	"context"
	"cynxhostagent/internal/dependencies"
	"cynxhostagent/internal/model/request"
	"cynxhostagent/internal/model/response"
	"cynxhostagent/internal/model/response/responsecode"
	"cynxhostagent/internal/model/response/responsedata"
	"cynxhostagent/internal/repository/database"
	"cynxhostagent/internal/usecase"
	"strconv"
)

type UserUseCaseImpl struct {
	tblUser    database.TblUser
	jwtManager *dependencies.JWTManager
	config     *dependencies.Config
}

func New(tblUser database.TblUser, jwtManager *dependencies.JWTManager, config *dependencies.Config) usecase.UserUseCase {
	return &UserUseCaseImpl{
		tblUser:    tblUser,
		jwtManager: jwtManager,
		config:     config,
	}
}

func (usecase *UserUseCaseImpl) BypassLoginUser(ctx context.Context, req request.BypassLoginUserRequest, resp *response.APIResponse) context.Context {

	if req.ClientIp != usecase.config.Central.PrivateIp && req.ClientIp != usecase.config.Central.PublicIp {
		resp.Code = responsecode.CodeAuthenticationError
		resp.Error = "Invalid client IP"
		return ctx
	}

	_, user, err := usecase.tblUser.GetUser(ctx, "id", strconv.Itoa(req.UserId))
	if err != nil {
		resp.Code = responsecode.CodeTblUserError
		resp.Error = err.Error()
		return ctx
	}

	if user == nil {
		resp.Code = responsecode.CodeNotFound
		resp.Error = "User not found"
		return ctx
	}

	token, err := usecase.jwtManager.GenerateToken(user.Id)
	if err != nil {
		resp.Code = responsecode.CodeJwtError
		resp.Error = err.Error()
		return ctx
	}

	resp.Code = responsecode.CodeSuccess
	resp.Data = responsedata.AuthResponseData{
		AccessToken: token.AccessToken,
		TokenType:   "Bearer",
	}
	return ctx
}
