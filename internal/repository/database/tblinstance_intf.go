package database

import (
	"context"
	"cynxhostagent/internal/model/entity"
)

type TblInstance interface {
	CreateInstance(ctx context.Context, instance entity.TblInstance) (context.Context, int, error)
	GetInstances(ctx context.Context, key string, value string) (context.Context, []entity.TblInstance, error)
	UpdateInstance(ctx context.Context, id int, instance entity.TblInstance) (context.Context, error)
	DeleteInstance(ctx context.Context, id int) (context.Context, error)
}
