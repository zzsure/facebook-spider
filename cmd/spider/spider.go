package spider

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/robfig/cron"
	"github.com/urfave/cli"
	"gitlab.azbit.cn/web/facebook-spider/conf"
	"gitlab.azbit.cn/web/facebook-spider/library/log"
	"gitlab.azbit.cn/web/facebook-spider/library/storage"
	"gitlab.azbit.cn/web/facebook-spider/library/util"
	"gitlab.azbit.cn/web/facebook-spider/modules/crawler"
	"time"
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
	storage.Init()

	// read urls from csv file
	fds, err := util.ReadUrlsFromCsv(conf.Config.Spider.CsvPath)
	if err != nil {
		logger.Error("read data from csv file err:", err)
		panic(err)
	}

	// start a crawl cron task
	cc := cron.New()
	str := fmt.Sprintf("%d %d * * *", conf.Config.Spider.StartHour, conf.Config.Spider.StartMin)
	_, _ = cc.AddFunc(str, func() {
		logger.Info("exec crawl cron unix time:", time.Now().Unix())
		err = crawler.StartBasicCrawlTask(fds)
		if err != nil {
			logger.Error("crawl err:", err)
		}
	})
	_, _ = cc.AddFunc("*/1 * * * *", func() {
		logger.Info("check alive....run 1 min cron")
	})
	cc.Start()
	defer cc.Stop()

	select {}
}
