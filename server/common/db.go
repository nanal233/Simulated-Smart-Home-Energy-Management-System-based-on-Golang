package common

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func PrepareDatabase(dsn string) {
	// 初始化数据库连接
	// 需要指定 parseTime=True，否则无法解析时间。
	_db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	DB = _db
}
