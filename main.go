package main

import (
	"log"
)

func main() {
	server := srv.NewHTTPServer()
	log.Fatal(server.ListenAndServe())
}
