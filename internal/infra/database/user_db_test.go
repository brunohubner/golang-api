package database_test

import (
	"testing"
	// "github.com/brunohubner/

	"github.com/brunohubner/golang-api/internal/entity"
	"github.com/brunohubner/golang-api/internal/infra/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	fakeUserName  = "John Doe"
	fakeUserEmail = "john@doe.com"
	fakeUserPass  = "@Pass123"
)

func TestCreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.User{})

	user, _ := entity.NewUser(fakeUserName, fakeUserEmail, fakeUserPass)
	userDB := database.NewUser(db)

	err = userDB.Create(user)
	assert.Nil(t, err)

	userFound := entity.User{}
	err = db.First(&userFound, "id = ?", user.ID).Error

	assert.Nil(t, err)

	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
	assert.Equal(t, user.Password, userFound.Password)
}
