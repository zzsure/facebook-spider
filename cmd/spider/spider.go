package spider

import (
    "os"
	"github.com/op/go-logging"
	"github.com/urfave/cli"
	"gitlab.azbit.cn/web/facebook-spider/conf"
	"gitlab.azbit.cn/web/facebook-spider/library/log"
    "github.com/gocarina/gocsv"
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

// csv file data format
type FileData struct {
	URL  string `csv:"url"`  // official accounts url
	Lang string `csv:"language"` // official accounts language
}

func run(c *cli.Context) {
	conf.Init(c.String("conf"), c.String("args"))
	log.Init()
	// read urls from csv file
	fds, err := readUrlsFromCsv(conf.Config.Spider.CsvPath)
	if err != nil {
		//logger.Error("read data from csv file err:", err)
		panic(err)
	}
	// start a crawl cron task
	startCrawlTask(fds)
}

// a cron tas
func startCrawlTask(fds []*FileData) {
	// TODO: add a cron
	for _, fd := range fds {
        logger.Info("crawl url:", fd.URL, " begin")
		err := crawl(fd.URL, fd.Lang)
		if err != nil {
            logger.Error("crawl url:", fd.URL, " err:", err)
		}
        logger.Info("crawl url:", fd.URL, " end")
	}
}

// read url and language information from csv file
func readUrlsFromCsv(path string) ([]*FileData, error) {
    file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    fds := []*FileData{}
    if err := gocsv.UnmarshalFile(file, &fds); err != nil {
        return nil, err
    }
    return fds, nil
}

// crawl official accounts article data
func crawl(url, lang string) error {
	// use colly to crwal
    return nil
}
