package main

import (
	"fmt"
	"net/http"

	"github.com/danielkrainas/tinkersnest/api/client"
	"github.com/danielkrainas/tinkersnest/api/v1"
)

func main() {
	const ENDPOINT = "http://localhost:9240"

	// Create a new client
	c := client.New(ENDPOINT, http.DefaultClient)
	fmt.Printf("created new client to %q\n", ENDPOINT)

	// Check V1 endpoint is good and healthy
	//=====================================
	err := c.Ping()
	if err != nil {
		panic("error sending ping")
	}

	fmt.Println("sent ping")

	// TODO: blog example
}