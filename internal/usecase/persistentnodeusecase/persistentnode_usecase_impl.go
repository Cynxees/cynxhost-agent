package persistentnodeusecase

import (
	"context"
	"cynxhostagent/internal/dependencies"
	"cynxhostagent/internal/model/request"
	"cynxhostagent/internal/model/response"
	"cynxhostagent/internal/repository/database"
	"cynxhostagent/internal/usecase"

	"github.com/sirupsen/logrus"
)

type PersistentNodeUseCaseImpl struct {
	tblPersistentNode database.TblPersistentNode
	tblInstance       database.TblInstance
	tblInstanceType   database.TblInstanceType
	tblStorage        database.TblStorage
	tblServerTemplate database.TblServerTemplate

	awsClient *dependencies.AWSClient
	log       *logrus.Logger
	config    *dependencies.Config
}

func New(tblPersistentNode database.TblPersistentNode, tblInstance database.TblInstance, tblInstanceType database.TblInstanceType, tblStorage database.TblStorage, tblServerTemplate database.TblServerTemplate, awsClient *dependencies.AWSClient, logger *logrus.Logger, config *dependencies.Config) usecase.PersistentNodeUseCase {

	return &PersistentNodeUseCaseImpl{
		tblPersistentNode: tblPersistentNode,
		tblStorage:        tblStorage,
		tblServerTemplate: tblServerTemplate,
		tblInstance:       tblInstance,
		tblInstanceType:   tblInstanceType,

		awsClient: awsClient,
		log:       logger,
		config:    config,
	}
}

func (usecase *PersistentNodeUseCaseImpl) RunPersistentNodeScript(ctx context.Context, req request.RunPersistentNodeScriptRequest, resp *response.APIResponse) {

}
