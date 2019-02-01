package meta

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/guregu/null"
	gtmpl "github.com/kiwamunet/gens/template"
	"github.com/knq/snaker"
)

type ModelInfo struct {
	PackageName     string
	Enum            string
	StructName      string
	ShortStructName string
	TableName       string
	Fields          []string
}

var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
}

var intToWordMap = []string{
	"zero",
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}

// Constants for return types of golang
const (
	golangByteArray          = "[]byte"
	gureguNullInt            = "null.Int"
	sqlNullInt               = "sql.NullInt64"
	golangInt                = "int"
	golangInt64              = "int64"
	gureguNullFloat          = "null.Float"
	sqlNullFloat             = "sql.NullFloat64"
	golangFloat              = "float"
	golangFloat32            = "float32"
	golangFloat64            = "float64"
	gureguNullString         = "null.String"
	sqlNullString            = "sql.NullString"
	gureguNulljsonRawMessage = "*json.RawMessage"
	sqlNulljsonRawMessage    = "*json.RawMessage"
	gureguNullTime           = "null.Time"
	golangTime               = "*time.Time"
	gureguNullBool           = "null.Bool"
	sqlNullBool              = "sql.NullBool"
	golangBool               = "bool"
)

// GenerateStruct generates a struct for the given table.
func GenerateStruct(db *sql.DB, databaseName string, tableName string, structName string, pkgName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) (*ModelInfo, error) {

	columnDataTypes, err := getColumnsFromMysqlTable(db, databaseName, tableName)
	if err != nil {
		fmt.Println("Error in selecting column data information from mysql information schema")
		return &ModelInfo{}, err
	}

	fields, enum := generateFieldsTypes(db, *columnDataTypes, 0, jsonAnnotation, gormAnnotation, gureguTypes, tableName)

	return &ModelInfo{
		PackageName:     pkgName,
		Enum:            enum,
		StructName:      structName,
		TableName:       tableName,
		ShortStructName: strings.ToLower(string(structName[0])),
		Fields:          fields,
	}, nil
}

// getColumnsFromMysqlTable Select column details from information schema and return map of map
func getColumnsFromMysqlTable(db *sql.DB, databaseName, tableName string) (*map[string]map[string]string, error) {

	// Store colum as map of maps
	columnDataTypes := make(map[string]map[string]string)
	// Select columnd data from INFORMATION_SCHEMA
	// columnDataTypeQuery := "SELECT COLUMN_NAME, COLUMN_KEY, COLUMN_TYPE, COLUMN_DEFAULT, DATA_TYPE, IS_NULLABLE, ORDINAL_POSITION, EXTRA FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND table_name = ?"

	columnDataTypeQuery := "SELECT C.COLUMN_NAME, C.COLUMN_KEY, C.COLUMN_TYPE, C.COLUMN_DEFAULT, C.DATA_TYPE, C.IS_NULLABLE, C.ORDINAL_POSITION, C.EXTRA, S.INDEX_NAME FROM INFORMATION_SCHEMA.COLUMNS AS C LEFT JOIN INFORMATION_SCHEMA.STATISTICS AS S ON C.TABLE_NAME = S.TABLE_NAME AND C.COLUMN_NAME = S.COLUMN_NAME AND C.TABLE_SCHEMA = S.TABLE_SCHEMA WHERE C.TABLE_SCHEMA = ? AND C.TABLE_NAME = ?"

	rows, err := db.Query(columnDataTypeQuery, databaseName, tableName)
	if err != nil {
		fmt.Println("Error selecting from db: " + err.Error())
		return nil, err
	}
	if rows != nil {
		defer rows.Close()
	} else {
		return nil, errors.New("No results returned for table")
	}

	for rows.Next() {
		var column string
		var columnKey string
		var columnType string
		var columnDefault null.String
		var dataType string
		var nullable string
		var ordinalPos string
		var extra string
		var index null.String
		rows.Scan(&column, &columnKey, &columnType, &columnDefault, &dataType, &nullable, &ordinalPos, &extra, &index)
		columnDataTypes[column] = map[string]string{"value": dataType, "nullable": nullable, "primary": columnKey, "position": ordinalPos, "extra": extra, "columnType": columnType, "columnDefault": NullString(columnDefault), "index": NullString(index)}
	}
	return &columnDataTypes, err
}

// Generate fields string
func generateFieldsTypes(db *sql.DB, obj map[string]map[string]string, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, tableName string) ([]string, string) {

	enumStr := ""
	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	m := map[int]string{}
	for _, key := range keys {
		fieldName := fmtFieldName(stringifyFirstChar(key))
		mysqlType := obj[key]

		nullable := false
		nullStr := ";NOT NULL"
		if mysqlType["nullable"] == "YES" {
			nullable = true
			nullStr = ";NULL"
		}

		valueType := sqlTypeToGoType(mysqlType["value"], nullable, gureguTypes, fmt.Sprintf("%v_%v", strings.ToUpper(key[:1])+key[1:], tableName))

		primary := ""
		if mysqlType["primary"] == "PRI" {
			primary = ";primary_key"
		}
		columnType := fmt.Sprintf("type:%s", mysqlType["columnType"])

		if strings.Contains(mysqlType["columnType"], "enum") {
			enumStr = createEnum(key, enumStr, mysqlType["columnType"], tableName)
		}

		if mysqlType["primary"] == "PRI" && valueType == "int" {
			if strings.Contains(mysqlType["extra"], "auto_increment") {
				primary = primary + ";auto_increment:true"
				columnType = ""
				if gureguTypes {
					valueType = gureguNullInt
				} else {
					valueType = sqlNullInt
				}
			} else {
				primary = primary + ";auto_increment:false"
			}
		}

		index := ""
		if strings.ToLower(mysqlType["index"]) == strings.ToLower(fieldName) {
			index = ";index"
		}

		defaultValue := ""
		if mysqlType["columnDefault"] != "" {
			defaultValue = fmt.Sprintf(" DEFAULT %s", mysqlType["columnDefault"])
		} else if mysqlType["columnDefault"] == "" && columnType == "type:timestamp" && nullable {
			defaultValue = " NULL DEFAULT NULL"
		}

		if mysqlType["extra"] != "" && mysqlType["extra"] != "auto_increment" {
			defaultValue = fmt.Sprintf("%s %s", defaultValue, mysqlType["extra"])
		}

		pos, _ := strconv.Atoi(mysqlType["position"])

		var annotations []string
		if gormAnnotation == true {
			if columnType[5:9] == "enum" { //type:enum('ALLOW','DISABLE')
				annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s\" sql:\"%s\"", key, columnType)) //`json:"ecosystem" sql:"type:ENUM('NONE','APPLYING','COMPLETE')"`
			} else {
				annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s%s%s;%s%s\"", key, nullStr, primary, index, columnType, defaultValue))
			}
		}
		if jsonAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("json:\"%s\"", key))
		}
		if len(annotations) > 0 {
			row := fmt.Sprintf("%s %s `%s`",
				fieldName,
				valueType,
				strings.Join(annotations, " "))
			m[pos] = row

		} else {
			row := fmt.Sprintf("%s %s",
				fieldName,
				valueType)
			m[pos] = row
		}
	}

	// log.Println(len(m))
	fields := make([]string, 0, len(m))
	for i := 1; i < len(m)+1; i++ {
		fields = append(fields, m[i])
	}
	return fields, enumStr
}

func sqlTypeToGoType(mysqlType string, nullable bool, gureguTypes bool, enumStr string) string {
	switch mysqlType {
	case "int", "smallint", "mediumint":
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt
		}
		return golangInt
	case "tinyint":
		if nullable {
			if gureguTypes {
				return gureguNullBool
			}
			return sqlNullBool
		}
		return golangBool
	case "bigint":
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt
		}
		return golangInt64
	case "char", "varchar", "mediumtext", "text", "tinytext":
		if nullable {
			if gureguTypes {
				return gureguNullString
			}
			return sqlNullString
		}
		return "string"
	case "enum":
		if nullable {
			if gureguTypes {
				return enumStr
			}
			return enumStr
		}
		return enumStr
	case "longtext":
		if nullable {
			if gureguTypes {
				return gureguNulljsonRawMessage
			}
			return sqlNulljsonRawMessage
		}
		return "json.RawMessage"
	case "date", "datetime", "time", "timestamp":
		if nullable && gureguTypes {
			return gureguNullTime
		}
		return golangTime
	case "decimal", "double":
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			return sqlNullFloat
		}
		return golangFloat64
	case "float":
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			return sqlNullFloat
		}
		return golangFloat32
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		return golangByteArray
	}
	return ""
}

// fmtFieldName formats a string as a struct key
//
// Example:
// 	fmtFieldName("foo_id")
// Output: FooID
func fmtFieldName(s string) string {
	name := lintFieldName(s)
	runes := []rune(name)
	for i, c := range runes {
		ok := unicode.IsLetter(c) || unicode.IsDigit(c)
		if i == 0 {
			ok = unicode.IsLetter(c)
		}
		if !ok {
			runes[i] = '_'
		}
	}
	return string(runes)
}

type EnumInfo struct {
	Key    string
	Type   string
	Member []string
}

func createEnum(typeStr string, enumStr string, columnEnum string, tableName string) string {
	columnEnum = strings.Replace(columnEnum, "enum(", "", 1)
	columnEnum = strings.Replace(columnEnum, ")", "", 1)
	columnEnum = strings.Replace(columnEnum, "'", "", -1)
	columnEnums := strings.Split(columnEnum, ",")

	typeStr = strings.ToUpper(typeStr[:1]) + typeStr[1:]
	enumInfo := &EnumInfo{
		Key:    tableName,
		Type:   typeStr,
		Member: columnEnums,
	}

	t, err := getTemplate(gtmpl.EnumTmpl)
	if err != nil {
		fmt.Println("Error in loading model template: " + err.Error())
		return enumStr
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, enumInfo)
	if err != nil {
		fmt.Println("Error in rendering model: " + err.Error())
		return enumStr
	}
	return enumStr + buf.String()
}
func getTemplate(t string) (*template.Template, error) {
	var funcMap = template.FuncMap{
		"title":       strings.Title,
		"toLower":     strings.ToLower,
		"toSnakeCase": snaker.CamelToSnake,
	}

	tmpl, err := template.New("model").Funcs(funcMap).Parse(t)

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
