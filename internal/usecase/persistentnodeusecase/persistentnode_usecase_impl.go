package persistentnodeusecase

import (
	"context"
	"cynxhostagent/internal/constant/types"
	"cynxhostagent/internal/dependencies"
	"cynxhostagent/internal/model/request"
	"cynxhostagent/internal/model/response"
	"cynxhostagent/internal/model/response/responsecode"
	"cynxhostagent/internal/repository/database"
	"cynxhostagent/internal/usecase"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

type PersistentNodeUseCaseImpl struct {
	tblPersistentNode database.TblPersistentNode
	tblInstance       database.TblInstance
	tblInstanceType   database.TblInstanceType
	tblStorage        database.TblStorage
	tblServerTemplate database.TblServerTemplate

	awsClient   *dependencies.AWSClient
	log         *logrus.Logger
	config      *dependencies.Config
	osManager   *dependencies.OSManager
	tmuxManager *dependencies.TmuxManager
}

func New(tblPersistentNode database.TblPersistentNode, tblInstance database.TblInstance, tblInstanceType database.TblInstanceType, tblStorage database.TblStorage, tblServerTemplate database.TblServerTemplate, awsClient *dependencies.AWSClient, logger *logrus.Logger, config *dependencies.Config, osManager *dependencies.OSManager, tmuxManager *dependencies.TmuxManager) usecase.PersistentNodeUseCase {

	return &PersistentNodeUseCaseImpl{
		tblPersistentNode: tblPersistentNode,
		tblStorage:        tblStorage,
		tblServerTemplate: tblServerTemplate,
		tblInstance:       tblInstance,
		tblInstanceType:   tblInstanceType,

		awsClient:   awsClient,
		log:         logger,
		config:      config,
		osManager:   osManager,
		tmuxManager: tmuxManager,
	}
}

func (usecase *PersistentNodeUseCaseImpl) RunPersistentNodeTemplateScript(ctx context.Context, req request.RunPersistentNodeTemplateScriptRequest, resp *response.APIResponse) {

	// Get the persistent node
	ctx, persistentNodes, err := usecase.tblPersistentNode.GetPersistentNodes(ctx, "id", strconv.Itoa(req.PersistentNodeId))
	if err != nil {
		resp.Code = responsecode.CodeTblPersistentNodeError
		resp.Error = err.Error()
		return
	}

	if len(persistentNodes) == 0 {
		resp.Code = responsecode.CodeNotFound
		resp.Error = "Persistent node not found"
		return
	}

	persistentNode := persistentNodes[0]

	// Check User
	if persistentNode.OwnerId != req.SessionUser.Id {
		resp.Code = responsecode.CodeAuthenticationError
		resp.Error = "Unauthorized user"
		return
	}

	// Get the script
	ctx, serverTemplate, err := usecase.tblServerTemplate.GetServerTemplate(ctx, "id", strconv.Itoa(persistentNode.ServerTemplateId))
	if err != nil {
		resp.Code = responsecode.CodeTblServerTemplateError
		resp.Error = err.Error()
		return
	}

	// Get the script
	var script string

	switch req.ScriptType {
	case string(types.ScriptTypeSetup):
		script = serverTemplate.Script.SetupScript

	case string(types.ScriptTypeStart):
		script = serverTemplate.Script.StartScript

	case string(types.ScriptTypeStop):
		script = serverTemplate.Script.StopScript

	case string(types.ScriptTypeShutdown):
		script = serverTemplate.Script.ShutdownScript

	default:
		resp.Code = responsecode.CodeValidationError
		resp.Error = "Invalid script type"
		return
	}

	// Run the script
	err = usecase.osManager.RunBashScript(script)

	if err != nil {
		resp.Code = responsecode.CodeOsError
		resp.Error = "Error running script " + err.Error()
		return
	}

	resp.Code = responsecode.CodeSuccess
}

func (usecase *PersistentNodeUseCaseImpl) SendCommand(ctx context.Context, req request.SendCommandRequest, resp *response.APIResponse) {

	err := usecase.tmuxManager.SendCommand(usecase.config.Tmux.SessionName, req.Command)

	if err != nil {
		resp.Code = responsecode.CodeTmuxError
		resp.Error = "Error running script " + err.Error()
		return
	}

	resp.Code = responsecode.CodeSuccess
}

func (usecase *PersistentNodeUseCaseImpl) GetServerProperties(ctx context.Context, resp *response.APIResponse) {
	properties, err := usecase.osManager.ReadServerProperties(usecase.config.Files.MinecraftServerProperties)
	if err != nil {
		resp.Error = fmt.Sprintf("Error reading server.properties: %v", err)
		resp.Code = responsecode.CodeFailed
		return
	}
	resp.Code = responsecode.CodeSuccess
	resp.Data = properties
}

func (usecase *PersistentNodeUseCaseImpl) SetServerProperties(ctx context.Context, req request.SetServerPropertiesRequest, resp *response.APIResponse) {

	err := usecase.osManager.SetServerProperties(usecase.config.Files.MinecraftServerProperties, req.ServerProperties)
	if err != nil {
		resp.Error = fmt.Sprintf("Error reading server.properties: %v", err)
		resp.Code = responsecode.CodeFailed
		return
	}
	resp.Code = responsecode.CodeSuccess
}
