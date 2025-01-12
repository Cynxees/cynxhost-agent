package tbluser

import (
	"context"
	"cynxhostagent/internal/model/entity"
	"cynxhostagent/internal/model/request"
	"cynxhostagent/internal/repository/database"

	"gorm.io/gorm"
)

type TblUserImpl struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) database.TblUser {
	return &TblUserImpl{
		DB: DB,
	}
}

func (database *TblUserImpl) InsertUser(ctx context.Context, user entity.TblUser) (context.Context, int, error) {
	err := database.DB.WithContext(ctx).Create(&user).Error
	if err != nil {
		return ctx, 0, err
	}

	return ctx, user.Id, nil
}

func (database *TblUserImpl) CheckUserExists(ctx context.Context, key string, value string) (context.Context, bool, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&entity.TblUser{}).Where(key+" = ?", value).Count(&count).Error
	if err != nil {
		return ctx, false, err
	}

	return ctx, count > 0, nil
}

func (database *TblUserImpl) GetUser(ctx context.Context, key string, value string) (context.Context, *entity.TblUser, error) {
	var user entity.TblUser
	err := database.DB.WithContext(ctx).Where(key+" = ?", value).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx, nil, nil
		}
		return ctx, &user, err
	}

	return ctx, &user, nil
}

func (database *TblUserImpl) PaginateUser(ctx context.Context, req request.PaginateRequest) (context.Context, []entity.TblUser, error) {
	var users []entity.TblUser
	offset := (req.Page - 1) * req.Size
	err := database.DB.WithContext(ctx).Limit(req.Size).Offset(offset).Find(&users).Error
	if err != nil {
		return ctx, nil, err
	}

	return ctx, users, nil
}
