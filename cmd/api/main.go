package main

import (
	"ecommerce-api/internal/platform/cache"
	"ecommerce-api/internal/platform/config"
	"ecommerce-api/internal/platform/database"
	"ecommerce-api/internal/platform/middleware"
	"ecommerce-api/internal/platform/response"
	"fmt"
	"log"
	"net/http"
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

	mux := http.NewServeMux()

	handler := middleware.Logger(mux)
	handler = middleware.Recovery(handler)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, "service is healthy", map[string]string{
			"status": "ok",
		})
	})

	address := fmt.Sprintf(":%s", cfg.Server.Port)

	log.Printf("server running on %s", address)

	err = http.ListenAndServe(address, handler)
	if err != nil {
		log.Fatal(err)
	}
}
