
# db2go

`db2go` is a CLI utility written in Go that connects to a MySQL database, reads the schema of a specified table, and generates a Go struct representing the table's structure. The tool uses JSON struct tags to facilitate easy marshaling and unmarshaling of data.

## Features

- Connects to a MySQL database using the provided connection details.
- Reads the schema of a specified table using MySQL's `DESCRIBE` command.
- Generates Go structs based on the table's schema, with optional JSON struct tags.

## Installation

1. Make sure you have [Go installed](https://golang.org/dl/).
2. Clone this repository:
   ```bash
   git clone https://github.com/yourusername/db2go.git
   cd db2go
   ```
3. Build the binary:
   ```bash
   go build -o db2go
   ```

## Usage

Run the utility by providing the required arguments for connecting to the MySQL database and specifying the table:

```bash
./db2go \
  --host <db_host> \
  --port <db_port> \
  --user <db_user> \
  --password <db_password> \
  --db <db_name> \
  --timeout <timeout_seconds> \
  --table <table_name>
```

### Arguments

| Argument      | Short Flag | Description                            | Required | Default |
|---------------|------------|----------------------------------------|----------|---------|
| `--host`      | `-h`       | Hostname or IP address of the database | Yes      |         |
| `--port`      | `-p`       | Port number of the database            | No       | 3306    |
| `--user`      | `-u`       | Username for database authentication   | Yes      |         |
| `--password`  | `-s`       | Password for database authentication   | Yes      |         |
| `--db`        | `-d`       | Name of the database                   | Yes      |         |
| `--timeout`   | `-t`       | Connection timeout in seconds          | No       | 10      |
| `--table`     | `-b`       | Name of the table to process           | Yes      |         |

### Example

Generate a Go struct for a table named `users` in a MySQL database:

```bash
./db2go \
  --host 127.0.0.1 \
  --port 3306 \
  --user root \
  --password secret \
  --db my_database \
  --timeout 10 \
  --table users
```

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

## Dependencies

- [ardanlabs/conf](https://github.com/ardanlabs/conf) - For argument parsing.
- [bitsbuster/db2go](https://github.com/bitsbuster/db2go) - For database connection and schema processing.

## License

This project is licensed under the [Apache 2.0 License](LICENSE).

---

Feel free to adapt or extend this tool as needed! Contributions are welcome. 😊
