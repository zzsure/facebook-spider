package storage

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/redisstorage"
	"gitlab.azbit.cn/web/facebook-spider/conf"
)

var StorageIns *redisstorage.Storage

func Init() {
	StorageIns = &redisstorage.Storage{
		Address:  conf.Config.Redis.Address,
		Password: conf.Config.Redis.Password,
		DB:       conf.Config.Redis.DB,
		Prefix:   conf.Config.Redis.Prefix,
	}
	c := colly.NewCollector()
	err := c.SetStorage(StorageIns)
	if err != nil {
		//panic(err)
	}
}
