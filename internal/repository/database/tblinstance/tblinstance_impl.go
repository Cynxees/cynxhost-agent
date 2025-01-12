package tblinstance

import (
	"context"
	"cynxhostagent/internal/model/entity"
	"cynxhostagent/internal/repository/database"

	"gorm.io/gorm"
)

type TblInstanceImpl struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) database.TblInstance {
	return &TblInstanceImpl{
		DB: DB,
	}
}

func (database *TblInstanceImpl) CreateInstance(ctx context.Context, instance entity.TblInstance) (context.Context, int, error) {
	err := database.DB.WithContext(ctx).Create(&instance).Error
	if err != nil {
		return ctx, 0, err
	}

	return ctx, instance.Id, nil
}

func (database *TblInstanceImpl) GetInstances(ctx context.Context, key string, value string) (context.Context, []entity.TblInstance, error) {
	var instances []entity.TblInstance
	err := database.DB.WithContext(ctx).Where(key+" = ?", value).Find(&instances).Error
	if err != nil {
		return ctx, nil, err
	}

	return ctx, instances, nil
}

func (database *TblInstanceImpl) UpdateInstance(ctx context.Context, id int, instance entity.TblInstance) (context.Context, error) {
	err := database.DB.WithContext(ctx).Model(&entity.TblInstance{}).Where("id = ?", id).Updates(instance).Error
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (database *TblInstanceImpl) DeleteInstance(ctx context.Context, id int) (context.Context, error) {
	err := database.DB.WithContext(ctx).Where("id = ?", id).Delete(&entity.TblInstance{}).Error
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}
