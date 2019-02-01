package template

var EnumTmpl = `type {{.Type}} int
{{$type := .Type}}{{$key := .Key}}
const (
	{{range $index, $var := .Member}}{{ if eq $index 0 }}{{$var}}_{{$key}} {{$type}} = iota
	{{else}}{{$var}}_{{$key}}
	{{end}}{{end}}
)
	
func (em {{.Type}}) String() string {
	switch em {
	{{range $index, $var := .Member}}case {{$var}}_{{$key}}:
		return "{{$var}}"
	{{end}}default:
		return "UNKNOWN"
	}	
}

`
