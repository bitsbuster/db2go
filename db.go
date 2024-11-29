package db2go

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectionString defines the details required to establish a connection to a database.
type ConnectionString struct {
	// Host specifies the hostname or IP address of the database server.
	Host string
	// Port is the port number on which the database server is listening.
	Port uint16
	// Timeout is the maximum amount of time (in seconds) to wait for the database connection to be established.
	Timeout uint16
	// User is the username used for authenticating to the database.
	User string
	// Password is the password associated with the User for database authentication.
	Password string
	// DatabaseName is the name of the specific database to connect to on the server.
	DatabaseName string
}

// TableDescriptor represents the schema details of a single column in a database table.
type TableDescriptor struct {
	// Field is the name of the column in the table.
	Field string
	// Type is the data type of the column, as defined in the database schema (e.g., INT, VARCHAR(255)).
	Type string
	// Null indicates whether the column can contain NULL values ("YES" or "NO").
	Null string
	// Key specifies if the column is part of a key (e.g., "PRI" for primary key, "UNI" for unique key).
	Key string
	// Default is the default value assigned to the column, if any. A nil value indicates no default.
	Default *string
	// Extra contains additional information about the column, such as auto-increment settings.
	Extra string
}

// GetDbConnection returns the connection to the database using ConnectionString
func GetDbConnection(c *ConnectionString) *sql.DB {

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&timeout=%ds", c.User, c.Password, c.Host, c.Port, c.DatabaseName, c.Timeout)
	conn, err := sql.Open("mysql", dbURI)

	if err != nil {
		fmt.Println("failed creating connection to DB")
		panic(err)
	}

	if err = conn.Ping(); err != nil {
		fmt.Println("cannot stablish connection with DB")
		panic(err)
	}

	return conn
}

// GetTableDescriptor returns the information from sql "describe <tableName>" query in a TableDescriptor struct
func GetTableDescriptor(conn *sql.DB, tableName string) []TableDescriptor {

	rows, err := conn.Query(fmt.Sprintf("describe %s", tableName))
	if err != nil {
		fmt.Println("failed querying table description")
		panic(err)
	}

	defer rows.Close()

	result := make([]TableDescriptor, 0)
	for rows.Next() {
		r := TableDescriptor{}

		err = rows.Scan(&r.Field, &r.Type, &r.Null, &r.Key, &r.Default, &r.Extra)
		if err != nil {
			fmt.Println("failed scanning table description row")
			panic(err)
		}

		result = append(result, r)
	}

	return result
}
