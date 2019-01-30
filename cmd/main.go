package main

// import (
// 	"database/sql"
// 	"log"

// 	"github.com/guregu/null"
// 	"github.com/jinzhu/gorm"
// 	_ "github.com/jinzhu/gorm/dialects/mysql"
// 	"github.com/kiwamunet/gens/cmd/dd"
// )

// func main() {
// 	db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/**********?parseTime=true")
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	defer db.Close()
// 	log.Println(db.HasTable(&dd.Reward{}))
// 	m := dd.Reward{
// 		RewardCode: null.String{
// 			NullString: sql.NullString{
// 				String: "aaa",
// 				Valid:  true,
// 			},
// 		},
// 	}
// 	db.Create(&m)
// }
