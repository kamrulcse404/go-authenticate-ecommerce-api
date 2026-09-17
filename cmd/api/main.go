package main

import (
	"ecommerce-api/internal/auth"
	"ecommerce-api/internal/platform/cache"
	"ecommerce-api/internal/platform/config"
	"ecommerce-api/internal/platform/database"
	"ecommerce-api/internal/platform/middleware"
	"ecommerce-api/internal/platform/response"
	"ecommerce-api/internal/user"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("postgres connected successfully")

	redisClient, err := cache.Connect(cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}

	defer redisClient.Close()

	log.Println("redis connected successfully")

	r := chi.NewRouter()

	r.Use(middleware.Recovery)
	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(
			w,
			http.StatusOK,
			"service is healthy",
			map[string]string{
				"status": "ok",
			},
		)
	})

	// user mount
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	r.Mount("/users", user.Routes(userHandler))

	//auth mount
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	r.Mount("/auth", auth.Routes(authHandler))

	address := fmt.Sprintf(":%s", cfg.Server.Port)

	log.Printf("server running on %s", address)

	err = http.ListenAndServe(address, r)
	if err != nil {
		log.Fatal(err)
	}
}
