package dbrepo

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type FileDB struct {
	file      *os.File
	fileMutex *sync.RWMutex
}

func NewFileDB(file *os.File) (*FileDB, error) {
	db := &FileDB{
		file:      file,
		fileMutex: &sync.RWMutex{},
	}
	return db, nil
}

func (db *FileDB) CreateRecord(data map[string]interface{}) (map[string]interface{}, error) {
	db.fileMutex.Lock()
	defer db.fileMutex.Unlock()

	var id uint32

	// Read the file to determine the last ID used
	scanner := bufio.NewScanner(db.file)
	for scanner.Scan() {
		var record map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, err
		}
		if recordID, ok := record["id"].(float64); ok {
			id = uint32(recordID)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Increment ID for new record
	data["id"] = id + 1

	// Marshal the new record to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Append newline character to JSON data
	recordString := string(jsonData) + "\n"

	// Write the new record string to the file
	if _, err := db.file.WriteString(recordString); err != nil {
		return nil, err
	}

	return data, nil
}

func (db *FileDB) ReadRecord(id uint32) (map[string]interface{}, error) {
	db.fileMutex.RLock()
	defer db.fileMutex.RUnlock()

	// Rewind the file to the beginning
	if _, err := db.file.Seek(0, 0); err != nil {
		return nil, err
	}

	// Read the file line by line
	scanner := bufio.NewScanner(db.file)
	for scanner.Scan() {
		var record map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, err
		}
		if recordID, ok := record["id"].(float64); ok && uint32(recordID) == id {
			return record, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, ErrRecordNotFound
}

func (db *FileDB) UpdateRecord() {}
func (db *FileDB) DeleteRecord() {}
