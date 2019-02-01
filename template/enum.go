package template

var EnumTmpl = `type {{.Type}} int
{{$type := .Type}}
const (
	{{range $index, $var := .Member}}{{ if eq $index 0 }}{{$var}}_{{$type}} {{$type}} = iota
	{{else}}{{$var}}_{{$type}}
	{{end}}{{end}}
)
	
func (em {{.Type}}) String() string {
	switch em {
	{{range $index, $var := .Member}}case {{$var}}_{{$type}}:
		return "{{$var}}"
	{{end}}default:
		return "UNKNOWN"
	}	
}

`
