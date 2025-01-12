package request

type BypassLoginUserRequest struct {
	ClientIp string `validate:"required"`

	UserId int `json:"user_id" validate:"required"`
}

type RunPersistentNodeScriptRequest struct {
	ClientIp string `validate:"required"`

	PersistentNodeId string `json:"persistent_node_id" validate:"required"`
	ScriptType       string `json:"script_type" validate:"required"`
}

type PaginateRequest struct {
	Page int `json:"page" validate:"required"`
	Size int `json:"size" validate:"required"`
}
