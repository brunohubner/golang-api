package entity_test

import (
	"testing"

	"github.com/brunohubner/golang-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

var (
	fakeProductName  = "Product Test"
	fakeProductPrice = 10.0
)

func TestNewProduct(t *testing.T) {
	p, err := entity.NewProduct(fakeProductName, fakeProductPrice)
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.NotEmpty(t, p.ID)
	assert.Equal(t, fakeProductName, p.Name)
	assert.Equal(t, fakeProductPrice, p.Price)
}

func TestProductWhenNameIsRequired(t *testing.T) {
	p, err := entity.NewProduct("", fakeProductPrice)
	assert.NotNil(t, err)
	assert.Nil(t, p)
	assert.Equal(t, entity.ErrNameIsRequired, err)
}

func TestProductWhenPriceIsRequired(t *testing.T) {
	p, err := entity.NewProduct(fakeProductName, 0)
	assert.NotNil(t, err)
	assert.Nil(t, p)
	assert.Equal(t, entity.ErrPriceIsRequired, err)
}

func TestProductWhenPriceIsInvalid(t *testing.T) {
	p, err := entity.NewProduct(fakeProductName, -1)
	assert.NotNil(t, err)
	assert.Nil(t, p)
	assert.Equal(t, entity.ErrInvalidIdPrice, err)
}
