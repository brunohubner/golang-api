package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/brunohubner/golang-api/internal/dto"
	"github.com/brunohubner/golang-api/internal/entity"
	"github.com/brunohubner/golang-api/internal/infra/database"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NeProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
