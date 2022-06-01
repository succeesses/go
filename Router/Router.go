package Router

import (
	"awesomeProject/API"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	//router := gin.Default()
	router := gin.New()
	// 接口

	//接口1
	router.GET("/input1/:MemberId/:Type/:Limit/:Offset", API.APInput1)

	//接口2
	router.POST("/input2/:Type", API.APInput2)

	//接口3
	router.POST("/input3/:time1/:time2", API.APInput3)

	return router
}
