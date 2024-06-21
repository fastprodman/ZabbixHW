package main

import (
	"log"
	"net/http"
	"zabbixhw/pkg/repository"
)

type application struct {
	DB repository.DatabaseRepo
}

func main() {
	var app application

	err := http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
