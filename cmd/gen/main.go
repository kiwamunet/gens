package main

import (
	"github.com/droundy/goopt"
	gen "github.com/kiwamunet/gens"
)

var (
	sqlConnStr  = goopt.String([]string{"-c", "--connstr"}, "", "database connection string")
	sqlDatabase = goopt.String([]string{"-d", "--database"}, "", "Database to for connection")
	sqlTable    = goopt.String([]string{"-t", "--table"}, "", "Table to build struct from")

	packageName = goopt.String([]string{"-p", "--package"}, "", "name to set for package")
	outputPath  = goopt.String([]string{"-o", "--output"}, "", "name to set for package")

	json   = goopt.Flag([]string{"--json"}, []string{"--no-json"}, "Add json annotations (default)", "Disable json annotations")
	gorm   = goopt.Flag([]string{"--gorm"}, []string{}, "Add gorm annotations (tags)", "")
	guregu = goopt.Flag([]string{"--guregu"}, []string{}, "Add guregu null types", "")
)

func init() {
	// Setup goopts
	goopt.Description = func() string {
		return "ORM and RESTful API generator for Mysql"
	}
	goopt.Version = "0.1"
	goopt.Summary = `gen --connstr "user:password@/dbname" --package pkgName --database databaseName --table tableName --output outputdir [--json] [--gorm] [--guregu]`

	//Parse options
	goopt.Parse(nil)
}

func main() {
	g := gen.NewGen()
	db := gen.DBConfig{
		ConnStr:  sqlConnStr,
		Database: sqlDatabase,
		Table:    sqlTable,
	}
	g.DB = db
	g.Package = packageName
	g.OutPutPath = outputPath

	g.IsJSON = json
	g.IsGorm = gorm
	g.IsGuregu = guregu
	g.Gen()
}
