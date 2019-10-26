package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gitlab.azbit.cn/web/facebook-spider/conf"
	"time"
)

var DB *gorm.DB

func Init() {
	var err error
	// 默认设置了隔离级别为RC，以避免间隙锁导致的死锁。
	// 间隙锁可以避免幻读，如果担心幻读，可以使用默认的隔离级别RR，把&tx_isolation=%%27READ-COMMITTED%%27去掉即可
	dsn := fmt.Sprintf("%s@%s/%s?charset=utf8&parseTime=True&loc=Local&timeout=3s&tx_isolation=%%27READ-COMMITTED%%27",
		conf.Config.Database.UserPassword, conf.Config.Database.HostPort, conf.Config.Database.DB)
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Println(dsn)
		panic(err)
	}
	DB.DB().SetConnMaxLifetime(time.Duration(conf.Config.Database.Conn.MaxLifeTime) * time.Second)
	DB.DB().SetMaxIdleConns(conf.Config.Database.Conn.MaxIdle)
	DB.DB().SetMaxOpenConns(conf.Config.Database.Conn.MaxOpen)
}
