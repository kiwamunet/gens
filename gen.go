package gen

import (
	"bytes"
	"database/sql"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jimsmart/schema"
	"github.com/kiwamunet/gens/meta"
	gtmpl "github.com/kiwamunet/gens/template"
	"github.com/knq/snaker"
)

type DBConfig struct {
	ConnStr  *string
	Database *string
	Table    *string
}

type Gen struct {
	DB         DBConfig
	Package    *string
	OutPutPath *string
	IsJSON     *bool
	IsGorm     *bool
	IsGuregu   *bool
	IsVerBose  *bool
}

func NewGen() *Gen {
	return &Gen{}
}

func (g *Gen) Gen() {

	// check param
	if g.DB.ConnStr == nil || *g.DB.ConnStr == "" {
		log.Println("sql connection string is required!")
		return
	}
	if g.DB.Database == nil || *g.DB.Database == "" {
		log.Println("database string is required!")
		return
	}

	db, err := sql.Open("mysql", *g.DB.ConnStr)
	if err != nil {
		fmt.Println("Error in open database: " + err.Error())
		return
	}
	defer db.Close()

	// parse or read tables
	var tables []string
	if *g.DB.Table != "" {
		tables = strings.Split(*g.DB.Table, ",")
	} else {
		tables, err = schema.TableNames(db)
		if err != nil {
			log.Println("Error in fetching tables information from mysql information schema")
			return
		}
	}
	// if packageName is not set we need to default it
	if g.Package == nil || *g.Package == "" {
		*g.Package = "model"
	}

	dir := ""
	if g.OutPutPath != nil && *g.OutPutPath != "" {
		dir = *g.OutPutPath + "/"
	}
	os.Mkdir(dir+*g.Package, 0777)

	t, err := getTemplate(gtmpl.ModelTmpl)
	if err != nil {
		fmt.Println("Error in loading model template: " + err.Error())
		return
	}

	var structNames []string
	// generate go files for each table
	for _, tableName := range tables {
		structName := meta.FmtFieldName(tableName)
		// structName = inflection.Singular(strucstName)
		structNames = append(structNames, structName)

		modelInfo, err := meta.GenerateStruct(db, *g.DB.Database, tableName, structName, *g.Package, *g.IsJSON, *g.IsGorm, *g.IsGuregu)
		if err != nil {
			fmt.Println("Error in creating struct from json: " + err.Error())
			return
		}

		var buf bytes.Buffer
		err = t.Execute(&buf, modelInfo)
		if err != nil {
			fmt.Println("Error in rendering model: " + err.Error())
			return
		}

		data, err := format.Source(buf.Bytes())

		if err != nil {
			fmt.Println("Error in formating source: " + err.Error())
			return
		}

		ioutil.WriteFile(filepath.Join(dir+*g.Package, tableName+".go"), data, 0777)
	}
}

func getTemplate(t string) (*template.Template, error) {
	var funcMap = template.FuncMap{
		"title":            strings.Title,
		"toLower":          strings.ToLower,
		"toLowerCamelCase": camelToLowerCamel,
		"toSnakeCase":      snaker.CamelToSnake,
	}

	tmpl, err := template.New("model").Funcs(funcMap).Parse(t)

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func camelToLowerCamel(s string) string {
	ss := strings.Split(s, "")
	ss[0] = strings.ToLower(ss[0])
	return strings.Join(ss, "")
}
