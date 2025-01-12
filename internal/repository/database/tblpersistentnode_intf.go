package database

import (
	"context"
	"cynxhostagent/internal/model/entity"
)

type TblPersistentNode interface {
	CreatePersistentNode(ctx context.Context, persistentNode entity.TblPersistentNode) (context.Context, int, error)
	GetPersistentNodes(ctx context.Context, key string, value string) (context.Context, []entity.TblPersistentNode, error)
	UpdatePersistentNode(ctx context.Context, id int, persistentNode entity.TblPersistentNode) (context.Context, error)
	DeletePersistentNode(ctx context.Context, id int) (context.Context, error)
}
