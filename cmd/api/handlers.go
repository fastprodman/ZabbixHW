package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// postRecordHandler handles the creation of a new record
func (app *application) postRecordHandler(w http.ResponseWriter, r *http.Request) {

	// Decode the JSON body into a map
	var record map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&record)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Check if the JSON contains the field "id"
	if _, exists := record["id"]; exists {
		http.Error(w, "Field 'id' is not allowed", http.StatusBadRequest)
		return
	}

	// Create the record in the database
	err = app.DB.CreateRecord(record)
	if err != nil {
		http.Error(w, "Error creating record", http.StatusInternalServerError)
		return
	}

	// Respond back with the created record
	response, err := json.Marshal(record)
	if err != nil {
		http.Error(w, "Error encoding response JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// getRecordHandler handles fetching a record by ID
func (app *application) getRecordHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Read the record from the database
	record, err := app.DB.ReadRecord(uint32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Respond back with the fetched record
	response, err := json.Marshal(record)
	if err != nil {
		http.Error(w, "Error encoding response JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// putRecordHandler handles updating a record by ID
func (app *application) putRecordHandler(w http.ResponseWriter, r *http.Request) {
	// Getting id to update from URL parameters
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Decode the JSON body into a map
	var record map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&record)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Check if the JSON contains the field "id"
	if _, exists := record["id"]; exists {
		http.Error(w, "Field 'id' is not allowed", http.StatusBadRequest)
		return
	}

	// Update the record in the database
	err = app.DB.UpdateRecord(uint32(id), record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Respond back with the updated record
	response, err := json.Marshal(record)
	if err != nil {
		http.Error(w, "Error encoding response JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// deleteRecordHandler handles deleting a record by ID
func (app *application) deleteRecordHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Delete the record from the database
	err = app.DB.DeleteRecord(uint32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Respond with no content status
	w.WriteHeader(http.StatusNoContent)
}
