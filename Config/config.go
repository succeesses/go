package Config

import (
	"fmt"
	"github.com/go-ini/ini"
)

var ConfGlobal = &Config{}

type Config struct {
	//AppMode  string `ini:"app_mode"`
	//LogLevel string `ini:"log_level"`

	Mysql MysqlConfig `ini:"mysql"`
	Redis RedisConfig `ini:"redis"`
}

type MysqlConfig struct {
	User     string `ini:"user"`
	Password string `ini:"password"`
	Host     string `ini:"host"`
	Port     uint   `ini:"port"`
	Db       string `ini:"db"`
}
type RedisConfig struct {
	Ip   string `ini:"ip"`
	Port uint   `ini:"port"`
}

func InitConfig() {
	err := ini.MapTo(ConfGlobal, "Config/config.ini")
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("配置文件已经成功生效: ", ConfGlobal)
}
