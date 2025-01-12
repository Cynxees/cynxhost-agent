package types

type StorageStatus string
type InstanceTypeStatus string
type InstanceStatus string
type PersistentNodeStatus string

const (
	StorageStatusInUseRunningInstance  StorageStatus = "IN_USE:RUNNING_INSTANCE"
	StorageStatusInUseCreatingSnapshot StorageStatus = "IN_USE:CREATING_SNAPSHOT"
	StorageStatusReady                 StorageStatus = "READY"
	StorageStatusNew                   StorageStatus = "NEW"
)

const (

	// No Process
	PersistentNodeStatusRunning  PersistentNodeStatus = "RUNNING"
	PersistentNodeStatusStopped  PersistentNodeStatus = "STOPPED"
	PersistentNodeStatusShutdown PersistentNodeStatus = "SHUTDOWN"

	// Processing ( in the middle of running script or something )
	PersistentNodeStatusCreating     PersistentNodeStatus = "CREATING"
	PersistentNodeStatusSetup        PersistentNodeStatus = "SETUP"
	PersistentNodeStatusStarting     PersistentNodeStatus = "STARTING"
	PersistentNodeStatusStopping     PersistentNodeStatus = "STOPPING"
	PersistentNodeStatusShuttingDown PersistentNodeStatus = "SHUTTING_DOWN"
)

const (
	InstanceStatusCreate   InstanceStatus = "CREATING"
	InstanceStatusActive   InstanceStatus = "ACTIVE"
	InstanceStatusInactive InstanceStatus = "INACTIVE"
)
