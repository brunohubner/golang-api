package entity_test

import (
	"testing"

	"github.com/brunohubner/golang-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

var (
	fakeUserName     = "John Doe"
	fakeUserEmail    = "john@doe.com"
	fakeUserPassword = "@Pass123"
)

func TestNewUser(t *testing.T) {
	user, err := entity.NewUser(fakeUserName, fakeUserEmail, fakeUserPassword)

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, fakeUserName, user.Name)
	assert.Equal(t, fakeUserEmail, user.Email)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
}

func TestUser_ValidatePassword(t *testing.T) {
	user, _ := entity.NewUser(fakeUserName, fakeUserEmail, fakeUserPassword)

	assert.True(t, user.ValidatePassword(fakeUserPassword))
	assert.False(t, user.ValidatePassword("wrongPassword"))
	assert.NotEqual(t, fakeUserPassword, user.Password)
}
