package template

var EnumTmpl = `{{$type := .Type}}{{$key := .Key}}
type {{.Type}}_{{$key}} int64

const (
	{{range $index, $var := .Member}}{{ if eq $index 0 }}{{$var}}_{{$key}} {{$type}}_{{$key}} = iota
	{{else}}{{$var}}_{{$key}}
	{{end}}{{end}}
)

var type{{.Type}}_{{$key}} = [...]string{
	{{range $index, $var := .Member}}"{{$var}}",
	{{end}}
}

func (em {{.Type}}_{{$key}}) String() string {
	return type{{.Type}}_{{$key}}[em]
}

func (e *{{.Type}}_{{$key}}) Scan(value interface{}) error {
	*e = {{.Type}}_{{$key}}(value.(int64))
	return nil
}

func (e {{.Type}}_{{$key}}) Value() (driver.Value, error) { return int64(e), nil }

`
