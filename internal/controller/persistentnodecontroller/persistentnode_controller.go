package persistentnodecontroller

import (
	"cynxhostagent/internal/usecase"

	"github.com/go-playground/validator/v10"
)

type PersistentNodeController struct {
	persistentNodeUsecase usecase.PersistentNodeUseCase
	validator             *validator.Validate
}

func New(
	persistentNodeUseCase usecase.PersistentNodeUseCase,
	validate *validator.Validate,
) *PersistentNodeController {
	return &PersistentNodeController{
		persistentNodeUsecase: persistentNodeUseCase,
		validator:             validate,
	}
}


