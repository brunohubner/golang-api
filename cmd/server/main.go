package main

import (
	"net/http"

	"github.com/brunohubner/golang-api/configs"
	"github.com/brunohubner/golang-api/internal/entity"
	"github.com/brunohubner/golang-api/internal/infra/database"
	"github.com/brunohubner/golang-api/internal/infra/webserver/handlers"
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

	http.HandleFunc("/api/v1/products", productHandler.CreateProduct)

	http.ListenAndServe(":8001", nil)
}
