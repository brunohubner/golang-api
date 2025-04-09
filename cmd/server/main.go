package main

import (
	"fmt"
	"net/http"

	"app/configs"
	"app/internal/entity"
	"app/internal/infra/database"
	"app/internal/infra/webserver/handlers"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "app/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Bruno Hubner co-chi API
// @version         1.0
// @description     lorem impsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
// @termsOfService  http://swagger.io/terms

// @contact.name   Bruno Hubner
// @contact.url    https://brunohubner.com
// @contact.email  brunohubnerdev@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8003
// @BasePath  /

// @securityDefinitions.apiKey ApiKeyAuth
// @in          	header
// @name        	Authorization
// @description  	Authorization header with JWT Bearer token

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
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
	r.Use(middleware.Recoverer)

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

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8003/docs/doc.json")))

	http.ListenAndServe(fmt.Sprintf(":%s", config.WebServerPort), r)
}
