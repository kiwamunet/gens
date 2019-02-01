package template

var EnumTmpl = `{{$type := .Type}}{{$key := .Key}}
type {{.Type}}_{{$key}} int

const (
	{{range $index, $var := .Member}}{{ if eq $index 0 }}{{$var}}_{{$key}} {{$type}}_{{$key}} = iota
	{{else}}{{$var}}_{{$key}}
	{{end}}{{end}}
)
	
func (em {{.Type}}_{{$key}}) String() string {
	switch em {
	{{range $index, $var := .Member}}case {{$var}}_{{$key}}:
		return "{{$var}}"
	{{end}}default:
		return "UNKNOWN"
	}	
}

`
