package entity

import (
	"errors"
	"time"

	"github.com/brunohubner/golang-api/pkg/entity"
)

var (
	ErrIdIsRequired    = errors.New("id is required")
	ErrInvalidId       = errors.New("id is invalid")
	ErrNameIsRequired  = errors.New("name is required")
	ErrPriceIsRequired = errors.New("price is required")
	ErrInvalidIdPrice  = errors.New("price is invalid")
)

type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

func NewProduct(name string, price float64) (*Product, error) {
	p := &Product{
		ID:        entity.NewID().String(),
		Name:      name,
		Price:     price,
		CreatedAt: time.Now(),
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Product) Validate() error {
	if p.ID == "" {
		return ErrIdIsRequired
	}

	if _, err := entity.ParseID(string(p.ID)); err != nil {

	}

	if p.Name == "" {
		return ErrNameIsRequired
	}

	if p.Price == 0 {
		return ErrPriceIsRequired
	}

	if p.Price < 0 {
		return ErrInvalidIdPrice
	}

	return nil
}
