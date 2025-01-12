package tblscript

import (
	"context"
	"cynxhostagent/internal/model/entity"
	"cynxhostagent/internal/repository/database"

	"gorm.io/gorm"
)

type TblScriptImpl struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) database.TblScript {
	return &TblScriptImpl{
		DB: DB,
	}
}

func (database *TblScriptImpl) CreateScript(ctx context.Context, script entity.TblScript) (context.Context, int, error) {
	err := database.DB.WithContext(ctx).Create(&script).Error
	if err != nil {
		return ctx, 0, err
	}

	return ctx, script.Id, nil
}

func (database *TblScriptImpl) GetScript(ctx context.Context, key string, value string) (context.Context, entity.TblScript, error) {
	var script entity.TblScript
	err := database.DB.WithContext(ctx).Where(key+" = ?", value).First(&script).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx, script, nil
		}
		return ctx, script, err
	}

	return ctx, script, nil
}

func (database *TblScriptImpl) DeleteScript(ctx context.Context, key string, value string) (context.Context, error) {
	err := database.DB.WithContext(ctx).Where(key+" = ?", value).Delete(&entity.TblScript{}).Error
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}
