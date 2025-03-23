package database

import (
	"github.com/brunohubner/golang-api/internal/entity"
	"gorm.io/gorm"
)

type Product struct {
	DB *gorm.DB
}

func NewProduct(db *gorm.DB) *Product {
	return &Product{DB: db}
}

func (p *Product) Create(product *entity.Product) error {
	return p.DB.Create(product).Error
}

func (p *Product) FindAll(page, limit int, sort string) ([]entity.Product, error) {
	var products []entity.Product
	var err error

	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if limit != 0 && page != 0 {
		err = p.DB.Offset((page - 1) * limit).Limit(limit).Order("created_at" + sort).Find(&products).Error
	} else {
		err = p.DB.Order("created_at" + sort).Find(&products).Error
	}

	return products, err
}

func (p *Product) FindById(id string) (*entity.Product, error) {
	product := &entity.Product{}
	err := p.DB.First(product, "id = ?", id).Error
	return product, err
}

func (p *Product) Update(product *entity.Product) error {
	product, err := p.FindById(product.ID.String())
	if err != nil {
		return err
	}

	return p.DB.Save(product).Error
}

func (p *Product) Delete(id string) error {
	product, err := p.FindById(id)
	if err != nil {
		return err
	}

	return p.DB.Delete(product).Error
}
