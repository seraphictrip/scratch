package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"scratch/service/products"
	"scratch/service/users"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run(ctx context.Context) error {
	router := mux.NewRouter()
	// version our routes to protect from breaking changes
	subrouter := router.PathPrefix("/api/v1/").Subrouter()

	userStore := users.NewStore(s.db)
	usersHandler := users.NewHandler(userStore)
	usersHandler.RegisterRoutes(subrouter)

	productsStore := products.NewStore(s.db)
	productsHandler := products.NewHandler(productsStore)
	productsHandler.RegisterRoutes(subrouter)

	log.Println("listening on", s.addr)
	return http.ListenAndServe(s.addr, subrouter)
}
