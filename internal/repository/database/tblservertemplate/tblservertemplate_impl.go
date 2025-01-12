package tblservertemplate

import (
	"context"
	"cynxhostagent/internal/model/entity"
	"cynxhostagent/internal/repository/database"

	"gorm.io/gorm"
)

type TblServerTemplateImpl struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) database.TblServerTemplate {
	return &TblServerTemplateImpl{
		DB: DB,
	}
}

func (database *TblServerTemplateImpl) CreateServerTemplate(ctx context.Context, serverTemplate entity.TblServerTemplate) (context.Context, int, error) {
	err := database.DB.WithContext(ctx).Create(&serverTemplate).Error
	if err != nil {
		return ctx, 0, err
	}

	return ctx, serverTemplate.Id, nil
}

func (database *TblServerTemplateImpl) GetServerTemplate(ctx context.Context, key string, value string) (context.Context, entity.TblServerTemplate, error) {
	var serverTemplate entity.TblServerTemplate
	err := database.DB.WithContext(ctx).Preload("Script").Where(key+" = ?", value).First(&serverTemplate).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx, serverTemplate, nil
		}
		return ctx, serverTemplate, err
	}

	return ctx, serverTemplate, nil
}

func (database *TblServerTemplateImpl) PaginateServerTemplate(ctx context.Context, page int, size int) (context.Context, []entity.TblServerTemplate, error) {
	var serverTemplates []entity.TblServerTemplate
	offset := (page - 1) * size
	err := database.DB.WithContext(ctx).Limit(size).Offset(offset).Find(&serverTemplates).Error
	if err != nil {
		return ctx, nil, err
	}

	return ctx, serverTemplates, nil
}

func (database *TblServerTemplateImpl) DeleteServerTemplate(ctx context.Context, key string, value string) (context.Context, error) {
	err := database.DB.WithContext(ctx).Where(key+" = ?", value).Delete(&entity.TblServerTemplate{}).Error
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}
