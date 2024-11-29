
# db2go

`db2go` is a utility written in Go that connects to a MySQL database, reads the schema of a specified table, and generates a Go struct representing the table's structure. The tool uses JSON struct tags to facilitate easy marshaling and unmarshaling of data.

## Features

- Connects to a MySQL database using the provided connection details.
- Reads the schema of a specified table using MySQL's `DESCRIBE` command.
- Generates Go structs based on the table's schema, with optional JSON struct tags.

## Usaage

   ```bash
   go get github.com/bitsbuster/db2go
   ```

### Example

Generate a Go struct for a table named `users` in a MySQL database:

Sample output:
```go
type Users struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Created  string `json:"created"`
    Modified string `json:"modified"`
}
```

## How It Works

1. Parses command-line arguments using the `conf` package.
2. Connects to the MySQL database using the provided connection details.
3. Uses MySQL's `DESCRIBE` command to retrieve the schema of the specified table.
4. Generates a Go struct from the schema, including JSON struct tags if the feature is enabled.

## License

This project is licensed under the [Apache 2.0 License](LICENSE).

---

