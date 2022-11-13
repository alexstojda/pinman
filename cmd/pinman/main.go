package main

import (
	"os"

	"pinman/web"
)

func main() {
	spaPath := os.Getenv("SPA_PATH")
	server := web.NewServer(spaPath)

	server.StartServer()
}
