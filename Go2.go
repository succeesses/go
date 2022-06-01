package main

import (
	"awesomeProject/Router"
)

func main() {

	//Mysql.MysqlInit()
	//Redis.RedisInit()
	//Logger.InitLogger()
	router := Router.InitRouter()
	router.Run(":8000") // 打开8000端口
}

//常驻进程
/*
func main() {
	cntxt := &daemon.Context{
		PidFileName: "sample.pid",
		PidFilePerm: 0644,
		LogFileName: "sample.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[go-daemon sample]"},
	}
	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()
	log.Print("- - - - - - - - - - - - - - -")
	log.Print("daemon started")
	worker()
}
func worker() {
	router := router.InitRouter()
	router.Run()
}
*/
