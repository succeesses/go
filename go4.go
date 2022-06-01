package main

import (
	"awesomeProject/GoMission4"
	_ "github.com/go-sql-driver/mysql"
	"math"
)

func main4() {

	MainThread()
}

const IntMaxFlag = math.MaxInt32 // 定义一个最大的变量
const IntMinFlag = math.MinInt32 // 定义一个最小的变量

func MainThread() {

	for i := 0; i < 500; i++ {
		// 开始时刻初始化所有的交易对
		GoMission4.TradeResult[i].Num = 0
		GoMission4.TradeResult[i].Price = 0
		GoMission4.TradeResult[i].Nums5Min = 0
		GoMission4.TradeResult[i].MaxPrices5Min = IntMinFlag
		GoMission4.TradeResult[i].MaxPrice = IntMinFlag
		GoMission4.TradeResult[i].MinPrice = IntMaxFlag
		GoMission4.TradeResult[i].MinPrices5Min = IntMaxFlag
		GoMission4.TradeResult[i].InitialPrice = -100      // 初始化一个不可能的标志变量
		GoMission4.TradeResult[i].InitialPrices5min = -100 // 初始化一个不能的标志变量
	}

	GoMission4.Mysqlinit() // 初始化mysql

	GoMission4.RedisConnect() // 连接Redis

	go GoMission4.RedisGenerate() // 随机生成数据发布到Redis的频道

	GoMission4.MissionOne()
	GoMission4.MissionTwo()
	GoMission4.MissionThree()
	GoMission4.MissionFour()

}
