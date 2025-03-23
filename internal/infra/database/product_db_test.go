package database_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/brunohubner/golang-api/internal/entity"
	"github.com/brunohubner/golang-api/internal/infra/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func makeInMemoryProductDB(t *testing.T) *database.Product {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})

	return database.NewProduct(db)
}

func TestNewProduct(t *testing.T) {
	productDB := makeInMemoryProductDB(t)

	p, err := entity.NewProduct("Product 1", 10.5)
	assert.Nil(t, err)
	err = productDB.Create(p)
	assert.Nil(t, err)
	assert.NotNil(t, p.ID)
	assert.NotEmpty(t, p.ID)
	assert.Equal(t, "Product 1", p.Name)
	assert.Equal(t, 10.5, p.Price)
	assert.NotZero(t, p.CreatedAt)
}

func TestFindAllProducts(t *testing.T) {
	productDB := makeInMemoryProductDB(t)

	for i := 1; i < 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*1000)
		assert.NoError(t, err)
		productDB.Create(product)
	}

	products, err := productDB.FindAll(1, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = productDB.FindAll(2, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)
}
