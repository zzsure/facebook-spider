package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/op/go-logging"
	"gitlab.azbit.cn/web/facebook-spider/conf"
	"gitlab.azbit.cn/web/facebook-spider/consts"
	"gitlab.azbit.cn/web/facebook-spider/library/util"
	"gitlab.azbit.cn/web/facebook-spider/models"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var logger = logging.MustGetLogger("modules/crawler")

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
		//break
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
	// Request the HTML page.
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("new request error:", err)
	}
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	res, err := client.Do(req)
	if err != nil {
		logger.Error("http client do error:", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		logger.Error("status code error:", res.StatusCode, res.Status)
	}

	// TODO: for test save html to data
	//body, err := ioutil.ReadAll(res.Body)
	//_ = util.SaveStringToFile("./data", "res_20191030.html", string(body))

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
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
