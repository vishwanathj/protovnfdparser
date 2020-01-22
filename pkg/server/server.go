package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/vishwanathj/protovnfdparser/pkg/config"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/vishwanathj/protovnfdparser/pkg/models"
)

// Server struct to hold the mux router object
type Server struct {
	router *mux.Router
	config *config.WebServerConfig
}

// NewServer creates a new mux server
func NewServer(v models.VnfdService, cfg *config.WebServerConfig) *Server {
	log.Debug()
	s := Server{router: mux.NewRouter(), config: cfg}
	NewVnfdRouter(v, s.newSubrouter(cfg.WebServerBasePath))
	return &s
}

// Start starts the web server
func (s *Server) Start() {
	log.Debug()
	webServerPort := fmt.Sprintf("%d", s.config.WebServerPort)
	//webServerSecurePort := fmt.Sprintf("%d", s.config.WebServerSecurePort)
	fmt.Println("Listening on port: ", webServerPort)
	//go http.ListenAndServeTLS(":"+webServerSecurePort, "/etc/ssl/certs/vnfdsvc.crt",
	//"/etc/ssl/certs/vnfdsvc.key", handlers.LoggingHandler(os.Stdout, s.router))
	if err := http.ListenAndServe(":"+webServerPort, handlers.LoggingHandler(os.Stdout, s.router)); err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}
}

func (s *Server) newSubrouter(path string) *mux.Router {
	log.Debug()
	return s.router.PathPrefix(path).Subrouter()
}
