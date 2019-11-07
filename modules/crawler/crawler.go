package crawler

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/op/go-logging"
	"gitlab.azbit.cn/web/facebook-spider/conf"
	"gitlab.azbit.cn/web/facebook-spider/consts"
	"gitlab.azbit.cn/web/facebook-spider/library/util"
	"gitlab.azbit.cn/web/facebook-spider/models"
	"math/rand"
	"strings"
	"time"
)

var logger = logging.MustGetLogger("modules/crawler")

// a cron task
func StartBasicCrawlTask(fds []*models.FileData) error {
	_, err := crawlByColly("https://mbasic.facebook.com/", "en")
	if err != nil {
		logger.Error("crawl by colly err:", err)
	}
	return nil

	for _, fd := range fds {
		url := util.GetMobilePostURL(fd.URL)
		logger.Info("crawl url:", url, " begin")
		/*b, err := util.RequestUrl(url)
		if err != nil {
			logger.Error("request url:", url, " err:", err)
		}
		// TODO: for test save html to data
		_ = util.SaveStringToFile("./data", "basic_index.html", string(b))*/

		//html, _ := util.ReadStringFromFile("./data/basic.html")

		break
	}
	if !login() {
		return errors.New("login error")
	}
	return nil
}

// check login
func login() bool {
	b, _ := util.RequestUrl(consts.LOGIN_CHECK_URL)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(b)))
	if err != nil {
		logger.Error("document reader error:", err)
	}
	isLogin := true
	title := doc.Find("header title").Text()
	if strings.Contains(title, consts.LOGIN_CHECK_STRING) {
		isLogin = false
	}
	logger.Info("user login:", isLogin)
	if !isLogin {

	}
	return isLogin
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

func crawlByColly(url, lang string) ([]byte, error) {
	c := colly.NewCollector()
	//extensions.RandomUserAgent(c)
	//extensions.Referer(c)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", "facebook.com")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", "https://mbasic.facebook.com/")
		//r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Referer", "https://mbasic.facebook.com/")
		//r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9;en-US,en;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.8")
		r.Headers.Set("Content-Type", "text/html")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36")
		r.ResponseCharacterEncoding = "utf-8"
	})
	// https://mbasic.facebook.com/login/device-based/regular/login/?refsrc=https%3A%2F%2Fmbasic.facebook.com%2F&lwv=100&refid=8
	//lsd: AVr7ySCP
	//jazoest: 2671
	//m_ts: 1573108965
	//li: trzDXTJqUodlele3wo2bPUCq
	//try_number: 0
	//unrecognized_tries: 0
	//email: 18801090613
	//pass: I7*FUZd5
	//login: Log In
	// datr=jrnDXdWs41-YAcCqwuo7uNdS; sb=jrnDXbWnmlDBdnx8M-Nyjy6Q; c_user=100042468376443; xs=9%3AoN3FER5Grt_86g%3A2%3A1573108637%3A-1%3A-1; fr=32Np1ktt80sXSbimH.AWX1YD2zvO9LSbnIfQPPD_kCzHQ.Bdw7ud.O8.AAA.0.0.Bdw7ud.AWW7OU0J
	var err error

	c.OnHTML("#login_form", func(e *colly.HTMLElement) {
		loginURL, exists := e.DOM.Attr("action")
		if !exists {
			err = errors.New("doesn't have action label")
			return
		}
		loginURL = fmt.Sprintf("https://mbasic.facebook.com%s", loginURL)
		logger.Info("login url is:", loginURL)
		reqMap := make(map[string]string)
		e.DOM.Find("input").Each(func(i int, s *goquery.Selection) {
			name, _ := s.Attr("name")
			value, _ := s.Attr("value")
			if name != "" && value != "" && name != "sign_up" {
				reqMap[name] = value
			}
		})
		reqMap["email"] = "18810572605"
		reqMap["pass"] = "4ocjR&SN"
		logger.Info("req map:", reqMap)
		for _, cookie := range c.Cookies("https://mbasic.facebook.com") {
			logger.Info("cookie", cookie.Value)
		}
		err = c.Post(loginURL, reqMap)
		logger.Error("post err:", err)
	})

	c.OnResponse(func(resp *colly.Response) {
		logger.Info(string(resp.Body))
		// cookie: datr=mrzDXdgeLArnpjRXz98ll3YE; sb=mrzDXRO23kG6ft4c3U37GeAu
		//_ = util.SaveStringToFile("./data", "basic_index.html", string(resp.Body))
		for _, cookie := range c.Cookies("https://mbasic.facebook.com") {
			logger.Info("cookie", cookie.Value)
		}
	})

	c.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
	})

	err = c.Visit(url)

	return nil, err
}
