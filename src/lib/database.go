package lib

import (
	"log"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"fmt"
)

var DB *sql.DB
var err error

//初始化数据库链接
func init(){

	dataBaseContent := userName+":"+password+"@tcp("+host+")/"+database

	DB , err = sql.Open(driver,dataBaseContent)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Mysql is ok!")
}