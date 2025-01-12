package tblinstancetype

import (
	"context"
	"cynxhostagent/internal/model/entity"
	"cynxhostagent/internal/model/request"
	"cynxhostagent/internal/repository/database"

	"gorm.io/gorm"
)

type TblInstanceTypeImpl struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) database.TblInstanceType {
	return &TblInstanceTypeImpl{
		DB: DB,
	}
}

func (database *TblInstanceTypeImpl) PaginateInstanceType(ctx context.Context, req request.PaginateRequest) (context.Context, []entity.TblInstanceType, error) {
	var instanceTypes []entity.TblInstanceType
	offset := (req.Page - 1) * req.Size
	err := database.DB.WithContext(ctx).Limit(req.Size).Offset(offset).Find(&instanceTypes).Error
	if err != nil {
		return ctx, nil, err
	}

	return ctx, instanceTypes, nil
}

func (database *TblInstanceTypeImpl) GetInstanceType(ctx context.Context, key string, value string) (context.Context, entity.TblInstanceType, error) {
	var instanceType entity.TblInstanceType
	err := database.DB.WithContext(ctx).Where(key+" = ?", value).First(&instanceType).Error
	if err != nil {
		return ctx, entity.TblInstanceType{}, err
	}

	return ctx, instanceType, nil
}
