package services

import (
	"avito-task/entity"
	"avito-task/initializers"
	"fmt"

	"gorm.io/gorm"
)

var Users = UserService{initializers.DB}

type UserService struct {
	DB *gorm.DB
}

func CreateUserService(db *gorm.DB) UserService {
	return UserService{db}
}

func (u UserService) GetAll() ([]*entity.User, error) {
	var users []*entity.User
	results := u.DB.Preload("Segments").Find(&users)
	if results.Error != nil {
		return users, results.Error
	}
	return users, nil
}

func (u UserService) GetOne(id uint) (*entity.User, error) {
	var user *entity.User
	result := u.DB.Preload("Segments").First(&user, "id = ?", fmt.Sprint(id))
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (u UserService) AddOne(id uint) (*entity.User, error) {
	user := &entity.User{ID: id}
	result := u.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
