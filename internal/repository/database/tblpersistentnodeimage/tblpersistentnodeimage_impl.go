package tblpersistentnodeimage

import (
	"context"
	"cynxhost/internal/model/entity"
	"cynxhost/internal/repository/database"

	"gorm.io/gorm"
)

type TblPersistentNodeImageImpl struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) database.TblPersistentNodeImage {
	return &TblPersistentNodeImageImpl{
		DB: DB,
	}
}

func (database *TblPersistentNodeImageImpl) CreatePersistentNodeImage(ctx context.Context, persistentNodeImage entity.TblPersistentNodeImage) (context.Context, int, error) {
	err := database.DB.WithContext(ctx).Create(&persistentNodeImage).Error
	if err != nil {
		return ctx, 0, err
	}

	return ctx, persistentNodeImage.Id, nil
}

func (database *TblPersistentNodeImageImpl) GetPersistentNodeImages(ctx context.Context, key string, value string) (context.Context, []entity.TblPersistentNodeImage, error) {
	var persistentNodeImages []entity.TblPersistentNodeImage
	err := database.DB.WithContext(ctx).Where(key+" = ?", value).Find(&persistentNodeImages).Error
	if err != nil {
		return ctx, nil, err
	}

	return ctx, persistentNodeImages, nil
}

func (database *TblPersistentNodeImageImpl) UpdatePersistentNodeImage(ctx context.Context, id int, persistentNodeImage entity.TblPersistentNodeImage) (context.Context, error) {
	err := database.DB.WithContext(ctx).Model(&entity.TblPersistentNodeImage{}).Where("id = ?", id).Updates(persistentNodeImage).Error
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (database *TblPersistentNodeImageImpl) DeletePersistentNodeImage(ctx context.Context, id int) (context.Context, error) {
	err := database.DB.WithContext(ctx).Where("id = ?", id).Delete(&entity.TblPersistentNodeImage{}).Error
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}
