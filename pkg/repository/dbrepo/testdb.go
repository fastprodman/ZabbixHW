package dbrepo

import "errors"

type TestDB struct {
	Data []map[string]interface{}
}

// Adds id to record and writes it into db
func (db *TestDB) CreateRecord(data map[string]interface{}) error {
	var newID uint32 = 1 // Default ID if the database is empty

	// Check if the database is not empty
	if len(db.Data) > 0 {
		// Get the ID of the last record
		lastRecord := db.Data[len(db.Data)-1]
		if id, ok := lastRecord["id"].(uint32); ok {
			newID = id + 1
		} else {
			return errors.New("invalid ID type in last record")
		}
	}

	// Set the new record's ID
	data["id"] = newID

	// Add the new record to the database
	db.Data = append(db.Data, data)

	return nil
}

// ReadRecord retrieves a record by its ID
func (db *TestDB) ReadRecord(id uint32) (map[string]interface{}, error) {
	// Iterate through the records to find the matching ID
	for _, record := range db.Data {
		if recordID, ok := record["id"].(uint32); ok {
			if recordID == id {
				return record, nil
			}
		} else {
			return nil, errors.New("invalid ID type in record")
		}
	}
	// If no matching record is found, return an error
	return nil, errors.New("record not found")
}

func (db *TestDB) UpdateRecord(id uint32, data map[string]interface{}) error {
	// Iterate through the records to find the matching ID
	data["id"] = id
	for i, record := range db.Data {
		if recordID, ok := record["id"].(uint32); ok {
			if recordID == id {
				// Update the record with new data, preserving the ID
				for key, value := range data {
					record[key] = value
				}
				db.Data[i] = record
				return nil
			}
		} else {
			return errors.New("invalid ID type in record")
		}
	}
	// If no matching record is found, return an error
	return errors.New("record not found")
}

// DeleteRecord deletes a record by its ID
func (db *TestDB) DeleteRecord(id uint32) error {
	// Iterate through the records to find the matching ID
	for i, record := range db.Data {
		if recordID, ok := record["id"].(uint32); ok {
			if recordID == id {
				// Remove the record from the slice
				db.Data = append(db.Data[:i], db.Data[i+1:]...)
				return nil
			}
		} else {
			return errors.New("invalid ID type in record")
		}
	}
	// If no matching record is found, return an error
	return errors.New("record not found")
}
