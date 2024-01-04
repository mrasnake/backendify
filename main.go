package main

import (
	"fmt"
	"log"
	"net/http"
)

// main creates the Service layer object, uses it to create the Transport layer object,
// attaches the router to the base endpoint and begins listening and serving port 9000.
func main() {

	service, err := NewService()
	if err != nil {
		panic(err)
	}
	server := NewServer(service)

	http.Handle("/", server.Router)
	fmt.Println("Starting token service at port 9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
	return
}