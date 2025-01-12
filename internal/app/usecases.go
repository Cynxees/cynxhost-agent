package app

import (
	"cynxhostagent/internal/usecase"
	"cynxhostagent/internal/usecase/persistentnodeusecase"
	"cynxhostagent/internal/usecase/userusecase"
)

type UseCases struct {
	UserUseCase           usecase.UserUseCase
	PersistentNodeUseCase usecase.PersistentNodeUseCase
}

func NewUseCases(repos *Repos, dependencies *Dependencies) *UseCases {

	return &UseCases{
		UserUseCase:           userusecase.New(repos.TblUser, repos.JWTManager, dependencies.Config),
		PersistentNodeUseCase: persistentnodeusecase.New(repos.TblPersistentNode, repos.TblInstance, repos.TblInstanceType, repos.TblStorage, repos.TblServerTemplate, dependencies.AWSClient, dependencies.Logger, dependencies.Config),
	}
}
