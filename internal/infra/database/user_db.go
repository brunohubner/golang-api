package database

import (
	"app/internal/entity"

	"gorm.io/gorm"
)

type User struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{DB: db}
}

func (u *User) Create(user *entity.User) error {
	return u.DB.Create(user).Error
}

func (u *User) FindByEmail(email string) (*entity.User, error) {
	user := &entity.User{}
	err := u.DB.Where("email = ?", email).First(user).Error
	return user, err
}
