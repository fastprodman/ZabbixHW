package main

import (
	"log"
	"net/http"
)

func main() {

	err := http.ListenAndServe(":8080", routes())
	if err != nil {
		log.Fatal(err)
	}
}
