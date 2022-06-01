package Model

import (
	"awesomeProject/Logger"
	"awesomeProject/Mysql"
	"awesomeProject/Redis"
	"context"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// 充币表结构体
type iwalaDeposit struct {
	Id         string `gorm:"column:id"`
	CurrencyId string `gorm:"column:currency_id"`
	Num        string `gorm:"column:num"`
	AddTime    string `gorm:"column:add_time"`
}

// 提币表结构体
type iwalaWithdraw struct {
	Id         string `gorm:"column:id"`
	CurrencyId string `gorm:"column:currency_id"`
	Num        string `gorm:"column:num"`
	AddTime    string `gorm:"column:add_time"`
}
type IwalaMember struct {
	MemberId int64  `gorm:"column:member_id"`
	UserName string `gorm:"column:user_name"`
	Email    string `gorm:"column:email"`
	Name     string `gorm:"column:name"`
}

// IwalaTrade type IwalaTrade struct {
//	CurrencyId int
//	Sum        float32
//	Fee        float32
//	AddTime    string
//}

type IwalaTrade struct {
	Id  int64   `gorm:"column:currency_id"`
	Num float64 `gorm:"column:sum(num)"`
	Fee float64 `gorm:"column:sum(fee)"`
}

var Member IwalaMember

func HandleAPI1Deposit(name int, limited int, offset int) (iwalaDep []iwalaDeposit, err error) {

	Mysql.MysqlInit()
	db := Mysql.DB // 连接mysql

	if err = db.Where("member_id = ?", name).Limit(limited).Offset(offset).Find(&iwalaDep).Error; err != nil {
		return
	}
	return
}

func HandleAPI1Withdraw(name int, limited int, offset int) (iwalaWith []iwalaWithdraw, err error) {

	Mysql.MysqlInit()
	db := Mysql.DB

	if err = db.Where("member_id = ?", name).Limit(limited).Offset(offset).Find(&iwalaWith).Error; err != nil {
		return
	}
	return
}

func HandleAPI2Insert(member IwalaMember) (err error) {

	Mysql.MysqlInit()
	db := Mysql.DB
	//添加数据
	result := db.Create(&member)
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

func HandleAPI2Delete(id int64) (err error) {

	Mysql.MysqlInit()
	db := Mysql.DB

	if err = db.Where("Member_Id = ?", id).Delete(&IwalaMember{}).Error; err != nil {
		fmt.Println("删除失败")
		return
	}
	fmt.Println("删除成功")
	return
}

func HandleAPI2Select(id int64) (members []IwalaMember, err error) {
	Mysql.MysqlInit()
	db := Mysql.DB
	if err = db.Where("Member_Id = ?", id).Take(&members).Error; err != nil {
		fmt.Println("查询失败")
		return
	}
	return
}

func HanldeAPI3(times1 string, times2 string) (trades []IwalaTrade, err error) {

	Mysql.MysqlInit()
	db := Mysql.DB
	Redis.RedisInit()
	Logger.InitLogger()
	Redis.CacheInit()

	if err = db.Select("currency_id,sum(num),sum(fee)").Where("add_time >= ? AND add_time <= ?", times1, times2).Group("currency_id").Find(&trades).Error; err != nil {
		return
	}
	data, _ := json.Marshal(trades)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	s1 := "接口3" + "_" + times1 + "_" + times2

	err = Redis.Redis.Set(ctx, s1, data, 0).Err()
	if err != nil {
		panic(err)
	}
	Logger.Logger.Info("查询成功", zap.String("货币Id", strconv.Itoa(int(trades[0].Id))), zap.String("交易费用", strconv.Itoa(int(trades[0].Fee))))

	Redis.GlobalCache.Set(s1, data, cache.DefaultExpiration)
	return
}
