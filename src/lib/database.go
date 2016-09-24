package lib

import (
	"log"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
)

const driver string = "mysql"
const userName string  ="next_tech"
const password string  ="00e4398aa6"
const host string = "tapi01.nomiss.hb02.allydata.cn:3306"
const database string = "nomiss"

var DB *sql.DB
var err error

//初始化数据库链接
func init(){

	dataBaseContent := userName+":"+password+"@tcp("+host+")/"+database

	DB , err = sql.Open(driver,dataBaseContent)

	if err != nil {
		log.Fatal(err.Error())
	}
}