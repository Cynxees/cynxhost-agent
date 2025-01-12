package request

import contextmodel "cynxhostagent/internal/model/context"

type PaginateRequest struct {
	Page int `json:"page" validate:"required"`
	Size int `json:"size" validate:"required"`
}

type BypassLoginUserRequest struct {
	ClientIp string `validate:"required"`

	UserId int `json:"user_id" validate:"required"`
}

type RunPersistentNodeTemplateScriptRequest struct {
	SessionUser contextmodel.User `validate:"required"`

	PersistentNodeId int    `json:"persistent_node_id" validate:"required"`
	ScriptType       string `json:"script_type" validate:"required"`
}

type GetPersistentNodeRealTimeLogsRequest struct {
	SessionUser contextmodel.User `validate:"required"`

	PersistentNodeId int `json:"persistent_node_id" validate:"required"`
}
