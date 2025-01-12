package database

import (
	"context"
	"cynxhostagent/internal/model/entity"
	"cynxhostagent/internal/model/request"
)

type TblInstanceType interface {
	PaginateInstanceType(ctx context.Context, req request.PaginateRequest) (context.Context, []entity.TblInstanceType, error)
	GetInstanceType(ctx context.Context, key string, value string) (context.Context, entity.TblInstanceType, error)
}
