package types

type LaunchCallbackPersistentNodeType string

const (
	LaunchCallbackPersistentNodeTypeInitialLaunch LaunchCallbackPersistentNodeType = "INITIAL_LAUNCH"
)

type StatusCallbackPersistentNodeType string

const (
	SetupSuccessCallbackPersistentNodeType StatusCallbackPersistentNodeType = "SETUP_SUCCESS"
)

type ScriptType string

const (
	ScriptTypeSetup    ScriptType = "SETUP"
	ScriptTypeStart    ScriptType = "START"
	ScriptTypeStop     ScriptType = "STOP"
	ScriptTypeShutdown ScriptType = "SHUTDOWN"
)
