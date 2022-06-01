package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type iwalaCurrency struct {
	CurrencyId   int `gorm:"primaryKey"`
	CurrencyName string
}

//func (v iwalaCurrency) TableName() string {
//	return "iwalaCurrency"
//}
// todo 看看是否需要绑定表名

type iwalaDeposit struct {
	Id         int `gorm:"primaryKey"`
	CurrencyId int
	Num        int //数量
	AddTime    time.Time
}
type iwalaWithdraw struct {
	Id         int `gorm:"primaryKey"`
	CurrencyId int
	Num        int // 数量
	AddTime    time.Time
}

type Result struct {
	CurrencyId int
	Auditor    int
	TotalNum   float32
	TIME       string
}

func SendMsg(msg string) {
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

func main1() {
	dsn := "szb02:iwala.netszb02@tcp(47.115.143.171:3306)/szb02?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	fmt.Println("------------------正在测试连接数据库---------------------")
	if err != nil {
		panic("连接数据库失败")
	} else {
		fmt.Println("数据库连接成功")
	}

	var total1 int64 = 0
	db.Model(iwalaCurrency{}).Count(&total1) // 打印每个表的总行数
	fmt.Println("表一的总记录条数", total1)

	var total2 int64 = 0
	db.Model(iwalaDeposit{}).Count(&total2) // 打印每个表的总行数
	fmt.Println("表二的总记录条数", total2)

	var total3 int64 = 0
	db.Model(iwalaWithdraw{}).Count(&total3) // 打印每个表的总行数
	fmt.Println("表三的总记录条数", total3)

	var results []Result
	sql := "select currency_id, auditor, sum(count1) as total_num, time from((SELECT currency_id, auditor, date_format(add_time, '%Y-%m-%d') as time, sum(num) as count1 from iwala_deposit group by currency_id) union all (SELECT currency_id, auditor, date_format(add_time, '%Y-%m-%d') as time, sum(num) as count1 from iwala_withdraw group by currency_id)) as total group by total.currency_id having total.auditor <> 0;"
	db.Raw(sql).Scan(&results)
	//fmt.Print(results) //打印数组结构体
	for i := 0; i < len(results); i++ {
		res1 := results[i].CurrencyId
		res2 := results[i].Auditor
		res3 := results[i].TotalNum
		res4 := results[i].TIME

		result := "今天的日期是:" + res4 + "  货币的ID是:" + strconv.Itoa(res1) + "  审核人的ID:" + strconv.Itoa(res2) + "  该货币的今天交易的总数量是(取整数)：" + strconv.Itoa(int(res3))
		//fmt.Println(result)
		SendMsg(result)
	}
}
