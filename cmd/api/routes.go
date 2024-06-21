package main

import "net/http"

func routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /records", postRecordHandler)
	mux.HandleFunc("GET /records/{id}", getRecordHandler)
	mux.HandleFunc("PUT /records/{id}", putRecordHandler)
	mux.HandleFunc("DELETE /records/{id}", deleteRecordHandler)

	return mux
}
