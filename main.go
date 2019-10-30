package main

import (
	"github.com/urfave/cli"
	"gitlab.azbit.cn/web/facebook-spider/cmd/spider"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "facebook-spider"
	app.Commands = []cli.Command{
        spider.Spider,
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
