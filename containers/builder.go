// Inspiration to create dependcy injection came from this post: https://blog.drewolson.org/dependency-injection-in-go/

package containers

import (
	"saft_parser/config"
	controllers "saft_parser/controllers/v1"
	"saft_parser/handlers"
	"saft_parser/middleware"
	"saft_parser/server"
	"saft_parser/services"

	"go.uber.org/dig"
)

// BuildContainer returns a container with all app dependencies built in
func BuildContainer() *dig.Container {
	container := dig.New()

	// config
	container.Provide(config.NewConfig)

	// persistance layer

	// services
	container.Provide(services.NewSAFTParser)
	container.Provide(services.NewSAFTService)

	// controllers
	container.Provide(controllers.NewSAFTController)

	// generic http layer
	container.Provide(middleware.NewMiddleware)
	container.Provide(handlers.NewHttpHandlers)

	// server
	container.Provide(server.NewServer)

	return container
}
