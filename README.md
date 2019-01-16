# gens 


The gens tool produces golang structs from a given database for use in a .go file.
It supports [gorm](https://github.com/jinzhu/gorm) tags and implements some usable methods.

Generated datatypes include support for nullable columns [sql.NullX types](https://golang.org/pkg/database/sql/#NullBool) or [guregu null.X types](https://github.com/guregu/null)
and the expected basic built in go types.


## Usage
#### Command
```BASH
$ go get -v github.com/kiwamunet/gens/cmd/gen
$ gen -c "root:@tcp(127.0.0.1:3306)/{database}?parseTime=true" -o "./" -d nickel -p model --json --gorm --guregu

$ gen --help
Usage of gen:
	gen --connstr "user:password@/dbname" --package pkgName --database databaseName --table tableName --output outputdir [--json] [--gorm] [--guregu]
Options:
  -c, --connstr=   database connection string
  -d, --database=  Database to for connection
  -t, --table=     Table to build struct from
  -p, --package=   name to set for package
  -o, --output=    name to set for package
  --json           Add json annotations (default)
  --no-json        Disable json annotations
  --gorm           Add gorm annotations (tags)
  --guregu         Add guregu null types
  -h, --help       Show usage message

```

#### Import package
```Golang
go get -v github.com/kiwamunet/gens

import "github.com/kiwamunet/gens"

g := gen.NewGen()
genConnStr = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", Username, Password, Endpoint, DBName)
genDatabase = DBName

dbCon := gen.DBConfig{
    ConnStr:  &genConnStr,
    Database: &genDatabase,
    Table:    &genTable,
}
g.DB = dbCon
g.Package = &genPackage
g.OutPutPath = &genOutPut

g.IsJSON = &genJSON
g.IsGorm = &genGorm
g.IsGuregu = &genGuregu
g.Gen()
```

## Supported Databases

Currently Supported
- MySQL

#### Supported Datatypes

Currently only a limited number of datatypes are supported. Initial support includes:
-  tinyint (sql.NullBool or null.Bool)
-  int      (sql.NullInt64 or null.Int)
-  smallint      (sql.NullInt64 or null.Int)
-  mediumint      (sql.NullInt64 or null.Int)
-  bigint (sql.NullInt64 or null.Int)
-  decimal (sql.NullFloat64 or null.Float)
-  float (sql.NullFloat64 or null.Float)
-  double (sql.NullFloat64 or null.Float)
-  datetime (null.Time)
-  time  (null.Time)
-  date (null.Time)
-  timestamp (null.Time)
-  var (sql.String or null.String)
-  enum (sql.String or null.String)
-  varchar (sql.String or null.String)
-  longtext (sql.String or null.String)
-  mediumtext (sql.String or null.String)
-  text (sql.String or null.String)
-  tinytext (sql.String or null.String)
-  binary
-  blob
-  longblob
-  mediumblob
-  varbinary