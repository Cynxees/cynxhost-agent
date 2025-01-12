package database

import (
	"context"
	"cynxhostagent/internal/model/entity"
)

type TblServerTemplate interface {
	CreateServerTemplate(ctx context.Context, serverTemplate entity.TblServerTemplate) (context.Context, int, error)
	GetServerTemplate(ctx context.Context, key string, value string) (context.Context, entity.TblServerTemplate, error)
	PaginateServerTemplate(ctx context.Context, page int, size int) (context.Context, []entity.TblServerTemplate, error)
	DeleteServerTemplate(ctx context.Context, key string, value string) (context.Context, error)
}
