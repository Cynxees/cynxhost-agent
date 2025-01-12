package tblstorage

import (
	"context"
	"cynxhostagent/internal/model/entity"
	"cynxhostagent/internal/repository/database"

	"gorm.io/gorm"
)

type TblStorageImpl struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) database.TblStorage {
	return &TblStorageImpl{
		DB: DB,
	}
}

func (database *TblStorageImpl) GetStorages(ctx context.Context, key string, value string) (context.Context, []entity.TblInstanceType, error) {
	var storages []entity.TblInstanceType
	err := database.DB.WithContext(ctx).Where(key+" = ?", value).Find(&storages).Error
	if err != nil {
		return ctx, nil, err
	}

	return ctx, storages, nil
}

func (database *TblStorageImpl) CreateStorage(ctx context.Context, storage entity.TblStorage) (context.Context, int, error) {
	err := database.DB.WithContext(ctx).Create(&storage).Error
	if err != nil {
		return ctx, 0, err
	}

	return ctx, storage.Id, nil
}

func (database *TblStorageImpl) UpdateStorage(ctx context.Context, id int, storage entity.TblStorage) (context.Context, error) {
	err := database.DB.WithContext(ctx).Model(&entity.TblStorage{}).Where("id = ?", id).Updates(storage).Error
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (database *TblStorageImpl) DeleteStorage(ctx context.Context, id int) (context.Context, error) {
	err := database.DB.WithContext(ctx).Where("id = ?", id).Delete(&entity.TblStorage{}).Error
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}
