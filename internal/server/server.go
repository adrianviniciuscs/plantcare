package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"plantcare/internal/endpoints"
)

func Start() {
	r := chi.NewRouter()

	endpoints.SetupRoutes(r)

	fmt.Printf("Server running in 8000")
	http.ListenAndServe(":8000", r)

}
