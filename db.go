package db2go

import (
	"database/sql"
	"fmt"

	// _ "github.com/ardanlabs/conf"
	_ "github.com/go-sql-driver/mysql"
)

type ConnectionString struct {
	Host         string
	Port         uint16
	User         string
	Password     string
	DatabaseName string
	Timeout      uint16
}

type TableDescriptor struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default *string
	Extra   string
}

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

func GetTable(conn *sql.DB, table string) []TableDescriptor {

	rows, err := conn.Query(fmt.Sprintf("describe %s", table))
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
