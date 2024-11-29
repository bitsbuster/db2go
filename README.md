
# db2go

`db2go` is a utility written in Go that connects to a MySQL database, reads the schema of a specified table, and generates a Go struct representing the table's structure. The tool uses JSON struct tags to facilitate easy marshaling and unmarshaling of data.

## Features

- Connects to a MySQL database using the provided connection details.
- Reads the schema of a specified table using MySQL's `DESCRIBE` command.
- Generates Go structs based on the table's schema, with optional JSON struct tags.

## Usage

   ```bash
   go get github.com/bitsbuster/db2go
   ```

### Example

Generate a Go struct for a table named `users` in a MySQL database:

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ardanlabs/conf"
	"github.com/bitsbuster/db2go"
)

type Arguments struct {
	Host         string `conf:"flag:host,short:h,required"`
	Port         uint16 `conf:"flag:port,short:p,default:3306"`
	User         string `conf:"flag:user,short:u,required"`
	Password     string `conf:"flag:password,short:s,required"`
	DatabaseName string `conf:"flag:db,short:d,required"`
	Timeout      uint16 `conf:"flag:timeout,short:t,default:10"`
	Table        string `conf:"flag:table,short:b,required"`
}

func main() {

	arguments := &Arguments{}
	//check for execution argument to start app as bridge or redis dispatcher
	if err := conf.Parse(os.Args[1:], "", arguments); err != nil {
		log.Panic(err)
	}

	connection := db2go.GetDbConnection(&db2go.ConnectionString{
		Host:         arguments.Host,
		Port:         arguments.Port,
		User:         arguments.User,
		Password:     arguments.Password,
		DatabaseName: arguments.DatabaseName,
		Timeout:      arguments.Timeout,
	})

	defer connection.Close()

	tableDescriptor := db2go.GetTable(connection, arguments.Table)

	st := db2go.CreateStruct(tableDescriptor, arguments.Table, true)

	fmt.Printf("%s\n", st)
}
```

Build the binary:
   ```bash
   go build -o db2go
   ```

Then you can execute it:
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

## License

This project is licensed under the [Apache 2.0 License](LICENSE).

---

