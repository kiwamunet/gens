package template

var ModelTmpl = `package {{.PackageName}}

import (
    "database/sql"
    "database/sql/driver"
    "encoding/json"
    "time"

    "github.com/guregu/null"
)

var (
    _ = time.Second
    _ = sql.LevelDefault
    _ = null.Bool{}
    _ = json.Decoder{}
    _ = driver.NamedValue{} 
)

type {{.StructName}} struct {
    {{range .Fields}}{{.}}
    {{end}}
}

// TableName sets the insert table name for this struct type
func ({{.ShortStructName}} *{{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
`
