package server

import (
	"fmt"
	"net/http"
    "plantcare/internal/endpoints"
	"github.com/go-chi/chi/v5"
)

func Start() {
    r := chi.NewRouter()

    endpoints.SetupRoutes(r)


    fmt.Printf("Server running in 8000")
    http.ListenAndServe(":8000", r)

}
