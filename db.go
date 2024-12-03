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

// GetDbConnection establishes and returns a connection to a MySQL database.
//
// This function creates a database connection using the provided `ConnectionString`
// object, formats the connection URI, and verifies the connection by pinging the database.
//
// Parameters:
//   - c: *ConnectionString - A pointer to a `ConnectionString` struct containing
//     the database connection details, including user, password, host, port, database name,
//     and timeout.
//
// Returns:
//   - *sql.DB: A pointer to an established SQL database connection.
//
// Behavior:
//   - The function formats the connection string to include parsing of time values and a timeout.
//   - If the connection cannot be created or the database cannot be reached, the function
//     logs the error message and panics.
//
// Notes:
//   - The caller is responsible for closing the returned connection to avoid resource leaks.
//   - This function assumes a MySQL database and uses the Go `sql` package along with the
//     MySQL driver.
//   - Ensure the `ConnectionString` struct contains valid and properly formatted connection parameters.
//
// Example Usage:
//
//	connString := &ConnectionString{
//	    User:         "root",
//	    Password:     "password",
//	    Host:         "localhost",
//	    Port:         3306,
//	    DatabaseName: "my_database",
//	    Timeout:      5,
//	}
//	db := GetDbConnection(connString)
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

// GetTableDescriptor retrieves the column descriptors for a specified table.
//
// This function executes a "DESCRIBE" query on the provided table name using the
// database connection `conn`. It retrieves the column details and stores them
// as a slice of `TableDescriptor` objects, where each object contains metadata
// about a single column.
//
// Parameters:
//   - conn: *sql.DB - A pointer to an open SQL database connection.
//   - tableName: string - The name of the table to describe.
//
// Returns:
//   - []TableDescriptor: A slice of `TableDescriptor` objects containing metadata
//     about the columns of the specified table.
//
// Notes:
//   - This function will panic if there is an error executing the query or scanning
//     the rows. Ensure proper error handling and valid table names are used before
//     calling this function.
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

// GetDescriptorsForAllTables retrieves table descriptors for all tables in a database.
//
// This function queries the database connection `conn` to get the names of all tables
// using the `GetDbTableNames` function. It then iterates over each table name and
// retrieves its descriptors using the `GetTableDescriptor` function. The results
// are stored in a map where the keys are table names and the values are slices of
// `TableDescriptor` objects.
//
// Parameters:
//   - conn: *sql.DB - A pointer to an open SQL database connection.
//
// Returns:
//   - map[string][]TableDescriptor: A map where the key is the table name (string)
//     and the value is a slice of `TableDescriptor` containing metadata for the respective table.
func GetDescriptorsForAllTables(conn *sql.DB) map[string][]TableDescriptor {

	tables := GetDbTableNames(conn)

	result := make(map[string][]TableDescriptor)

	for _, t := range tables {

		result[t] = GetTableDescriptor(conn, t)

	}

	return result
}

// GetDbTableNames retrieves the names of all tables in the connected database.
//
// This function executes a "SHOW TABLES" query on the provided database connection `conn`
// to list all tables in the current database. It processes the query results, scans each
// table name, and appends it to a slice of strings.
//
// Parameters:
//   - conn: *sql.DB - A pointer to an open SQL database connection.
//
// Returns:
//   - []string: A slice containing the names of all tables in the database.
//
// Notes:
//   - This function will panic if there is an error executing the query or scanning
//     the rows. Ensure error handling and proper database connection setup before calling this function.
func GetDbTableNames(conn *sql.DB) []string {
	rows, err := conn.Query("show tables")
	if err != nil {
		fmt.Println("failed querying tables")
		panic(err)
	}

	defer rows.Close()

	result := make([]string, 0)
	for rows.Next() {
		r := ""

		err = rows.Scan(&r)
		if err != nil {
			fmt.Println("failed scanning table name row")
			panic(err)
		}

		result = append(result, r)
	}

	return result
}
