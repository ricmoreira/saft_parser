package main

import (
	"saft_parser/containers"
	"saft_parser/server"
)

func main() {

	container := containers.BuildContainer()

	err := container.Invoke(func(server *server.Server) {
		server.Run()
	})

	if err != nil {
		panic(err)
	}
}
