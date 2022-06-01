package Mysql

import (
	"awesomeProject/Config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// MysqlInit 初始化mysql
func MysqlInit() {

	Config.InitConfig()
	Config := Config.ConfGlobal
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Config.Mysql.User, Config.Mysql.Password, Config.Mysql.Host, Config.Mysql.Port, Config.Mysql.Db)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	if err != nil {
		panic("连接数据库失败")
	} else {
		fmt.Println("数据库连接成功")
	}
	DB = db
}
