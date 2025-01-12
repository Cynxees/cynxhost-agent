package database

import (
	"context"
	"cynxhostagent/internal/model/entity"
)

type TblStorage interface {
	GetStorages(ctx context.Context, key string, value string) (context.Context, []entity.TblInstanceType, error)
	CreateStorage(ctx context.Context, storage entity.TblStorage) (context.Context, int, error)
	UpdateStorage(ctx context.Context, id int, storage entity.TblStorage) (context.Context, error)
	DeleteStorage(ctx context.Context, id int) (context.Context, error)
}
