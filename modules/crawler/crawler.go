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
func StartBasicCrawlTask(fds []*models.FileData) error {
	// TODO: for test
	html, _ := util.ReadStringFromFile("./data/vogue.html")
	pds, moreURL, err := parsePost([]byte(html))
	logger.Info("show more url: ", moreURL)

	// save article data to file
	err = savePostDataToFile(pds, "https://www.facebook.com/Vogue/")
	if err != nil {
		logger.Error("save article data err:", err)
	}

	_, err = getPostComments(pds)
	return nil
}

// a cron task
func StartBasicCrawlTask1(fds []*models.FileData) error {
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

		// crawl valid date post
		pds, err := recursiveCrawlPost(content)
		if err != nil {
			logger.Error("crawl url: ", fd.URL, ", err:", err)
			continue
		}
		// save post data to file
		err = savePostDataToFile(pds, fd.URL)
		if err != nil {
			logger.Error("save url:", fd.URL, ", post data err:", err)
		}

		// get all post's comments
		cds, err := getPostComments(pds)
		if err != nil {
			logger.Error("crawl url: ", fd.URL, ", err:", err)
			continue
		}
		err = saveCommentDataToFile(cds, fd.URL)
		if err != nil {
			logger.Error("save url:", fd.URL, ", comment data err:", err)
		}

		// TODO: need delete
		break
	}
	return nil
}

func crawlSleep() {
	rs := rand.Intn(consts.MAX_SLEEP_TIME)
	logger.Info("random post sleep seconds:", rs)
	time.Sleep(time.Duration(rs) * time.Second)
}

func getPostComments(pds []*models.PostData) ([]*models.CommentData, error) {
	var acds []*models.CommentData
	for _, pd := range pds {
		if pd.CommentURL != "" {
			/*content, err := crawlByColly(pd.CommentURL)
			if err != nil {
				logger.Error("craw comment url: ", pd.CommentURL, ", err: ", err)
				continue
			}*/
			// TODO: for test
			html, _ := util.ReadStringFromFile("./data/comment.html")

			cds, err := parseComment([]byte(html))
			if err != nil {
				logger.Error("parse comment err:", err)
				continue
			}

			acds = append(acds, cds...)

			// TODO: need delete
			break
		}
	}
	return acds, nil
}

func recursiveCrawlPost(content []byte) ([]*models.PostData, error) {
	pds, moreURL, err := parsePost(content)
	logger.Info("more url:", string(moreURL))
	if err != nil {
		return nil, err
	}

	validDate := util.GetOffsetDateStr(-1 * conf.Config.Spider.RepeatDays)
	for _, pd := range pds {
		if pd.Date < validDate {
			return pds, nil
		} else {
			content, err := crawlByColly(moreURL)
			if err != nil {
				logger.Error("crawl more url: ", moreURL, ", err: ", err)
				return pds, nil
			}
			rpds, err := recursiveCrawlPost(content)
			if err == nil {
				pds = append(pds, rpds...)
			}
		}
	}

	return pds, err
}

func parsePost(b []byte) ([]*models.PostData, string, error) {
	var rets []*models.PostData
	var moreURL string

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(b)))
	if err != nil {
		return nil, "", err
	}
	doc.Find("div[role=article]").Each(func(i int, s *goquery.Selection) {
		post := ""
		s.Find("span p").Each(func(i int, s *goquery.Selection) {
			post += strings.TrimLeft(s.Text(), " ") + "\n"
		})

		logger.Info("post is: ", post)

		cellTime := s.Find("abbr").Text()
		logger.Info("time string is: ", cellTime)
		date := util.GetDateByCellTime(cellTime)
		logger.Info("date string is: ", date)

		var commentURL string
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.Text(), "Comment") {
				if util.IsContainNumber(s.Text()) {
					commentURL, _ = s.Attr("href")
					if commentURL != "" {
						commentURL = fmt.Sprintf("%s%s", strings.TrimRight(consts.BASIC_HTTPS_FACEBOOK_DOMAIN, "/"), commentURL)
					}
				}
				return
			}
		})
		logger.Info("comment url: ", commentURL)
		logger.Info("\n")

		ad := &models.PostData{
			Date:       date,
			Post:       post,
			CommentURL: commentURL,
		}
		rets = append(rets, ad)
	})

	doc.Find(".i a").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "Show more") {
			moreURL, _ = s.Attr("href")
			if moreURL != "" {
				moreURL = fmt.Sprintf("%s%s", strings.TrimRight(consts.BASIC_HTTPS_FACEBOOK_DOMAIN, "/"), moreURL)
			}
			return
		}
	})
	return rets, moreURL, nil
}

func parseComment(content []byte) ([]*models.CommentData, error) {
	return nil, nil
}

// save comment data to file
func saveCommentDataToFile(cds []*models.CommentData, url string) error {
	cm := make(map[string]string)

	for _, cd := range cds {
		if _, ok := cm[cd.Date]; ok {
			cm[cd.Date] += "\n"
		}
		if cd.Comment != "" {
			cm[cd.Date] += cd.Comment
		}
	}

	if len(cm) == 0 {
		logger.Info("len pm and cm is 0")
		return nil
	}

	articleDir, err := util.GetArticleDir(conf.Config.Spider.ArticleBaseDir, url)
	if err != nil {
		logger.Error("get posts path err:", err)
		return err
	}

	// save comment data to file
	for k, v := range cm {
		err = util.SaveStringToFile(articleDir, util.GetCommentFileName(k), v)
		if err != nil {
			logger.Error("save ", k, " posts err:", err)
			continue
		}
	}

	return nil
}

// save article data to file
func savePostDataToFile(pds []*models.PostData, url string) error {
	pm := make(map[string]string)

	for _, pd := range pds {
		if _, ok := pm[pd.Date]; ok {
			pm[pd.Date] += "\n"
		}
		if pd.Post != "" {
			pm[pd.Date] += pd.Post
		}
	}

	if len(pm) == 0 {
		logger.Info("len pm and cm is 0")
		return nil
	}

	postsDir, err := util.GetArticleDir(conf.Config.Spider.ArticleBaseDir, url)
	if err != nil {
		logger.Error("get posts path err:", err)
		return err
	}

	// save post data to file
	for k, v := range pm {
		err = util.SaveStringToFile(postsDir, util.GetPostFileName(k), v)
		if err != nil {
			logger.Error("save ", k, " posts err:", err)
			continue
		}
	}

	return nil
}

// crawl by colly
func crawlByColly(url string) ([]byte, error) {
	crawlSleep()

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
