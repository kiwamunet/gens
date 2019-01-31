package template

var EnumTmpl = `type {{.Type}} int
{{$type := .Type}}
const (
	{{range $index, $var := .Member}}{{ if eq $index 0 }}{{$var}} {{$type}} = iota
	{{else}}{{$var}}
	{{end}}{{end}}
)
	
func (em {{.Type}}) String() string {
	switch em {
	{{range $index, $var := .Member}}case {{$var}}:
		return "{{$var}}"
	{{end}}default:
		return "UNKNOWN"
	}	
}

`
