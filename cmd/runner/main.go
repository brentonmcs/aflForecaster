package main

import (
	"log"

	"github.com/brentonmcs/afl/aflServer"
)

func main() {
	log.Print("Starting Server")
	aflServer.StartHTTPServer()
}
