package server

import (
	"saft_parser/config"
	"saft_parser/controllers/v1"
	"saft_parser/handlers"
	"saft_parser/middleware"

	"github.com/gin-gonic/gin"
)

// Server is the http layer saft resource
type Server struct {
	config         *config.Config
	saftController *controllers.SAFTController
	middleware     *middleware.Middleware
	handlers       *handlers.HttpHandlers
}

// NewServer is the Server constructor
func NewServer(cf *config.Config,
	saftC *controllers.SAFTController,
	mid *middleware.Middleware,
	hand *handlers.HttpHandlers) *Server {

	return &Server{
		config:         cf,
		saftController: saftC,
		middleware:     mid,
		handlers:       hand,
	}
}

// Run loads server with its routes and starts the server
func (s *Server) Run() {
	// Instantiate a new router
	r := gin.Default()

	// cors
	r.Use(*s.middleware.Cors())

	// generic routes
	r.HandleMethodNotAllowed = false
	r.NoRoute(s.handlers.NotFound)

	// SAFT resource
	saftApi := r.Group("/api/v1/saft")
	{
		// Upload a file
		saftApi.POST("/upload", s.saftController.FileToKafkaAction)
	}

	// Fire up the server
	r.Run(s.config.Host)
}
