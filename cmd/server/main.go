package main

import (
	"net/http"

	"github.com/brunohubner/golang-api/configs"
	"github.com/brunohubner/golang-api/internal/entity"
	"github.com/brunohubner/golang-api/internal/infra/database"
	"github.com/brunohubner/golang-api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	config, err := configs.LoadConfig("../")
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

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB, config.TokenAuth, config.JwtExpiresIn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	prefix := "/api/v1"

	r.Route(prefix, func(r chi.Router) {
		r.Route("/products", func(r chi.Router) {
			r.Use(jwtauth.Verifier(config.TokenAuth))
			r.Use(jwtauth.Authenticator)

			r.Post("/", productHandler.CreateProduct)
			r.Get("/", productHandler.FindManyProducts)
			r.Get("/{id}", productHandler.GetProduct)
			r.Put("/{id}", productHandler.UpdateProduct)
			r.Delete("/{id}", productHandler.DeleteProduct)
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Post("/generate-jwt", userHandler.GetJwt)
		})
	})

	http.ListenAndServe(":8001", r)
}
