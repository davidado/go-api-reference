// Package api : API server
package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/davidado/go-api-reference/service/cart"
	"github.com/davidado/go-api-reference/service/order"
	"github.com/davidado/go-api-reference/service/product"
	"github.com/davidado/go-api-reference/service/user"
	"github.com/gorilla/mux"
)

// Server is the main struct for the API server
type Server struct {
	addr string
	db   *sql.DB
}

// NewServer creates a new APIServer instance
func NewServer(addr string, db *sql.DB) *Server {
	return &Server{addr: addr, db: db}
}

// Run starts the API server
func (s *Server) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	orderStore := order.NewStore(s.db)

	cartHandler := cart.NewHandler(orderStore, productStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
