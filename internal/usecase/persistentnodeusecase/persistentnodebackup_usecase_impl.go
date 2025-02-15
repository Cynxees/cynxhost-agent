package persistentnodeusecase

import (
	"context"
	"cynxhostagent/internal/model/response"
	"cynxhostagent/internal/model/response/responsecode"
	"cynxhostagent/internal/model/response/responsedata"
	"fmt"
	"strconv"
)

func (uc *PersistentNodeUseCaseImpl) BackupImage(ctx context.Context, resp *response.APIResponse) {

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
