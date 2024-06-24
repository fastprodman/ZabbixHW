package dbrepo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type FileDB struct {
	data      []map[string]interface{}
	file      *os.File
	fileMutex *sync.RWMutex
}

func NewFileDB(file *os.File) (*FileDB, error) {
	fileMutex := &sync.RWMutex{}
	fileMutex.Lock()
	data, err := readJSONFile(file)
	fileMutex.Unlock()
	if err != nil {
		return nil, err
	}

	db := &FileDB{
		data:      data,
		file:      file,
		fileMutex: fileMutex,
	}
	return db, nil
}

func (db *FileDB) CreateRecord(data map[string]interface{}) error {
	db.fileMutex.Lock()
	defer db.fileMutex.Unlock()

	var newID uint32 = 1 // Default ID if the database is empty

	// Check if the database is not empty
	if len(db.data) > 0 {
		// Get the ID of the last record
		lastRecord := db.data[len(db.data)-1]
		if id, ok := lastRecord["id"].(uint32); ok {
			newID = id + 1
		} else {
			return errors.New("invalid ID type in last record")
		}
	}

	// Set the new record's ID
	data["id"] = newID

	// Add the new record to the database
	db.data = append(db.data, data)

	err := rewriteJSONFile(db.file, db.data)
	if err != nil {
		return fmt.Errorf("error writing to file: %s", err)
	}

	return nil
}

func (db *FileDB) ReadRecord(id uint32) (map[string]interface{}, error) {
	db.fileMutex.RLock()
	defer db.fileMutex.RUnlock()

	// Iterate through the records to find the matching ID
	for _, record := range db.data {
		if recordID, ok := record["id"].(uint32); ok {
			if recordID == id {
				return record, nil
			}
		} else {
			return nil, errors.New("invalid ID type in record")
		}
	}
	// If no matching record is found, return an error
	return nil, ErrRecordNotFound
}

func (db *FileDB) UpdateRecord(id uint32, data map[string]interface{}) error {
	db.fileMutex.Lock()
	defer db.fileMutex.Unlock()

	// Iterate through the records to find the matching ID
	data["id"] = id
	for i, record := range db.data {
		if recordID, ok := record["id"].(uint32); ok {
			if recordID == id {
				// Update the record with new data, preserving the ID
				for key, value := range data {
					record[key] = value
				}
				db.data[i] = record
				err := rewriteJSONFile(db.file, db.data)
				if err != nil {
					return err
				}
				return nil
			}
		} else {
			return errors.New("invalid ID type in record")
		}
	}
	// If no matching record is found, return an error
	return ErrRecordNotFound
}

// DeleteRecord deletes a record by its ID
func (db *FileDB) DeleteRecord(id uint32) error {
	db.fileMutex.Lock()
	defer db.fileMutex.Unlock()
	// Iterate through the records to find the matching ID
	for i, record := range db.data {
		if recordID, ok := record["id"].(uint32); ok {
			if recordID == id {
				// Remove the record from the slice
				db.data = append(db.data[:i], db.data[i+1:]...)
				err := rewriteJSONFile(db.file, db.data)
				if err != nil {
					return err
				}
				return nil
			}
		} else {
			return errors.New("invalid ID type in record")
		}
	}
	// If no matching record is found, return an error
	return ErrRecordNotFound
}

// rewriteJSONFile truncates the given file and writes the provided data to it
func rewriteJSONFile(file *os.File, data []map[string]interface{}) error {
	// Truncate the file
	err := file.Truncate(0)
	if err != nil {
		return err
	}

	// Move the file pointer to the beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	// Create a JSON encoder and encode the data
	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

// readJSONFile reads an array of JSON objects from the given file and returns it as []map[string]interface{}
func readJSONFile(file *os.File) ([]map[string]interface{}, error) {
	// Read the entire file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// If the file is empty, return an empty slice
	if len(fileContent) == 0 {
		return []map[string]interface{}{}, nil
	}

	// Parse the JSON content
	var data []map[string]interface{}
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
