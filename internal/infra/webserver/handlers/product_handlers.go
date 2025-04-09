package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"app/internal/dto"
	"app/internal/entity"
	"app/internal/infra/database"
	utils "app/pkg/entity"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NeProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// Create Product godoc
// @Summary      Create product
// @Description  Create products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request     body      dto.CreateProductInput  true  "product request"
// @Success      201
// @Failure      500         {object}  Error
// @Router       /api/v1/products [post]
// @Security ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: "Invalid request"}
		json.NewEncoder(w).Encode(err)
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: "Invalid product data"}
		json.NewEncoder(w).Encode(err)
		return
	}

	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := Error{Message: "Internal Server Error"}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetProduct godoc
// @Summary      Get a product
// @Description  Get a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "product ID" Format(uuid)
// @Success      200  {object}  entity.Product
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/products/{id} [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !utils.IsUUID(id) {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: "Invalid UUID"}
		json.NewEncoder(w).Encode(err)
		return
	}

	product, err := h.ProductDB.FindById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := Error{Message: "Product not found"}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Update a product
// @Tags         productsa
// @Accept       json
// @Produce      json
// @Param        id        	path      string                  true  "product ID" Format(uuid)
// @Param        request     body      dto.UpdateProductInput  true  "product request"
// @Success      200
// @Failure      400			{object}  Error
// @Failure      404			{object}  Error
// @Failure      500      {object}  Error
// @Router       /api/v1/products/{id} [put]
// @Security ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !utils.IsUUID(id) {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: "Invalid UUID"}
		json.NewEncoder(w).Encode(err)
		return
	}

	var product dto.UpdateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: "Invalid request"}
		json.NewEncoder(w).Encode(err)
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: "Invalid product data"}
		json.NewEncoder(w).Encode(err)
		return
	}

	p.ID = id

	err = h.ProductDB.Update(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := Error{Message: "Internal Server Error"}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Delete a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id        path      string                  true  "product ID" Format(uuid)
// @Success      200
// @Failure      404 			{object}  Error
// @Failure      500      {object}  Error
// @Router       /api/v1/products/{id} [delete]
// @Security ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !utils.IsUUID(id) {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: "Invalid UUID"}
		json.NewEncoder(w).Encode(err)
		return
	}

	err := h.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := Error{Message: "Internal Server Error"}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// List Products godoc
// @Summary      List products
// @Description  get all products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page      query     string  false  "page number"
// @Param        limit     query     string  false  "limit"
// @Success      200       {array}   entity.Product
// @Failure      404       {object}  Error
// @Failure      500       {object}  Error
// @Router       /api/v1/products [get]
// @Security ApiKeyAuth
func (h *ProductHandler) FindManyProducts(w http.ResponseWriter, r *http.Request) {
	pageString := r.URL.Query().Get("page")
	limitString := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")

	page, err := strconv.Atoi(pageString)
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(limitString)
	if err != nil {
		limit = 10
	}

	products, err := h.ProductDB.FindAll(page, limit, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := Error{Message: "Internal Server Error"}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
