package main

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type CorpInfo struct {
	Id         int64    `json:"id" gorm:"column:id;type:bigint(20) NOT NULL AUTO_INCREMENT;primaryKey"`
	CorpId     string   `json:"corpId" gorm:"type:varchar(50) NOT NULL"`
	CorpName   string   `json:"corpName" gorm:"type:varchar(255) NOT NULL"`
	Admin      string   `json:"admin" gorm:"type:varchar(50) NOT NULL"`
	StartTime  JsonTime `json:"startTime" gorm:"type:datetime NOT NULL"`
	StopTime   JsonTime `json:"stopTime" gorm:"type:datetime NOT NULL"`
	Phone      string   `json:"phone" gorm:"type:varchar(15) DEFAULT NULL"`
	CreateTime JsonTime `json:"createTime" gorm:"type:datetime NOT NULL"`
	UpdateTime JsonTime `json:"updateTime" gorm:"type:datetime NOT NULL;autoUpdateTime:milli"`
}

func (CorpInfo) TableName() string {
	return "t_corpinfo"
}

// 定义类型别名
type JsonTime time.Time

// 实现它的json序列化方法
func (this JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(this).Format("2006-01-02 15:04:05")) // Format内即是你想转换的格式
	return []byte(stamp), nil
}

var json_iterator = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	dsn := "root:king@tcp(127.0.0.1:3306)/proxy_server?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("open mysql fail")
		return
	}
	sqlDB, err := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(2)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(5)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	searchTest(db)

}

func searchTest(db *gorm.DB) {
	var corpInfo CorpInfo
	result := db.First(&corpInfo)
	//result := db.Where("corp_id = ?", "dingd53c3cce0ff62a0f35c2f4657eb6378f").Find(&corpInfo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("not fond data")
			return
		}
		fmt.Println("search err ", result.Error)
		return
	}
	s, _ := json_iterator.Marshal(corpInfo)
	fmt.Printf("first corp: \n%s\n", string(s))

	var corpInfos []CorpInfo
	result = db.Find(&corpInfos)

	s, _ = json_iterator.Marshal(corpInfos)
	fmt.Printf("all corp: \n%s\n", string(s))
}
