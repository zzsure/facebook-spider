package spider

import (
	"github.com/op/go-logging"
	"github.com/urfave/cli"
	"gitlab.azbit.cn/web/facebook-spider/conf"
	"gitlab.azbit.cn/web/facebook-spider/library/log"
	"gitlab.azbit.cn/web/facebook-spider/library/util"
	"gitlab.azbit.cn/web/facebook-spider/modules/crawler"
)

var Spider = cli.Command{
	Name:  "spider",
	Usage: "facebook spider",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "conf, c",
			Value: "config.toml",
			Usage: "toml配置文件",
		},
		cli.StringFlag{
			Name:  "args",
			Value: "",
			Usage: "multi config cmd line args",
		},
	},
	Action: run,
}

// log generation
var logger = logging.MustGetLogger("cmd/spider")

func run(c *cli.Context) {
	conf.Init(c.String("conf"), c.String("args"))
	log.Init()

	// read urls from csv file
	fds, err := util.ReadUrlsFromCsv(conf.Config.Spider.CsvPath)
	if err != nil {
		logger.Error("read data from csv file err:", err)
		panic(err)
	}

	// start a crawl cron task
	//crawler.StartCrawlTask(fds)
	err = crawler.StartBasicCrawlTask(fds)
	if err != nil {
		logger.Error("crawl err:", err)
	}
}
