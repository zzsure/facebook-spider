package crawler

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/op/go-logging"
	"gitlab.azbit.cn/web/facebook-spider/conf"
	"gitlab.azbit.cn/web/facebook-spider/consts"
	"gitlab.azbit.cn/web/facebook-spider/library/util"
	"gitlab.azbit.cn/web/facebook-spider/models"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

var logger = logging.MustGetLogger("modules/crawler")

// a cron task
func StartBasicCrawlTask(fds []*models.FileData) {
	for _, fd := range fds {
		url := util.GetMobilePostURL(fd.URL)
		logger.Info("crawl url:", url, " begin")
		b, err := util.RequestUrl(url)
		if err != nil {
			logger.Error("request url:", url, " err:", err)
		}
		// TODO: for test save html to data
		_ = util.SaveStringToFile("./data", "basic.html", string(b))

		break
	}
}

// a cron tas
func StartCrawlTask(fds []*models.FileData) {
	// TODO: add a cron
	for _, fd := range fds {
		//err := crawl(fd.URL, fd.Lang)
		url, err := util.GetOfficialAccountPostURL(fd.URL)
		if err != nil {
			logger.Error("parse url err:", err)
			continue
		}

		// crawl data to ads
		logger.Info("crawl url:", url, " begin")
		ads, err := crawlByGoquery(url, "en")
		//ads, err := crawlByColly(url, "en")

		if err != nil {
			logger.Error("crawl url:", url, " err:", err)
			continue
		}
		logger.Info("crawl url:", url, " end")

		// save article data to file
		err = saveArticleDataToFile(ads, fd.URL)
		if err != nil {
			logger.Error("save article data err:", err)
		}

		rs := rand.Intn(consts.MAX_SLEEP_TIME)
		logger.Info("random sleep seconds:", rs)
		time.Sleep(time.Duration(rs) * time.Second)
		break
	}
}

// save article data to file
func saveArticleDataToFile(ads []*models.ArticleData, url string) error {
	// get all posts and comments
	pm := make(map[string]string)
	cm := make(map[string]string)

	for _, ad := range ads {
		if _, ok := pm[ad.Date]; ok {
			pm[ad.Date] += "\n"
		}
		if ad.Posts != "" {
			pm[ad.Date] += ad.Posts
		}
		if len(ad.Comments) > 0 {
			cm[ad.Date] += strings.Join(ad.Comments, "\n")
		}
	}

	if len(pm) == 0 && len(cm) == 0 {
		logger.Info("len pm and cm is 0")
		return nil
	}

	postsDir, err := util.GetPostsDir(conf.Config.Spider.ArticleBaseDir, url)
	if err != nil {
		logger.Error("get posts path err:", err)
		return err
	}

	// save posts data to file
	for k, v := range pm {
		err = util.SaveStringToFile(postsDir, util.GetPostFileName(k), v)
		if err != nil {
			logger.Error("save ", k, " posts err:", err)
			continue
		}
	}

	// save comments data to file
	for k, v := range cm {
		err = util.SaveStringToFile(postsDir, util.GetCommentsFileName(k), v)
		if err != nil {
			logger.Error("save ", k, " comments err:", err)
			continue
		}
	}

	return nil
}

// craw data by goquery
func crawlByGoquery(url, lang string) ([]*models.ArticleData, error) {
	// request url get response
	b, err := util.RequestUrl(url)
	if err != nil {
		return nil, err
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(b)))
	if err != nil {
		logger.Error("document reader error:", err)
	}

	// TODO: for test
	//html, _ := util.ReadStringFromFile("./data/res_20191030.html")
	//doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	//if err != nil {
	//	logger.Error("document reader error:", err)
	//}

	var rets []*models.ArticleData
	// Find the review items
	doc.Find(".userContentWrapper").Each(func(i int, s *goquery.Selection) {
		posts := ""
		s.Find("p").Each(func(i int, s *goquery.Selection) {
			posts += strings.TrimLeft(s.Text(), " ") + "\n"
		})
		logger.Info("idx: ", i, "ret: ", posts)

		cellTime := s.Find(".timestampContent").Text()
		logger.Info("time string is:", cellTime)
		date := util.GetDateByCellTime(cellTime)
		logger.Info("date string is:", date)

		var comments []string
		ad := &models.ArticleData{
			Date:     date,
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

func crawlByColly(url, lang string) ([]*models.ArticleData, error) {
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	c.OnRequest(func(r *colly.Request) {
		//r.Headers.Set("Host", "facebook.com")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "/*")
		//r.Headers.Set("Origin", "http://facebook.com")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9;en-US,en;q=0.8")
		r.Headers.Set("Content-Type", "text/html")
		r.ResponseCharacterEncoding = "utf-8"
	})

	var err error

	c.OnResponse(func(resp *colly.Response) {
		ret, err := GzipDecode(resp.Body)
		if err == nil {
			fmt.Println(string(ret))
		} else {
			logger.Error("parse gzip err:", err)
		}
	})

	c.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
	})

	err = c.Visit(url)

	return nil, err
}

func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}
