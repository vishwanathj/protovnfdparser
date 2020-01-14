package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/vishwanathj/protovnfdparser/pkg/dataaccess"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/vishwanathj/protovnfdparser/pkg/models"
)

// VnfdBasePath the base URI path for REST operations
//const VnfdBasePath = "/vnfds"
const VnfdBasePath = "/"

// WebServerPort the port on which the web server runs
const WebServerPort = "8080"

// WebServerSecurePort the port on which the secure web server runs
const WebServerSecurePort = "443"

// Server struct to hold the mux router object
type Server struct {
	router *mux.Router
}

// NewServer creates a new mux server
func NewServer(u models.VnfdService, dal dataaccess.DataAccessLayer) *Server {
	log.Debug()
	s := Server{router: mux.NewRouter()}
	NewVnfdRouter(u, dal, s.newSubrouter(VnfdBasePath))
	return &s
}

// Start starts the web server
func (s *Server) Start() {
	log.Debug()
	fmt.Println("Listening on port: ", WebServerPort)
	go http.ListenAndServeTLS(":"+WebServerSecurePort, "/etc/ssl/certs/vnfdsvc.crt",
		"/etc/ssl/certs/vnfdsvc.key", handlers.LoggingHandler(os.Stdout, s.router))
	if err := http.ListenAndServe(":"+WebServerPort, handlers.LoggingHandler(os.Stdout, s.router)); err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}
}

func (s *Server) newSubrouter(path string) *mux.Router {
	log.Debug()
	return s.router.PathPrefix(path).Subrouter()
}
