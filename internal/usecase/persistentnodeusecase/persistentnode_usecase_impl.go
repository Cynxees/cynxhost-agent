package persistentnodeusecase

import (
	"context"
	"cynxhostagent/internal/constant/types"
	"cynxhostagent/internal/dependencies"
	"cynxhostagent/internal/model/request"
	"cynxhostagent/internal/model/response"
	"cynxhostagent/internal/model/response/responsecode"
	"cynxhostagent/internal/model/response/responsedata"
	"cynxhostagent/internal/repository/database"
	"cynxhostagent/internal/repository/micro/cynxhostcentral"
	"cynxhostagent/internal/usecase"
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/sirupsen/logrus"
)

type PersistentNodeUseCaseImpl struct {
	tblPersistentNode database.TblPersistentNode
	tblInstance       database.TblInstance
	tblInstanceType   database.TblInstanceType
	tblStorage        database.TblStorage
	tblServerTemplate database.TblServerTemplate

	awsClient       *dependencies.AWSClient
	cynxhostcentral *cynxhostcentral.CynxhostCentral

	log           *logrus.Logger
	config        *dependencies.Config
	osManager     *dependencies.OSManager
	dockerManager *dependencies.DockerManager
}

func New(tblPersistentNode database.TblPersistentNode, tblInstance database.TblInstance, tblInstanceType database.TblInstanceType, tblStorage database.TblStorage, tblServerTemplate database.TblServerTemplate, awsClient *dependencies.AWSClient, logger *logrus.Logger, config *dependencies.Config, osManager *dependencies.OSManager, dockerManager *dependencies.DockerManager, cynxhostCentral *cynxhostcentral.CynxhostCentral) usecase.PersistentNodeUseCase {

	return &PersistentNodeUseCaseImpl{
		tblPersistentNode: tblPersistentNode,
		tblStorage:        tblStorage,
		tblServerTemplate: tblServerTemplate,
		tblInstance:       tblInstance,
		tblInstanceType:   tblInstanceType,

		awsClient:       awsClient,
		cynxhostcentral: cynxhostCentral,

		log:           logger,
		config:        config,
		osManager:     osManager,
		dockerManager: dockerManager,
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

	// If shutdown, Call callback in central
	if req.ScriptType == string(types.PersistentNodeStatusShutdown) {

		err = usecase.cynxhostcentral.CallShutdownCallback()
		if err != nil {
			resp.Code = responsecode.CodeCentralError
			resp.Error = "Error calling central " + err.Error()
			return
		}
	}

	resp.Code = responsecode.CodeSuccess
}

func (usecase *PersistentNodeUseCaseImpl) GetServerProperties(ctx context.Context, resp *response.APIResponse) {
	properties, err := usecase.osManager.ReadServerProperties(usecase.config.DockerConfig.Files.MinecraftServerProperties)
	if err != nil {
		resp.Error = fmt.Sprintf("Error reading server.properties: %v", err)
		resp.Code = responsecode.CodeFailed
		return
	}
	resp.Code = responsecode.CodeSuccess
	resp.Data = properties
}

func (usecase *PersistentNodeUseCaseImpl) SetServerProperties(ctx context.Context, req request.SetServerPropertiesRequest, resp *response.APIResponse) {

	err := usecase.osManager.SetServerProperties(usecase.config.DockerConfig.Files.MinecraftServerProperties, req.ServerProperties)
	if err != nil {
		resp.Error = fmt.Sprintf("Error reading server.properties: %v", err)
		resp.Code = responsecode.CodeFailed
		return
	}
	resp.Code = responsecode.CodeSuccess
}

func (uc *PersistentNodeUseCaseImpl) StreamLogs(ctx context.Context, req request.GetPersistentNodeRealTimeLogsRequest, channel chan string) error {
	// Start streaming logs from a specific container
	go func() {
		// Start the interactive session
		err := uc.dockerManager.StreamOutput(req.SessionId, channel)
		if err != nil {
			fmt.Println("Error starting interactive session:", err)
			channel <- fmt.Sprintf("Error starting interactive session: %v", err)
			return
		}

		// Stream the logs and send them to the channel
		// Assuming the `StartInteractiveSession` is configured to stream logs to a WebSocket channel
		for logLine := range channel {
			// Sending logs to the channel for WebSocket transmission
			channel <- logLine
		}
	}()

	return nil
}

func (uc *PersistentNodeUseCaseImpl) CreateSession(ctx context.Context, req request.StartSessionRequest, resp *response.APIResponse) {
	// Create a new interactive session
	uniqueCode := strconv.Itoa(rand.Int())
	fmt.Println("Creating new session with code ", uniqueCode)
	err := uc.dockerManager.CreateNewSession(uniqueCode, uc.config.DockerConfig.ContainerName, req.Shell)
	if err != nil {
		resp.Code = responsecode.CodeFailed
		resp.Error = fmt.Sprintf("Error creating new session: %v", err)
		return
	}

	resp.Code = responsecode.CodeSuccess
	resp.Data = responsedata.CreateSessionResponseData{
		SessionId: uniqueCode,
	}
}

func (uc *PersistentNodeUseCaseImpl) SendCommand(ctx context.Context, req request.SendCommandRequest, resp *response.APIResponse) {

	// Send the command to the Docker container
	err := uc.dockerManager.SendCommand(req.SessionId, req.Command, req.IsBase64Encoded)
	if err != nil {
		resp.Code = responsecode.CodeFailed
		resp.Error = fmt.Sprintf("Error sending command to Docker container: %v", err)
		return
	}

	resp.Code = responsecode.CodeSuccess
}
