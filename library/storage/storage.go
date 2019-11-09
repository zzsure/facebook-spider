package storage

import(
	"github.com/gocolly/colly"
    "github.com/gocolly/redisstorage"
)

var StorageIns *redisstorage.Storage

func Init() {
    StorageIns = &redisstorage.Storage{
        Address:  "127.0.0.1:6379",
        Password: "",
        DB:       3,
        Prefix:   "fb",
    }
    c := colly.NewCollector()
    err := c.SetStorage(StorageIns)
    if err != nil {
        panic(err)
    }
}
