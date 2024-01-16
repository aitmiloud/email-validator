package main

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/unrolled/render"
)

var renderer = render.New()

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	if err := startServer(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func startServer() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Mount("/email", EmailRoutes())

	return http.ListenAndServe(":3000", r)
}

func EmailRoutes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/check", checkEmail)
	return r
}

func checkEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	domain, err := getDomainFromEmail(email)
	if err != nil {
		renderer.JSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	domainChecks, err := checkDomain(domain)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Domain checks: %+v\n", domainChecks)
	renderer.JSON(w, http.StatusOK, domainChecks)
}

func getDomainFromEmail(email string) (string, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", errors.New("Invalid email address format")
	}
	return parts[1], nil
}
