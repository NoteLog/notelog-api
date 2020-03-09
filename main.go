package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	log.Printf("Starting API Server")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// esURL := flag.String("esurl", "", "Elasticsearch service URL")
	// esPWD := flag.String("espwd", "", "Elasticsearch service password")

	flag.Parse()

	r := chi.NewRouter()

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second),
	)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to notelog-api root."))
	})

	//Routes for GitHub Repos
	r.Route("/github", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Get("/", GetGitHubRepo)
		})
	})
	r.Route("/blog", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Get("/", GetBlog)
		})
	})

	log.Printf("Listening On Port: %s", port)

	// Start the HTTP Server
	http.ListenAndServe(":"+port, r)
}
