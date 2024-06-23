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

	// file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	err := http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
