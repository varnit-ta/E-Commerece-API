package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/varnit-ta/Ecom-API/service/products"
	"github.com/varnit-ta/Ecom-API/service/user"
)

/*
APIServer represents an API server that handles HTTP requests
and interacts with a database.

It includes:
`addr`: The address where the server will listen for incoming requests.
`db`: A pointer to an `sql.DB` instance for database operations.
*/
type APIServer struct {
	addr string
	db   *sql.DB
}

/*
NewAPIServer initializes a new instance of APIServer.

Parameters:
`addr`: The address (host:port) where the server should listen.
`db`: A pointer to an `sql.DB` instance representing the database connection.

Returns:
A pointer to an `APIServer` instance, ready to be used.
*/
func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

/*
 */
func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	productStore := products.NewStore(s.db)
	productHandler := products.NewHandler(productStore, userStore)
	productHandler.RegisterRoutes(subrouter)

	log.Println("Listening on ", s.addr)

	return http.ListenAndServe(s.addr, router)
}
