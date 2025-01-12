package tblpersistentnode

import (
	"context"
	"cynxhostagent/internal/model/entity"
	"cynxhostagent/internal/repository/database"

	"gorm.io/gorm"
)

type TblPersistentNodeImpl struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) database.TblPersistentNode {
	return &TblPersistentNodeImpl{
		DB: DB,
	}
}

func (database *TblPersistentNodeImpl) CreatePersistentNode(ctx context.Context, persistentNode entity.TblPersistentNode) (context.Context, int, error) {
	err := database.DB.WithContext(ctx).Create(&persistentNode).Error
	if err != nil {
		return ctx, 0, err
	}

	return ctx, persistentNode.Id, nil
}

func (database *TblPersistentNodeImpl) GetPersistentNodes(ctx context.Context, key string, value string) (context.Context, []entity.TblPersistentNode, error) {
	var persistentNodes []entity.TblPersistentNode
	err := database.DB.WithContext(ctx).Where(key+" = ?", value).Find(&persistentNodes).Error
	if err != nil {
		return ctx, nil, err
	}

	return ctx, persistentNodes, nil
}

func (database *TblPersistentNodeImpl) UpdatePersistentNode(ctx context.Context, id int, persistentNode entity.TblPersistentNode) (context.Context, error) {
	err := database.DB.WithContext(ctx).Model(&entity.TblPersistentNode{}).Where("id = ?", id).Updates(persistentNode).Error
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (database *TblPersistentNodeImpl) DeletePersistentNode(ctx context.Context, id int) (context.Context, error) {
	err := database.DB.WithContext(ctx).Where("id = ?", id).Delete(&entity.TblPersistentNode{}).Error
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}
