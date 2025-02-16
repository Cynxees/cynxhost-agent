package persistentnodeusecase

import (
	"context"
	"cynxhostagent/internal/constant"
	"cynxhostagent/internal/constant/types"
	"cynxhostagent/internal/helper"
	"cynxhostagent/internal/model/entity"
	"cynxhostagent/internal/model/response"
	"cynxhostagent/internal/model/response/responsecode"
	"cynxhostagent/internal/model/response/responsedata"
	"errors"
	"fmt"
	"strconv"
)

func (uc *PersistentNodeUseCaseImpl) PushImage(ctx context.Context, resp *response.APIResponse) {

	repositoryName := uc.config.Aws.Ecr.RepositoryPrefix + "/" + strconv.Itoa(uc.config.App.PersistentNode.Id)

	_, err := uc.awsClient.GetRepository(ctx, repositoryName)

	switch err {
	case nil:
	// Repository exists
	case errors.New(constant.NotFound):
		// Repository does not exist
		_, err = uc.awsClient.CreateRepository(ctx, repositoryName)
		if err != nil {
			resp.Code = responsecode.CodeAWSError
			resp.Error = fmt.Sprintf("Error creating repository: %v", err)
			return
		}

	default:
		resp.Code = responsecode.CodeAWSError
		resp.Error = fmt.Sprintf("Error getting repository: %v", err)
		return
	}

	// Get current timestamp for tagging the image in yyyymmddhhmmss
	imageTag := helper.GetCurrentTimestampYYYYMMDDHHMMSS()

	err = uc.dockerManager.UploadImageToAwsEcr(uc.config.DockerConfig.ContainerName, repositoryName, imageTag, uc.config.Aws.Ecr)
	if err != nil {
		resp.Code = responsecode.CodeDockerError
		resp.Error = fmt.Sprintf("Error pushing image: %v", err)
		return
	}

	// Save the image to the database
	_, _, err = uc.tblPersistentNodeImage.CreatePersistentNodeImage(ctx, entity.TblPersistentNodeImage{
		PersistentNodeId: uc.config.App.PersistentNode.Id,
		ImageTag:         imageTag,
		Status:           types.PersistentNodeImageStatus(types.PersistentNodeImageStatusActive),
	})
	if err != nil {
		resp.Code = responsecode.CodeTblPersistentNodeImageError
		resp.Error = fmt.Sprintf("Error saving image to database: %v", err)
		return
	}

	resp.Code = responsecode.CodeSuccess
}

func (uc *PersistentNodeUseCaseImpl) ListImages(ctx context.Context, resp *response.APIResponse) {

	ctx, images, err := uc.tblPersistentNodeImage.GetPersistentNodeImages(ctx, "persistent_node_id", strconv.Itoa(uc.config.App.PersistentNode.Id))
	if err != nil {
		resp.Code = responsecode.CodeTblPersistentNodeImageError
		resp.Error = fmt.Sprintf("Error listing images: %v", err)
		return
	}

	resp.Code = responsecode.CodeSuccess
	resp.Data = responsedata.ListPersistentNodeImagesResponseData{
		Images: images,
	}
}
