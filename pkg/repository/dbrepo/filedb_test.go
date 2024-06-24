package dbrepo

import (
	"os"
	"testing"
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

// func Test_NewFileDB(t *testing.T) {
// 	t.Run("Successful DB creation", func(t *testing.T) {
// 		filepath := "./../../../testdata/testBase.jsonl"
// 		file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
// 		if err != nil {
// 			t.Errorf("Failed to open file for test")
// 		}
// 		defer file.Close()

// 		// Test initializing the FileDB
// 		db, err := NewFileDB(file)
// 		if err != nil {
// 			t.Fatalf("Failed to initialize FileDB: %v", err)
// 		}

// 		// Check that the file and mutex are correctly initialized
// 		if db.file == nil {
// 			t.Errorf("FileDB.file is nil")
// 		}
// 		if db.fileMutex == nil {
// 			t.Errorf("FileDB.fileMutex is nil")
// 		}
// 	})
// }

// func Test_ReadRecord(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		recordID       uint32
// 		expectedRecord string
// 		expectedErr    error
// 	}{
// 		{
// 			name:           "Successful record read",
// 			recordID:       1,
// 			expectedRecord: `{"id": 1, "name": "Alice", "details": {"age": 30, "city": "New York"}}`,
// 			expectedErr:    nil,
// 		},
// 		{
// 			name:           "Record not existing",
// 			recordID:       999,
// 			expectedRecord: "",
// 			expectedErr:    ErrRecordNotFound,
// 		},
// 	}

// 	filepath := "./../../../testdata/testBase.jsonl"
// 	file, err := os.OpenFile(filepath, os.O_RDONLY, 0)
// 	if err != nil {
// 		t.Errorf("Failed to open file for test")
// 	}
// 	defer file.Close()

// 	// Test initializing the FileDB
// 	db, err := NewFileDB(file)
// 	if err != nil {
// 		t.Fatalf("Failed to initialize FileDB: %v", err)
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			actualRecord, err := db.ReadRecord(tt.recordID)

// 			if !errors.Is(err, tt.expectedErr) {
// 				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
// 			}

// 			if tt.expectedRecord != "" {
// 				equal, err := helpers.CompareJSONWithMap(tt.expectedRecord, actualRecord)
// 				if err != nil {
// 					t.Fatalf("got error comparing records: %s", err.Error())
// 				}

// 				if !equal {
// 					t.Errorf("records not equal as expected %v and %s", actualRecord, tt.expectedRecord)
// 				}
// 			} else {
// 				if actualRecord != nil {
// 					t.Errorf("expected no record, but got %v", actualRecord)
// 				}
// 			}
// 		})
// 	}
// }

// func Test_CreateRecord(t *testing.T) {
// 	t.Run("Successful record creation", func(t *testing.T) {
// 		// Create a temporary file for testing
// 		file, err := os.CreateTemp("./../../../testdata/", "testdb*.jsonl")
// 		if err != nil {
// 			t.Fatalf("Failed to create temp file: %v", err)
// 		}
// 		defer os.Remove(file.Name())
// 		defer file.Close()

// 		// Initialize FileDB with the temp file
// 		db, err := NewFileDB(file)
// 		if err != nil {
// 			t.Fatalf("Failed to initialize FileDB: %v", err)
// 		}
// 		testRecord := `{"name":"Bob","details":{"age":25,"city":"Los Angeles"}}`
// 		expectedDBRecord := `{"id":1,"name":"Bob","details":{"age":25,"city":"Los Angeles"}}`

// 		var testRecordMap map[string]interface{}
// 		err = json.Unmarshal([]byte(testRecord), &testRecordMap)
// 		if err != nil {
// 			t.Fatalf("error unmarshaling expected body: %v", err)
// 		}

// 		err = db.CreateRecord(testRecordMap)
// 		if err != nil {
// 			t.Fatalf("Failed to create record: %v", err)
// 		}
// 		// Check if record modified as expected
// 		equals, err := helpers.CompareJSONWithMap(expectedDBRecord, testRecordMap)
// 		if err != nil {
// 			t.Fatalf("Failed to create record: %v", err)
// 		}

// 		if !equals {
// 			t.Errorf("Record not assigned with expected id")
// 		}

// 		// Check if record is in fle
// 		found, err := jsonRecordInFile(expectedDBRecord, file)
// 		if err != nil {
// 			t.Fatalf("Failed to create record: %v", err)
// 		}

// 		if !found {
// 			t.Errorf("Expected record not in db")
// 		}
// 	})
// }

// func jsonRecordInFile(record string, file *os.File) (bool, error) {
// 	file.Seek(0, 0)
// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		recordInFile := scanner.Text()

// 		found, err := helpers.CompareJSONStrings(recordInFile, record)
// 		if err != nil {
// 			return false, err
// 		}

// 		if found {
// 			return true, nil
// 		}

// 	}
// 	return false, nil
// }
