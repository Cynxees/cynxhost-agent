package database

import (
	"context"
	"cynxhostagent/internal/model/entity"
)

type TblParameter interface {
	SelectTblParameters(ctx context.Context, ids []string) (context.Context, []entity.TblParameter, error)
}
