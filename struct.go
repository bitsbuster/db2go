package db2go

import (
	"fmt"
	"os"
	"strings"
)

// CreateAllTablesStructFile generates Go struct definitions for multiple database tables
// and writes them to a specified file.
//
// This function takes a map of table names to their descriptors, generates Go struct
// definitions for each table using the `CreateStruct` function, and writes all the
// generated code to a file. The resulting file includes the specified package name.
//
// Parameters:
//   - filename: string - The name of the file where the generated structs will be written.
//   - packageName: string - The name of the Go package to include at the top of the file.
//   - descriptors: map[string][]TableDescriptor - A map where the keys are table names,
//     and the values are slices of `TableDescriptor` objects containing metadata about
//     the table columns.
//   - withJson: bool - A flag indicating whether to include JSON tags for the struct fields.
//
// Notes:
//   - The function uses the `CreateStruct` function to generate each struct definition.
//   - The `writeToFile` helper function is used to write the generated code to the specified file.
//   - Ensure the provided `filename` is writable, and the `packageName` is a valid Go package name.
//   - The file will contain all the structs, separated by newlines, under the specified package.
func CreateAllTablesStructFile(filename string, packageName string, descriptors map[string][]TableDescriptor, withJson bool) {

	builder := strings.Builder{}

	builder.WriteString("package ")
	builder.WriteString(packageName)
	builder.WriteString("\n\n")

	for k, v := range descriptors {

		builder.WriteString(CreateStruct(v, k, withJson))
		builder.WriteString("\n\n")

	}

	writeToFile(builder.String(), filename)
}

// CreateStruct generates a Go struct definition based on the table descriptors.
//
// This function takes a slice of `TableDescriptor` objects, a table name, and an
// optional flag for including JSON tags. It generates a Go struct definition where
// each column in the table corresponds to a struct field. The field names are camel-cased,
// and their types are determined based on the column descriptors.
//
// Parameters:
//   - tt: []TableDescriptor - A slice of `TableDescriptor` objects containing metadata
//     about the columns of the table.
//   - tableName: string - The name of the table, used as the base name for the generated struct.
//   - withJson: bool - A flag indicating whether to include JSON tags for the struct fields.
//
// Returns:
//   - string: A string representation of the generated Go struct.
//
// Panics:
//   - The function panics if the provided table descriptor slice is empty.
//
// Notes:
//   - The struct fields are formatted for alignment, ensuring consistent spacing.
//   - JSON tags are included in the struct definition if `withJson` is set to `true`.
//   - Helper functions like `Camelize` and `getType` are expected to handle field name
//     conversion and type determination, respectively.
func CreateStruct(tt []TableDescriptor, tableName string, withJson bool) string {

	if len(tt) < 1 {
		panic("table descriptor is empty")
	}

	withField := 0
	withType := 0
	temp := make([][]string, 0)

	for _, t := range tt {
		row := make([]string, 0)

		row = append(row, Camelize(t.Field, true))
		row = append(row, getType(t))
		if withJson {
			row = append(row, Camelize(t.Field, false))
		}
		if len(row[0]) > withField {
			withField = len(row[0])
		}
		if len(row[1]) > withType {
			withType = len(row[1])
		}
		temp = append(temp, row)
	}

	template := fmt.Sprintf("    %%-%ds %%-%ds", withField, withType)

	result := strings.Builder{}
	result.WriteString(fmt.Sprintf("type %sData struc {\n", Camelize(tableName, true)))

	for _, t := range temp {
		result.WriteString(fmt.Sprintf(template, t[0], t[1]))
		if len(t) == 3 {
			result.WriteString(fmt.Sprintf("\t`json:\"%s\"`", t[2]))
		}
		result.WriteString("\n")
	}

	result.WriteString("}")

	return result.String()
}

// getType determines the Go type corresponding to a database column type.
//
// This function maps a database column's type, as described in the `TableDescriptor`,
// to an appropriate Go type. It handles various database-specific nuances, such as
// detecting unsigned types, removing parentheses, and handling nullable fields.
// The resulting type is returned as a string suitable for use in a Go struct definition.
//
// Parameters:
//   - t: TableDescriptor - A descriptor of the table column, including its type, nullability,
//     and other metadata.
//
// Returns:
//   - string: The Go type corresponding to the column type, including pointer notation
//     if the column allows NULL values.
//
// Notes:
//   - Unsigned numeric types are prefixed with `u` to indicate unsigned integer types
//     (e.g., `uint64` for `BIGINT UNSIGNED`).
//   - Nullable columns are represented as pointers to their respective Go types (e.g., `*string`).
//   - Default Go types are provided for unknown column types, defaulting to `interface{}`.
//   - Time-related types are mapped to `time.Time`, and binary data types are mapped to `[]byte`.
//
// Example Mappings:
//   - `VARCHAR(255)` -> `string`
//   - `BIGINT UNSIGNED` -> `uint64`
//   - `DATETIME` -> `time.Time`
//   - `BOOL` -> `bool`
func getType(t TableDescriptor) string {

	cleanType := strings.ToUpper(t.Type)

	// Detects UNSIGNED and removes
	isUnsigned := strings.Contains(cleanType, "UNSIGNED")
	cleanType = strings.ReplaceAll(cleanType, "UNSIGNED", "")
	cleanType = strings.TrimSpace(cleanType)

	//removes parantesis
	posParentesis := strings.Index(cleanType, "(")
	if posParentesis > 0 {
		cleanType = cleanType[0:posParentesis]
	}

	result := strings.Builder{}
	if t.Null == "YES" {
		result.WriteString("*")
	}

	switch cleanType {
	case "VARCHAR", "TEXT", "CHAR", "ENUM", "SET", "LONGTEXT", "MEDIUMTEXT", "TINYTEXT":
		result.WriteString("string")
	case "BIGINT":
		if isUnsigned {
			result.WriteString("u") //
		}
		result.WriteString("int64")
	case "INT", "MEDIUMINT":
		if isUnsigned {
			result.WriteString("u") //
		}
		result.WriteString("int32")
	case "SMALLINT":
		if isUnsigned {
			result.WriteString("u") //
		}
		result.WriteString("int16")
	case "TINYINT":
		if isUnsigned {
			result.WriteString("u") //
		}
		result.WriteString("int8")
	case "FLOAT", "DOUBLE", "DECIMAL":
		result.WriteString("float64")
	case "DATE", "DATETIME", "TIMESTAMP", "TIME", "YEAR":
		result.WriteString("time.Time")
	case "BLOB", "LONGBLOB", "MEDIUMBLOB", "TINYBLOB", "BINARY", "VARBINARY":
		result.Reset()
		result.WriteString("[]byte")
	case "BIT", "BOOL", "BOOLEAN":
		result.WriteString("bool")
	default:
		result.Reset()
		result.WriteString("interface{}") // If the type is not known returns generic interface
	}
	return result.String()
}

// writeToFile appends a string value to a specified file.
//
// This function opens (or creates) a file with the specified filename, appends
// the given string value to it, and ensures the file is properly closed afterward.
//
// Parameters:
//   - value: string - The string content to write to the file.
//   - filename: string - The name of the file to which the content will be written.
//
// Behavior:
//   - If the file does not exist, it will be created.
//   - If the file exists, the content will be appended to the end of the file.
//   - The file is opened with permissions set to allow reading, writing, and creation
//     with mode `0644`.
//
// Panics:
//   - The function panics if there is an error opening the file or writing to it.
//
// Notes:
//   - Ensure appropriate error handling or pre-validation of file paths in production use.
//   - This function is primarily designed for simple file operations; for larger or more
//     complex I/O tasks, consider additional error handling or buffering.
func writeToFile(value, filename string) {
	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err := f.WriteString(value); err != nil {
		panic(err)
	}

}
