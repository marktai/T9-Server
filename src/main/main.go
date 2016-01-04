package main

import (
	"flag"
	"game"
	"server"
)

func main() {
	var port int

	flag.IntVar(&port, "Port", 8081, "Port the server listens to")

	flag.Parse()

	game.Open()
	defer game.Close()
	server.Run(port)
}
