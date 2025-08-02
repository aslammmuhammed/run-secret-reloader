package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine // embed gin engine
	server      *http.Server
}

// New returns a configured gin engine
func New() *Server {
	// Disable gin's default logger
	gin.SetMode(gin.ReleaseMode)

	// Create a new engine without any middleware
	engine := gin.New()

	// Add recovery middleware
	engine.Use(gin.Recovery())

	return &Server{Engine: engine}
}

// Run starts the HTTP server
func (s *Server) Run(port string) error {
	s.server = &http.Server{
		Addr:    ":" + port,
		Handler: s.Engine,
	}
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
