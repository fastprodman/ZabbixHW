package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"zabbixhw/pkg/repository"
	"zabbixhw/pkg/repository/filedb"
)

type application struct {
	DB repository.DatabaseRepo
}

func main() {
	// Define the flags with default values
	filepath := flag.String("filepath", "./testdata/db.json", "Path to the file")
	port := flag.Int("port", 8080, "Port number")

	// Parse the flags
	flag.Parse()

	file, err := os.OpenFile(*filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	db, err := filedb.NewFileDB(file)
	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		DB: db,
	}

	addr := fmt.Sprintf(":%d", *port)
	err = http.ListenAndServe(addr, app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
