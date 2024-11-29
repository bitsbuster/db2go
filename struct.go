package db2go

import (
	"fmt"
	"strings"
)

// CreateStruct creates a new struct with the camelized name of the table from the table descriptor.
// If withJson is true, creates the struct tags for json
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

// getType returns the field type from the mysql column type,
// having in mind for numbers if it is unsigned or it can be null (then returns a pointer)
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
