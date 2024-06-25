
# CRUD API for JSON Objects

This project is a test task that implements CRUD functionality for storing JSON objects with random content. It provides the following HTTP endpoints:

- **POST /records**: Creates a record with a required JSON body and returns the JSON object with the assigned ID in it.
- **GET /records/{id}**: Returns the record with the specified ID if it exists.
- **PUT /records/{id}**: Rewrites the object with the given ID for the now given object and returns the updated object.
- **DELETE /records/{id}**: Removes the record with the specified ID if it exists.

### Prerequisites

- Go 1.22.3

### Flags

The application accepts two optional flags:

- `-port`: Specifies the port on which the server will run. Default is `8080`.
- `-filepath`: Specifies the file path for the database file. Default is `./dbfile/db.json`.

Example:

```sh
go run ./... -port=9090 -filepath=./dbfile/db.json
```
