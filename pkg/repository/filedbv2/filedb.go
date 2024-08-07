package filedbv2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Custom error messages for the package
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidIDType  = errors.New("invalid ID type in record")
)

const (
	MaxCachedUpdates = 5
	SyncInterval     = 5 * time.Second
)

// FileDB struct that represents the file-based database
type FileDB struct {
	data          []map[string]interface{} // In-memory data storage
	file          *os.File                 // File handler for the database file
	fileMutex     *sync.RWMutex            // Mutex for handling concurrent access to the file
	dataMutex     *sync.RWMutex            // Mutex for handling concurrent access to in-memory data
	cachedUpdates uint
	updateChan    chan bool
	doneChan      chan bool
}

// NewFileDB initializes a new FileDB instance and loads data from the provided file
func NewFileDB(filePath string) (*FileDB, error) {
	fileMutex := &sync.RWMutex{}
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	// Read initial data from the JSON file
	data, err := readJSONFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %w", err)
	}

	// Create a new FileDB instance
	db := &FileDB{
		data:          data,
		file:          file,
		fileMutex:     fileMutex,
		dataMutex:     &sync.RWMutex{},
		updateChan:    make(chan bool, 1), // Buffered channel to prevent blocking
		doneChan:      make(chan bool),
		cachedUpdates: 0,
	}

	go db.syncLoop()

	return db, nil
}

func (db *FileDB) syncLoop() {
	for {
		select {
		case <-db.updateChan:
			db.syncDBWithCache()
		case <-time.After(SyncInterval):
			db.syncDBWithCache()
		case <-db.doneChan:
			db.syncDBWithCache()
			close(db.doneChan)
			return
		}
	}
}

func (db *FileDB) Close() {
	close(db.updateChan)
	// Send close command
	db.doneChan <- true
	// Wait until changes are in sync and than close the file
	<-db.doneChan
	db.file.Close()
}

func (db *FileDB) syncDBWithCache() error {
	db.fileMutex.Lock()
	defer db.fileMutex.Unlock()

	db.dataMutex.RLock()
	defer db.dataMutex.RUnlock()

	// Write updated data back to the file
	if err := rewriteJSONFile(db.file, db.data); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	db.cachedUpdates = 0
	return nil
}

// CreateRecord adds a new record to the database
func (db *FileDB) CreateRecord(data map[string]interface{}) error {
	db.dataMutex.Lock()
	defer db.dataMutex.Unlock()

	// Determine the ID for the new record
	var newID float64 = 1 // Default ID if the database is empty

	if len(db.data) > 0 {
		lastRecord := db.data[len(db.data)-1]
		if id, ok := lastRecord["id"].(float64); ok {
			newID = float64(id + 1)
		} else {
			return ErrInvalidIDType
		}
	}

	// Set the new record's ID and add it to the database
	data["id"] = newID
	db.data = append(db.data, data)

	db.cachedUpdates++
	if db.cachedUpdates > MaxCachedUpdates {
		select {
		case db.updateChan <- true:
		default:
		}
	}

	return nil
}

// ReadRecord retrieves a record by its ID
func (db *FileDB) ReadRecord(id uint32) (map[string]interface{}, error) {
	db.dataMutex.RLock()
	defer db.dataMutex.RUnlock()

	// Search for the record with the specified ID
	for _, record := range db.data {
		if recordID, ok := record["id"].(float64); ok {
			if recordID == float64(id) {
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
	db.dataMutex.Lock()
	defer db.dataMutex.Unlock()

	// Ensure the data has the ID field
	data["id"] = float64(id)

	// Search for the record with the specified ID and update it
	for i, record := range db.data {
		if recordID, ok := record["id"].(float64); ok {
			if recordID == float64(id) {
				// Replace the old record with the new data
				db.data[i] = data
				db.cachedUpdates++
				if db.cachedUpdates > MaxCachedUpdates {
					select {
					case db.updateChan <- true:
					default:
					}
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
	db.dataMutex.Lock()
	defer db.dataMutex.Unlock()

	// Search for the record with the specified ID and delete it
	for i, record := range db.data {
		if recordID, ok := record["id"].(float64); ok {
			if recordID == float64(id) {
				db.data = append(db.data[:i], db.data[i+1:]...)
				db.cachedUpdates++
				if db.cachedUpdates > MaxCachedUpdates {
					select {
					case db.updateChan <- true:
					default:
					}
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
