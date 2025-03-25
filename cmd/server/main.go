package main

import (
	"net/http"

	"github.com/brunohubner/golang-api/configs"
	"github.com/brunohubner/golang-api/internal/entity"
	"github.com/brunohubner/golang-api/internal/infra/database"
	"github.com/brunohubner/golang-api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	_, err := configs.LoadConfig("../")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("sqlite.test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entity.User{}, &entity.Product{})

	productDB := database.NewProduct(db)
	productHandler := handlers.NeProductHandler(productDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	prefix := "/api/v1"

	r.Route(prefix, func(r chi.Router) {
		r.Post("/products", productHandler.CreateProduct)
		r.Get("/products/{id}", productHandler.GetProduct)
		r.Put("/products/{id}", productHandler.UpdateProduct)
		r.Delete("/products/{id}", productHandler.DeleteProduct)
	})

	http.ListenAndServe(":8001", r)
}
