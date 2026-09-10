package main

import (
	"ecommerce-api/internal/platform/config"
	"fmt"
	"log"
	"net/http"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

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
