package request

import (
	contextmodel "cynxhostagent/internal/model/context"
	"cynxhostagent/internal/model/entity"
	"mime/multipart"
)

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

	PersistentNodeId int    `json:"persistent_node_id" validate:"required"`
	SessionId        string `json:"session_id" validate:"required"`
}

type StartSessionRequest struct {
	Shell string `json:"shell" validate:"required"`
}

type SendCommandRequest struct {
	Command         string `json:"command" validate:"required"`
	SessionId       string `json:"session_id" validate:"required"`
	IsBase64Encoded bool   `json:"is_base64_encoded"`
}

type SetServerPropertiesRequest struct {
	ServerProperties []entity.ServerProperty `json:"server_properties" validate:"required"`
}

type SendSingleDockerCommandRequest struct {
	Command string `json:"command" validate:"required"`
}

type DownloadFileRequest struct {
	FilePath string `json:"file_path" validate:"required"`
}

type RemoveFileRequest struct {
	FilePath string `json:"file_path" validate:"required"`
}

type UploadFileRequest struct {
	DestinationPath string               `json:"destination_path" validate:"required"`
	FileName        string               `json:"file_name" validate:"required"`
	FileData        multipart.File       `json:"file_data" validate:"required"`
	FileHeader      multipart.FileHeader `json:"file_data	" validate:"required"`
}

type ListDirectoryRequest struct {
	DirectoryPath string `json:"directory_path" validate:"required"`
}
