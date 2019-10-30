package spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocarina/gocsv"
	"github.com/gocolly/colly"
	"github.com/op/go-logging"
	"github.com/urfave/cli"
	"gitlab.azbit.cn/web/facebook-spider/conf"
	"gitlab.azbit.cn/web/facebook-spider/library/log"
	"gitlab.azbit.cn/web/facebook-spider/library/util"
	"net/http"
	"os"
	"path"
	"strings"
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
	URL  string `csv:"url"`      // official accounts url
	Lang string `csv:"language"` // official accounts language
}

type ArticleData struct {
	Posts    string   `json:"content"`
	Comments []string `json:"comments"`
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
		//err := crawl(fd.URL, fd.Lang)
		url, err := util.GetOfficialAccountURL(fd.URL)
		if err != nil {
			logger.Error("parse url err:", err)
			continue
		}

		// crawl data to ads
		logger.Info("crawl url:", url, " begin")
		ads, err := crawlByGoquery(url, "en")

		if err != nil {
			logger.Error("crawl url:", url, " err:", err)
			continue
		}
		logger.Info("crawl url:", url, " end")

		// get all posts and comments
		posts := ""
		comments := ""
		for i, ad := range ads {
			if i != 0 {
				posts += "\n"
			}
			posts += ad.Posts
			comments += strings.Join(ad.Comments, "\n")
		}

		// save data to file
		postsDir, err := getPostsDir(fd.URL)
		if err != nil {
			logger.Error("get posts path err:", err)
			continue
		}
		err = util.SaveStringToFile(postsDir, getPostFileName(), posts)
		if err != nil {
			logger.Error("save posts err:", err)
		}

		break
	}
}

func getPostFileName() string {
	// TODO：日期改成爬取的
	return fmt.Sprintf("posts-%s", util.GetCurrentDate())
}

func getPostsDir(url string) (string, error) {
	//dir := strings.TrimRight(conf.Config.Spider.ArticleBaseDir, "/")
	name, err := util.GetOfficialAccountName(url)
	if err != nil {
		return "", err
	}
	return path.Join(conf.Config.Spider.ArticleBaseDir, name), nil
	//return fmt.Sprintf("%s%s/%s%s", dir, url, util.GetCurrentDate(), "posts")
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

func crawlByGoquery(url, lang string) ([]*ArticleData, error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		logger.Error("http error:", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		logger.Error("status code error:", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logger.Error("document reader error:", err)
	}

	var rets []*ArticleData
	// Find the review items
	doc.Find("._4-u2 .userContent").Each(func(i int, s *goquery.Selection) {
		posts := ""
		var comments []string
		s.Find("p").Each(func(i int, s *goquery.Selection) {
			posts += strings.TrimLeft(s.Text(), " ") + "\n"
		})
		logger.Info("idx: ", i, "ret: ", posts)
		ad := &ArticleData{
			Posts:    posts,
			Comments: comments,
		}
		rets = append(rets, ad)
	})
	//doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
	//	// For each item found, get the band and title
	//	band := s.Find("a").Text()
	//	title := s.Find("i").Text()
	//	fmt.Printf("Review %d: %s - %s\n", i, band, title)
	//	ret = title
	//})
	return rets, nil
}

// crawl official accounts article data
func crawlByColly(url, lang string) (string, error) {
	// use colly to crwal
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36"
	c.OnRequest(func(r *colly.Request) {
		//r.Headers.Set("Host", "facebook.com")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		//r.Headers.Set("Origin", "https://facebook.com")
		//r.Headers.Set("Referer", "page_internal")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
	})

	var ret string
	var err error

	c.OnResponse(func(resp *colly.Response) {
		ret = string(resp.Body)
	})

	c.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
	})

	err = c.Visit(url)

	return ret, err
}
