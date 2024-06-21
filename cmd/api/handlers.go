package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) postRecordHandler(w http.ResponseWriter, r *http.Request) {

	// Decode the JSON body into a map
	var data map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Check if the JSON contains the field "id"
	if _, exists := data["id"]; exists {
		http.Error(w, "Field 'id' is not allowed", http.StatusBadRequest)
		return
	}

	respData, err := app.DB.CreateRecord(data)
	if err != nil {
		http.Error(w, "Error creating record", http.StatusInternalServerError)
		return
	}

	// Respond back with the updated data
	response, err := json.Marshal(respData)
	if err != nil {
		http.Error(w, "Error encoding response JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (app *application) getRecordHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) putRecordHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deleteRecordHandler(w http.ResponseWriter, r *http.Request) {

}
