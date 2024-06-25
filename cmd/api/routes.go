package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /records", app.postRecordHandler)
	mux.HandleFunc("GET /records/{id}", app.getRecordHandler)
	mux.HandleFunc("PUT /records/{id}", app.putRecordHandler)
	mux.HandleFunc("DELETE /records/{id}", app.deleteRecordHandler)

	return mux
}
