package tblparameter

import (
	"context"
	"cynxhostagent/internal/model/entity"
	"cynxhostagent/internal/repository/database"

	"gorm.io/gorm"
)

type TblParameterImpl struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) database.TblParameter {
	return &TblParameterImpl{
		DB: DB,
	}
}

func (database *TblParameterImpl) SelectTblParameters(ctx context.Context, ids []string) (context.Context, []entity.TblParameter, error) {

	paramDatas := []entity.TblParameter{}

	err := database.DB.WithContext(ctx).Where("id IN (?)", ids).Find(&paramDatas).Error
	if err != nil {
		return ctx, nil, err
	}

	return ctx, paramDatas, nil
}
