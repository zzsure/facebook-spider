package crawler

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/op/go-logging"
	"gitlab.azbit.cn/web/facebook-spider/conf"
	"gitlab.azbit.cn/web/facebook-spider/consts"
	"gitlab.azbit.cn/web/facebook-spider/library/storage"
	"gitlab.azbit.cn/web/facebook-spider/library/util"
	"gitlab.azbit.cn/web/facebook-spider/models"
	"math/rand"
	"strings"
	"time"
)

var logger = logging.MustGetLogger("modules/crawler")

// a cron task
func StartBasicCrawlTask1(fds []*models.FileData) error {
	// TODO: for test
	html, _ := util.ReadStringFromFile("./data/vogue.html")
	ads, err := parseArticle([]byte(html))

	// save article data to file
	err = saveArticleDataToFile(ads, "https://www.facebook.com/Vogue/")
	if err != nil {
		logger.Error("save article data err:", err)
	}
	return nil
}

// a cron task
func StartBasicCrawlTask(fds []*models.FileData) error {
	if !isLogin() {
		err := login()
		if err != nil {
			return err
		}
	}

	for _, fd := range fds {
		url := util.GetMobilePostURL(fd.URL)
		logger.Info("crawl url:", url, " begin")

		content, err := crawlByColly(url)
		if err != nil {
			return err
		}

		ads, err := parseArticle(content)
		logger.Info("url:", url, ", content:", string(content))

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
	return nil
}

func parseArticle(b []byte) ([]*models.ArticleData, error) {
	var rets []*models.ArticleData

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}
	doc.Find("div[role=article]").Each(func(i int, s *goquery.Selection) {
		posts := ""
		s.Find("span p").Each(func(i int, s *goquery.Selection) {
			posts += strings.TrimLeft(s.Text(), " ") + "\n"
		})

		cellTime := s.Find("abbr").Text()
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
	return rets, nil
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

// crawl by colly
func crawlByColly(url string) ([]byte, error) {
	c := colly.NewCollector()
	_ = c.SetStorage(storage.StorageIns)
	c.AllowURLRevisit = true
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", "facebook.com")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.8")
		r.ResponseCharacterEncoding = "utf-8"
	})

	var err error
	var content []byte

	c.OnResponse(func(resp *colly.Response) {
		content = resp.Body
		//logger.Info("crawl:", string(resp.Body))
	})

	c.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
	})

	errVisit := c.Visit(url)
	if errVisit != nil {
		return nil, errVisit
	}

	return content, err
}

// check login
func isLogin() bool {
	c := colly.NewCollector()
	_ = c.SetStorage(storage.StorageIns)
	for _, cookie := range c.Cookies(consts.LOGIN_CHECK_URL) {
		logger.Info("cookie:", cookie.String())
		if strings.Contains(cookie.String(), "c_user") {
			logger.Info("have login")
			return true
		}
	}
	logger.Info("have not login")
	return false
}

// log in mbasic facebook
func login() error {
	c := colly.NewCollector()
	_ = c.SetStorage(storage.StorageIns)
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", "facebook.com")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.8")
		r.ResponseCharacterEncoding = "utf-8"
	})

	var err error
	c.OnHTML("#login_form", func(e *colly.HTMLElement) {
		loginURL, exists := e.DOM.Attr("action")
		if !exists {
			err = errors.New("doesn't have action label")
			return
		}
		loginURL = fmt.Sprintf("%s%s", strings.TrimRight(consts.LOGIN_CHECK_URL, "/"), loginURL)
		logger.Info("login url is:", loginURL)

		reqMap := make(map[string]string)
		e.DOM.Find("input").Each(func(i int, s *goquery.Selection) {
			name, _ := s.Attr("name")
			value, _ := s.Attr("value")
			if name != "" && value != "" && name != "sign_up" {
				reqMap[name] = value
			}
		})
		reqMap["email"] = conf.Config.FaceBook.Account
		reqMap["pass"] = conf.Config.FaceBook.Password
		logger.Info("login req map:", reqMap)
		err = c.Post(loginURL, reqMap)
	})

	c.OnResponse(func(resp *colly.Response) {
		logger.Info("login:", string(resp.Body))
	})

	c.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
	})

	errVisit := c.Visit(consts.LOGIN_CHECK_URL)
	if errVisit != nil {
		return errVisit
	}

	return err
}
