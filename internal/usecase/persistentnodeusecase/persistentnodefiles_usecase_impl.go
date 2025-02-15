package persistentnodeusecase

import (
	"context"
	"cynxhostagent/internal/model/request"
	"cynxhostagent/internal/model/response"
	"cynxhostagent/internal/model/response/responsecode"
	"cynxhostagent/internal/model/response/responsedata"
	"fmt"
)

func (uc *PersistentNodeUseCaseImpl) DownloadFile(ctx context.Context, req request.DownloadFileRequest, resp *response.APIResponse) (file []byte, err error) {
	// Download the file from the server
	fileData, err := uc.dockerManager.GetFile(uc.config.DockerConfig.ContainerName, req.FilePath)
	if err != nil {
		resp.Code = responsecode.CodeDockerError
		resp.Error = fmt.Sprintf("Error reading file: %v", err)
		return nil, err
	}

	resp.Code = responsecode.CodeSuccess
	return fileData, nil

}

func (uc *PersistentNodeUseCaseImpl) UploadFile(ctx context.Context, req request.UploadFileRequest, resp *response.APIResponse) {

	// Upload the file to the server
	err := uc.dockerManager.WriteFile(req.DestinationPath, req.FileData, req.FileHeader, uc.config.DockerConfig.ContainerName, req.FileName)
	if err != nil {
		resp.Code = responsecode.CodeDockerError
		resp.Error = fmt.Sprintf("Error writing file: %v", err)
		return
	}

	resp.Code = responsecode.CodeSuccess
}

func (uc *PersistentNodeUseCaseImpl) RemoveFile(ctx context.Context, req request.RemoveFileRequest, resp *response.APIResponse) {
	// Remove the file from the server
	err := uc.dockerManager.RemoveFile(uc.config.DockerConfig.ContainerName, req.FilePath)
	if err != nil {
		resp.Code = responsecode.CodeDockerError
		resp.Error = fmt.Sprintf("Error removing file: %v", err)
		return
	}

	resp.Code = responsecode.CodeSuccess
}

func (uc *PersistentNodeUseCaseImpl) ListDirectory(ctx context.Context, req request.ListDirectoryRequest, resp *response.APIResponse) {
	// List the directory contents
	files, err := uc.dockerManager.ListDirectory(uc.config.DockerConfig.ContainerName, req.DirectoryPath)
	if err != nil {
		resp.Code = responsecode.CodeDockerError
		resp.Error = fmt.Sprintf("Error listing directory: %v", err)
		return
	}

	resp.Code = responsecode.CodeSuccess
	resp.Data = responsedata.ListDirectoryResponseData{
		Files: files,
	}
}
