package filedb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

// Custom error messages for the package
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidIDType  = errors.New("invalid ID type in record")
)

// FileDB struct that represents the file-based database
type FileDB struct {
	data      []map[string]interface{} // In-memory data storage
	file      *os.File                 // File handler for the database file
	fileMutex *sync.Mutex              // Mutex for handling concurrent access to the file
}

// NewFileDB initializes a new FileDB instance and loads data from the provided file
func NewFileDB(file *os.File) (*FileDB, error) {
	fileMutex := &sync.Mutex{}
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// Read initial data from the JSON file
	data, err := readJSONFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %w", err)
	}

	// Create a new FileDB instance
	db := &FileDB{
		data:      data,
		file:      file,
		fileMutex: fileMutex,
	}
	return db, nil
}

// CreateRecord adds a new record to the database
func (db *FileDB) CreateRecord(data map[string]interface{}) error {
	db.fileMutex.Lock()
	defer db.fileMutex.Unlock()

	// Determine the ID for the new record
	var newID uint32 = 1 // Default ID if the database is empty

	if len(db.data) > 0 {
		lastRecord := db.data[len(db.data)-1]
		if id, ok := lastRecord["id"].(float64); ok {
			newID = uint32(id) + 1
		} else {
			return ErrInvalidIDType
		}
	}

	// Set the new record's ID and add it to the database
	data["id"] = newID
	db.data = append(db.data, data)

	// Write updated data back to the file
	if err := rewriteJSONFile(db.file, db.data); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

// ReadRecord retrieves a record by its ID
func (db *FileDB) ReadRecord(id uint32) (map[string]interface{}, error) {
	db.fileMutex.Lock()
	defer db.fileMutex.Unlock()

	// Search for the record with the specified ID
	for _, record := range db.data {
		if recordID, ok := record["id"].(float64); ok {
			if uint32(recordID) == id {
				return record, nil
			}
		} else {
			return nil, ErrInvalidIDType
		}
	}

	// Return an error if the record is not found
	return nil, ErrRecordNotFound
}

// UpdateRecord updates a record with the specified ID
func (db *FileDB) UpdateRecord(id uint32, data map[string]interface{}) error {
	db.fileMutex.Lock()
	defer db.fileMutex.Unlock()

	// Search for the record with the specified ID and update it
	for i, record := range db.data {
		if recordID, ok := record["id"].(float64); ok {
			if uint32(recordID) == id {
				for key, value := range data {
					record[key] = value
				}
				record["id"] = id
				db.data[i] = record

				// Write updated data back to the file
				if err := rewriteJSONFile(db.file, db.data); err != nil {
					return fmt.Errorf("error writing to file: %w", err)
				}
				return nil
			}
		} else {
			return ErrInvalidIDType
		}
	}

	// Return an error if the record is not found
	return ErrRecordNotFound
}

// DeleteRecord removes a record with the specified ID
func (db *FileDB) DeleteRecord(id uint32) error {
	db.fileMutex.Lock()
	defer db.fileMutex.Unlock()

	// Search for the record with the specified ID and delete it
	for i, record := range db.data {
		if recordID, ok := record["id"].(float64); ok {
			if uint32(recordID) == id {
				db.data = append(db.data[:i], db.data[i+1:]...)

				// Write updated data back to the file
				if err := rewriteJSONFile(db.file, db.data); err != nil {
					return fmt.Errorf("error writing to file: %w", err)
				}
				return nil
			}
		} else {
			return ErrInvalidIDType
		}
	}

	// Return an error if the record is not found
	return ErrRecordNotFound
}

// rewriteJSONFile truncates the file and writes the provided data to it
func rewriteJSONFile(file *os.File, data []map[string]interface{}) error {
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("error truncating file: %w", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("error seeking file: %w", err)
	}

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("error encoding JSON data: %w", err)
	}

	return nil
}

// readJSONFile reads JSON data from the file and returns it as a slice of maps
func readJSONFile(file *os.File) ([]map[string]interface{}, error) {
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if len(fileContent) == 0 {
		return []map[string]interface{}{}, nil
	}

	var data []map[string]interface{}
	if err := json.Unmarshal(fileContent, &data); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON data: %w", err)
	}

	return data, nil
}
