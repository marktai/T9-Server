package main

import (
	"db"
	"flag"
	"server"
)

func main() {
	var port int
	var disableAuth bool

	flag.IntVar(&port, "port", 8081, "Port the server listens to")
	flag.BoolVar(&disableAuth, "disableAuth", false, "Disables authentication requirements")

	flag.Parse()

	db.Open()
	defer db.Db.Close()
	server.Run(port, disableAuth)
}
