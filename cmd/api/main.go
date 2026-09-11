package main

import (
	"ecommerce-api/internal/platform/cache"
	"ecommerce-api/internal/platform/config"
	"ecommerce-api/internal/platform/database"
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

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	address := fmt.Sprintf(":%s", cfg.Server.Port)

	log.Printf("server running on %s", address)

	err = http.ListenAndServe(address, mux)
	if err != nil {
		log.Fatal(err)
	}
}
