package database

import (
	"context"
	"cynxhostagent/internal/model/entity"
)

type TblPersistentNodeImage interface {
	CreatePersistentNodeImage(ctx context.Context, persistentNodeImage entity.TblPersistentNodeImage) (context.Context, int, error)
	GetPersistentNodeImages(ctx context.Context, key string, value string) (context.Context, []entity.TblPersistentNodeImage, error)
	UpdatePersistentNodeImage(ctx context.Context, id int, persistentNodeImage entity.TblPersistentNodeImage) (context.Context, error)
	DeletePersistentNodeImage(ctx context.Context, id int) (context.Context, error)
}
