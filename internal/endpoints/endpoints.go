package endpoints

import (
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"path/filepath"
	"plantcare/internal/api"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Error uploading file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		result, err := api.RecognizePlant(file)
		if err != nil {
			http.Error(w, "Error recognizing plant", http.StatusInternalServerError)
			return
		}

        tmpl, err := template.ParseFiles("templates/upload.html")
		if err != nil {
			http.Error(w, "Error recognizing plant", http.StatusInternalServerError)
			return
		}

		data := struct {
			Result string
		}{
			Result: result,
		}
    
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, data)
		return
	}
}

func SetupRoutes(r *chi.Mux) {
	workDir, _ := filepath.Abs("../../templates/")
	filesDir := http.Dir(workDir)

	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(filesDir)))

	r.Post("/upload", UploadHandler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/upload.html")
	})

}
