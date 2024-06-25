package filedb

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"zabbixhw/pkg/helpers"
)

func Test_NewFileDB(t *testing.T) {
	t.Run("Valid JSON File", func(t *testing.T) {
		// Create a temporary file for testing
		tmpFile, err := os.CreateTemp("", "testdb_*.json")
		if err != nil {
			t.Fatalf("failed to create temp file: %s", err)
		}
		defer os.Remove(tmpFile.Name()) // Clean up after the test

		// Write some initial JSON data to the file
		initialData := `[{"id": 1, "name": "Alice", "age": 30}, {"id": 2, "name": "Bob", "age": 25}]`
		if _, err := tmpFile.Write([]byte(initialData)); err != nil {
			t.Fatalf("failed to write initial data to temp file: %s", err)
		}

		// Make sure the file is at the beginning for reading
		if _, err := tmpFile.Seek(0, 0); err != nil {
			t.Fatalf("failed to seek to the beginning of the file: %s", err)
		}

		// Test NewFileDB function
		db, err := NewFileDB(tmpFile)
		if err != nil {
			t.Fatalf("NewFileDB failed: %s", err)
		}

		// Check if the data was read correctly
		expectedData := []map[string]interface{}{
			{"id": float64(1), "name": "Alice", "age": float64(30)}, // JSON unmarshals numbers as float64
			{"id": float64(2), "name": "Bob", "age": float64(25)},
		}

		if len(db.data) != len(expectedData) {
			t.Fatalf("expected %d records, got %d", len(expectedData), len(db.data))
		}

		for i, record := range db.data {
			for key, expectedValue := range expectedData[i] {
				if value, ok := record[key]; !ok || value != expectedValue {
					t.Errorf("expected %v for key %s, got %v", expectedValue, key, value)
				}
			}
		}

		// Ensure the mutex is properly initialized
		if db.fileMutex == nil {
			t.Fatal("expected fileMutex to be initialized")
		}

		// Ensure the file is correctly assigned
		if db.file == nil {
			t.Fatal("expected file to be assigned")
		}
	})

	t.Run("Empty JSON File", func(t *testing.T) {
		// Create a temporary empty file for testing
		tmpFile, err := os.CreateTemp("", "testdb_empty_*.json")
		if err != nil {
			t.Fatalf("failed to create temp file: %s", err)
		}
		defer os.Remove(tmpFile.Name()) // Clean up after the test

		// Test NewFileDB function with an empty file
		db, err := NewFileDB(tmpFile)
		if err != nil {
			t.Fatalf("NewFileDB failed: %s", err)
		}

		// Check if the data slice is initialized but empty
		if db.data == nil {
			t.Fatal("expected data to be initialized")
		}
		if len(db.data) != 0 {
			t.Fatalf("expected 0 records, got %d", len(db.data))
		}

		// Ensure the mutex is properly initialized
		if db.fileMutex == nil {
			t.Fatal("expected fileMutex to be initialized")
		}

		// Ensure the file is correctly assigned
		if db.file == nil {
			t.Fatal("expected file to be assigned")
		}
	})

	t.Run("Invalid JSON File", func(t *testing.T) {
		// Create a temporary file with invalid JSON data for testing
		tmpFile, err := os.CreateTemp("", "testdb_invalid_*.json")
		if err != nil {
			t.Fatalf("failed to create temp file: %s", err)
		}
		defer os.Remove(tmpFile.Name()) // Clean up after the test

		// Write invalid JSON data to the file
		invalidData := `{"id": 1, "name": "Alice", "age": 30}` // Not an array
		if _, err := tmpFile.Write([]byte(invalidData)); err != nil {
			t.Fatalf("failed to write invalid data to temp file: %s", err)
		}

		// Make sure the file is at the beginning for reading
		if _, err := tmpFile.Seek(0, 0); err != nil {
			t.Fatalf("failed to seek to the beginning of the file: %s", err)
		}

		// Test NewFileDB function with invalid JSON data
		_, err = NewFileDB(tmpFile)
		if err == nil {
			t.Fatal("expected an error due to invalid JSON, but got none")
		}
	})
}

func Test_ReadRecord(t *testing.T) {
	tests := []struct {
		name         string
		initialData  string
		recordID     uint32
		expectedData map[string]interface{}
		expectedErr  error
	}{
		{
			name:        "Read existing record",
			initialData: `[{"id": 1, "name": "John Doe"}, {"id": 2, "name": "Jane Doe"}]`,
			recordID:    1,
			expectedData: map[string]interface{}{
				"id":   float64(1),
				"name": "John Doe",
			},
			expectedErr: nil,
		},
		{
			name:         "Read non-existing record",
			initialData:  `[{"id": 1, "name": "John Doe"}, {"id": 2, "name": "Jane Doe"}]`,
			recordID:     3,
			expectedData: nil,
			expectedErr:  ErrRecordNotFound,
		},
		{
			name:         "Read from empty file",
			initialData:  `[]`,
			recordID:     1,
			expectedData: nil,
			expectedErr:  ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file
			tempFile, err := os.CreateTemp("", "testdb*.json")
			if err != nil {
				t.Fatalf("Failed to create temporary file: %v", err)
			}
			defer os.Remove(tempFile.Name()) // Clean up the file afterwards

			// Write initial data to the temporary file
			if _, err := tempFile.WriteString(tt.initialData); err != nil {
				t.Fatalf("Failed to write initial data to temporary file: %v", err)
			}
			// Ensure the file pointer is at the beginning of the file
			if _, err := tempFile.Seek(0, 0); err != nil {
				t.Fatalf("Failed to seek to beginning of temporary file: %v", err)
			}

			// Initialize the FileDB
			db, err := NewFileDB(tempFile)
			if err != nil {
				t.Fatalf("Failed to initialize FileDB: %v", err)
			}

			// Test reading a record
			record, err := db.ReadRecord(tt.recordID)
			if err != tt.expectedErr {
				t.Fatalf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if tt.expectedErr == nil {
				ok, err := helpers.CompareMapsAsJSON(record, tt.expectedData)
				if err != nil {
					t.Fatalf("Error comparing records %s", err.Error())
				}
				if !ok {
					t.Errorf("Expected record %v, got %v", tt.expectedData, record)
				}
			}
		})
	}
}

func Test_CreateRecord(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "testdb*.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up the file afterwards

	// Write initial data to the temporary file
	initialData := `[{"id": 1, "name": "John Doe"}]`
	if _, err := tempFile.WriteString(initialData); err != nil {
		t.Fatalf("Failed to write initial data to temporary file: %v", err)
	}
	// Ensure the file pointer is at the beginning of the file
	if _, err := tempFile.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek to beginning of temporary file: %v", err)
	}

	// Initialize the FileDB
	db, err := NewFileDB(tempFile)
	if err != nil {
		t.Fatalf("Failed to initialize FileDB: %v", err)
	}

	// Create a new record
	newRecord := map[string]interface{}{
		"name": "Jane Doe",
	}

	err = db.CreateRecord(newRecord)
	if err != nil {
		t.Fatalf("Failed to create record: %v", err)
	}

	// Reopen the temporary file to check its contents
	tempFile.Seek(0, 0)
	var records []map[string]interface{}
	err = json.NewDecoder(tempFile).Decode(&records)
	if err != nil {
		t.Fatalf("Failed to decode records from temporary file: %v", err)
	}

	// Verify the new record was added correctly
	if len(records) != 2 {
		t.Fatalf("Expected 2 records, got %d", len(records))
	}

	expectedRecord := map[string]interface{}{
		"id":   float64(2), // IDs are unmarshalled as float64
		"name": "Jane Doe",
	}

	ok, err := helpers.CompareMapsAsJSON(records[1], expectedRecord)
	if err != nil {
		t.Fatalf("Error comparing records %s", err.Error())
	}

	if !ok {
		t.Errorf("Expected record %v, got %v", expectedRecord, records[1])
	}
}

func Test_UpdateRecord(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "testdb*.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up the file afterwards

	// Write initial data to the temporary file
	initialData := `[{"id": 1, "name": "John Doe"}, {"id": 2, "name": "Jane Doe"}]`
	if _, err := tempFile.WriteString(initialData); err != nil {
		t.Fatalf("Failed to write initial data to temporary file: %v", err)
	}
	// Ensure the file pointer is at the beginning of the file
	if _, err := tempFile.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek to beginning of temporary file: %v", err)
	}

	// Initialize the FileDB
	db, err := NewFileDB(tempFile)
	if err != nil {
		t.Fatalf("Failed to initialize FileDB: %v", err)
	}

	// Update an existing record
	updatedRecord := map[string]interface{}{
		"name": "John Smith",
	}

	err = db.UpdateRecord(1, updatedRecord)
	if err != nil {
		t.Fatalf("Failed to update record: %v", err)
	}

	// Reopen the temporary file to check its contents
	tempFile.Seek(0, 0)
	var records []map[string]interface{}
	err = json.NewDecoder(tempFile).Decode(&records)
	if err != nil {
		t.Fatalf("Failed to decode records from temporary file: %v", err)
	}

	// Verify the record was updated correctly
	expectedRecord := map[string]interface{}{
		"id":   float64(1), // IDs are unmarshalled as float64
		"name": "John Smith",
	}

	ok, err := helpers.CompareMapsAsJSON(records[0], expectedRecord)
	if err != nil {
		t.Fatalf("Error comparing records %s", err.Error())
	}

	if !ok {
		t.Errorf("Expected record %v, got %v", expectedRecord, records[0])
	}

	// Test updating a non-existing record
	nonExistentRecord := map[string]interface{}{
		"name": "Non Existent",
	}

	err = db.UpdateRecord(3, nonExistentRecord)
	if err == nil || err != ErrRecordNotFound {
		t.Fatalf("Expected error %v, got %v", ErrRecordNotFound, err)
	}
}

func Test_DeleteRecord(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "testdb*.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up the file afterwards

	// Write initial data to the temporary file
	initialData := `[{"id": 1, "name": "John Doe"}, {"id": 2, "name": "Jane Doe"}]`
	if _, err := tempFile.WriteString(initialData); err != nil {
		t.Fatalf("Failed to write initial data to temporary file: %v", err)
	}
	// Ensure the file pointer is at the beginning of the file
	if _, err := tempFile.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek to beginning of temporary file: %v", err)
	}

	// Initialize the FileDB
	db, err := NewFileDB(tempFile)
	if err != nil {
		t.Fatalf("Failed to initialize FileDB: %v", err)
	}

	// Delete an existing record
	err = db.DeleteRecord(1)
	if err != nil {
		t.Fatalf("Failed to delete record: %v", err)
	}

	// Reopen the temporary file to check its contents
	tempFile.Seek(0, 0)
	var records []map[string]interface{}
	err = json.NewDecoder(tempFile).Decode(&records)
	if err != nil {
		t.Fatalf("Failed to decode records from temporary file: %v", err)
	}

	// Verify the record was deleted correctly
	if len(records) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(records))
	}

	expectedRecord := map[string]interface{}{
		"id":   float64(2), // IDs are unmarshalled as float64
		"name": "Jane Doe",
	}

	ok, err := helpers.CompareMapsAsJSON(records[0], expectedRecord)
	if err != nil {
		t.Fatalf("Error comparing records %s", err.Error())
	}

	if !ok {
		t.Errorf("Expected record %v, got %v", expectedRecord, records[0])
	}

	// Test deleting a non-existing record
	err = db.DeleteRecord(3)
	if err == nil || err != ErrRecordNotFound {
		t.Fatalf("Expected error %v, got %v", ErrRecordNotFound, err)
	}
}

func Test_ConcurrentOperations(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "testdb*.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up the file afterwards

	// Write initial data to the temporary file
	initialData := `[{"id": 1, "name": "John Doe"}, {"id": 2, "name": "Jane Doe"}]`
	if _, err := tempFile.WriteString(initialData); err != nil {
		t.Fatalf("Failed to write initial data to temporary file: %v", err)
	}
	// Ensure the file pointer is at the beginning of the file
	if _, err := tempFile.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek to beginning of temporary file: %v", err)
	}

	// Initialize the FileDB
	db, err := NewFileDB(tempFile)
	if err != nil {
		t.Fatalf("Failed to initialize FileDB: %v", err)
	}

	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup

	// Number of concurrent operations
	numOps := 100

	// Run concurrent CreateRecord, ReadRecord, UpdateRecord, and DeleteRecord operations
	for i := 0; i < numOps; i++ {
		wg.Add(5)

		go func(i int) {
			defer wg.Done()
			newRecord := map[string]interface{}{
				"name": fmt.Sprintf("User %d", i),
			}
			if err := db.CreateRecord(newRecord); err != nil {
				t.Errorf("Failed to create record: %v", err)
			}
		}(i)

		go func(i int) {
			defer wg.Done()
			if _, err := db.ReadRecord(uint32(i + 1)); err != nil && err != ErrRecordNotFound {
				t.Errorf("Failed to read record: %v", err)
			}
		}(i)

		go func(i int) {
			defer wg.Done()
			if _, err := db.ReadRecord(uint32(i + 1)); err != nil && err != ErrRecordNotFound {
				t.Errorf("Failed to read record: %v", err)
			}
		}(i)

		go func(i int) {
			defer wg.Done()
			updatedRecord := map[string]interface{}{
				"name": fmt.Sprintf("Updated User %d", i),
			}
			if err := db.UpdateRecord(uint32(i+1), updatedRecord); err != nil && err != ErrRecordNotFound {
				t.Errorf("Failed to update record: %v", err)
			}
		}(i)

		go func(i int) {
			defer wg.Done()
			if err := db.DeleteRecord(uint32(i + 1)); err != nil && err != ErrRecordNotFound {
				t.Errorf("Failed to delete record: %v", err)
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}
