package GoMission4

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io"
	"math"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var RDB *redis.Client

func RedisConnect() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",  // no password set
		DB:       0,   // use default DB
		PoolSize: 100, // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Redis连接失败 %v", err)
	}
	RDB = rdb
}

func RedisGenerate() {
	for range time.Tick(time.Second / 10) {
		// 连续推送记录
		rand.Seed(time.Now().UnixNano())
		TradeId := rand.Intn(307) + 93
		TradePrice := rand.Intn(100) + 1000
		TradeNum := rand.Intn(100) + 1000
		//fmt.Println(TradeId, TradePrice, TradeNum)
		msg := strconv.Itoa(TradeId) + " " + strconv.Itoa(TradePrice) + " " + strconv.Itoa(TradeNum)
		RDB.Publish(context.Background(), "MyChannel", msg)

	}

}

var db *gorm.DB

func Mysqlinit() {
	dsn := "root:ihaxi.net@tcp(127.0.0.1:3306)/go_test?charset=utf8&parseTime=True&loc=Local"
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	if err != nil {
		fmt.Printf("数据库连接失败 %v", err)
	}

	if DB.Error != nil {
		fmt.Printf("数据库不存在 %v", DB.Error)
	}
	db = DB
}

type Trade struct {
	Num          int
	Price        int
	MaxPrice     int // 一分钟 ；；
	MinPrice     int //；；
	InitialPrice int //；；

	MaxPrices5Min     int //；；
	MinPrices5Min     int //；；
	InitialPrices5min int //；；
	Nums5Min          int
}

var TradeResult [500]Trade // 尽最大可能容纳所有的交易对，每个时刻

type OneMinuteData struct {
	TradeId      int       `gorm:"column:trade_id"`
	Price        int       `gorm:"column:price"`
	MaxPrice     int       `gorm:"column:max_price"`
	MinPrice     int       `gorm:"column:min_price"`
	InitialPrice int       `gorm:"column:initial_price"`
	Num          int       `gorm:"num"`
	AddTime      time.Time `gorm:"add_time"`
}

type FiveMinuteData struct {
	TradeId      int       `gorm:"column:trade_id"`
	Price        int       `gorm:"column:price"`
	MaxPrice     int       `gorm:"column:max_price"`
	MinPrice     int       `gorm:"column:min_price"`
	InitialPrice int       `gorm:"column:initial_price"`
	Num          int       `gorm:"num"`
	AddTime      time.Time `gorm:"add_time"`
}

func SendMsgTo4(msg string) { // 发送给lark信息的模板
	//*json
	contentType := "application/text"
	//*data
	sendData := `{
	  "msg_type": "text",
	  "content": {"text": "` + msg + `"}
	}`
	//*request
	data := strings.NewReader(sendData)
	// 自己的lark链接
	result, err := http.Post("https://open.larksuite.com/open-apis/bot/v2/hook/bea77277-2559-4468-9060-9cd847560f00", contentType, data)
	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(result.Body) // 关闭
}

func SafeRun(entry func()) { // 安全运行的模版
	/*
		该函数传入一个匿名函数或闭包后的执行函数，
		当传入函数以任何形式发生 panic 崩溃后，可以将崩溃发生的错误打印出来，
		同时允许后面的代码继续运行，不会造成整个进程的崩溃。
	*/
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		switch err.(type) {
		case runtime.Error: // 运行时错误
			SendMsgTo4("程序运行时发生错误")
		default: // 非运行时错误
			fmt.Println("error:", err)
		}
	}()
	entry()
}

func MissionOne() {
	//四个定时器
	ticker := time.NewTicker(10 * time.Second)
	//每个一分钟的断点，快照
	SafeRun(func() {
		go func() {
			for {
				<-ticker.C
				const IntMaxFlag = math.MaxInt32 // 定义一个最大的变量
				const IntMinFlag = math.MinInt32 // 定义一个最小的变量
				var OneMinute [401]OneMinuteData
				for i := 93; i < 400; i++ { // 交易对的范围从93到400
					OneMinute[i].TradeId = i
					OneMinute[i].Price = TradeResult[i].Price
					OneMinute[i].Num = TradeResult[i].Num

					if TradeResult[i].MaxPrice == IntMinFlag && TradeResult[i].MinPrice == IntMaxFlag { //并没有更新
						OneMinute[i].MaxPrice = 0
						OneMinute[i].MinPrice = 0
					} else {
						OneMinute[i].MaxPrice = TradeResult[i].MaxPrice
						OneMinute[i].MinPrice = TradeResult[i].MinPrice
					}

					OneMinute[i].InitialPrice = TradeResult[i].InitialPrice
					OneMinute[i].AddTime = time.Now() //加入现在的时间

					if err := db.Create(&OneMinute[i]).Error; err != nil {
						fmt.Println("插入失败", err) // 插入数据库
					}

					TradeResult[i].MaxPrice = IntMinFlag
					TradeResult[i].MinPrice = IntMaxFlag
					TradeResult[i].InitialPrice = -100
					TradeResult[i].Num = 0
				}

			}
		}()
	})
}

var newId int

func MissionTwo() {
	ticker1 := time.NewTicker(500 * time.Millisecond)
	//每0.5秒使用Redis并且publish出去
	SafeRun(func() {
		go func() {
			for {
				<-ticker1.C
				msg1 := strconv.Itoa(newId) + " " + strconv.Itoa(TradeResult[newId].MaxPrice) + " " + strconv.Itoa(TradeResult[newId].MinPrice) + " " + strconv.Itoa(TradeResult[newId].Price) + " " + strconv.Itoa(TradeResult[newId].Num)
				RDB.Publish(context.Background(), "NewChannel", msg1)
			}
		}()
	})
}

func MissionThree() {
	ticker2 := time.NewTicker(300 * time.Second)
	//每五分种设置的断点，快照
	SafeRun(func() {
		go func() {
			for {
				<-ticker2.C
				const IntMaxFlag = math.MaxInt32 // 定义一个最大的变量
				const IntMinFlag = math.MinInt32 // 定义一个最小的变量
				var FiveMinute [401]FiveMinuteData
				for i := 93; i < 400; i++ {
					FiveMinute[i].TradeId = i
					FiveMinute[i].Price = TradeResult[i].Price
					FiveMinute[i].Num = TradeResult[i].Nums5Min
					if TradeResult[i].MaxPrices5Min == IntMinFlag && TradeResult[i].MinPrices5Min == IntMaxFlag {
						FiveMinute[i].MaxPrice = 0 // 该交易对没有更新
						FiveMinute[i].MinPrice = 0
					} else {
						FiveMinute[i].MaxPrice = TradeResult[i].MaxPrices5Min
						FiveMinute[i].MinPrice = TradeResult[i].MinPrices5Min
					}

					FiveMinute[i].InitialPrice = TradeResult[i].InitialPrices5min
					FiveMinute[i].AddTime = time.Now()
					// 拿到快照之后，写入数据库
					if err := db.Create(&FiveMinute[i]).Error; err != nil {
						fmt.Println("插入"+
							""+
							"失败", err)
					}

					TradeResult[i].MaxPrices5Min = IntMinFlag
					TradeResult[i].MinPrices5Min = IntMaxFlag
					TradeResult[i].InitialPrices5min = -100
					TradeResult[i].Nums5Min = 0
				}
			}
		}()
	})
}

func MissionFour() {
	//处理channel收到的消息，频率每秒种十次，也就是说要处理全部的消息
	ticker3 := time.NewTicker(100 * time.Millisecond)
	SafeRun(func() {
		//go func() {
		for {
			<-ticker3.C
			pub := RDB.Subscribe(context.Background(), "MyChannel") //订阅
			_, err := pub.Receive(context.Background())             //接收
			if err != nil {
				panic(err)
			}
			// 用管道来接收全部的消息
			ch := pub.Channel()
			for msg := range ch {
				TempRec := strings.Fields(msg.Payload) // 切割字符串
				// 这是一个局部变量 仅仅在该函数内部使用
				LocalId, err := strconv.Atoi(TempRec[0])
				if err != nil {
					return
				}
				TradeResult[LocalId].Price, err = strconv.Atoi(TempRec[1])
				if err != nil {
					return
				}
				tempNum, err := strconv.Atoi(TempRec[2]) // 交易量
				if err != nil {
					return
				}

				newId = LocalId // 重新复制给全局变量， 交易对的id

				TradeResult[LocalId].Num = TradeResult[LocalId].Num + tempNum
				TradeResult[LocalId].Nums5Min = TradeResult[LocalId].Nums5Min + tempNum

				if TradeResult[LocalId].InitialPrice == -100 { // 一分钟内第一次更新
					TradeResult[LocalId].InitialPrice, _ = strconv.Atoi(TempRec[1]) // 初始值更新
				}
				if TradeResult[LocalId].MaxPrice < TradeResult[LocalId].Price {
					TradeResult[LocalId].MaxPrice, _ = strconv.Atoi(TempRec[1]) //最大值略小，最大值更新
				}
				if TradeResult[LocalId].MinPrice > TradeResult[LocalId].Price {
					TradeResult[LocalId].MinPrice, _ = strconv.Atoi(TempRec[1]) //最小值略大，最小值更新
				}

				if TradeResult[LocalId].InitialPrices5min == -100 { //五分钟内第一次更新
					TradeResult[LocalId].InitialPrices5min, _ = strconv.Atoi(TempRec[1]) // 初始值更新
				}
				if TradeResult[LocalId].MaxPrices5Min < TradeResult[LocalId].Price {
					TradeResult[LocalId].MaxPrices5Min, _ = strconv.Atoi(TempRec[1]) //最大值略小，最大值更新
				}
				if TradeResult[LocalId].MinPrices5Min > TradeResult[LocalId].Price {
					TradeResult[LocalId].MinPrices5Min, _ = strconv.Atoi(TempRec[1]) //最小值略大，最小值更新
				}

			}
		}

	})
}
