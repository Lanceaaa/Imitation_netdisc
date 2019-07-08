package mysql

import (
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
	"fmt"
	"os"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang?charset=utf8")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if(err != nil) {
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
}

// 返回数据库连接对象
func DBConn() *sql.DB {
	return db
}