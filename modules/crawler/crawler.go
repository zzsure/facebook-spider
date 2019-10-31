package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/op/go-logging"
	"gitlab.azbit.cn/web/facebook-spider/conf"
	"gitlab.azbit.cn/web/facebook-spider/library/util"
	"gitlab.azbit.cn/web/facebook-spider/models"
	"strings"
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
		postsDir, err := util.GetPostsDir(conf.Config.Spider.ArticleBaseDir, fd.URL)
		if err != nil {
			logger.Error("get posts path err:", err)
			continue
		}
		err = util.SaveStringToFile(postsDir, util.GetPostFileName(), posts)
		if err != nil {
			logger.Error("save posts err:", err)
		}

		break
	}
}

// craw data by goquery
func crawlByGoquery(url, lang string) ([]*models.ArticleData, error) {
	// Request the HTML page.
	/*res, err := http.Get(url)
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
	}*/

	// TODO: for test
	html, _ := util.ReadStringFromFile("./data/res.html")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		logger.Error("document reader error:", err)
	}

	var rets []*models.ArticleData
	// Find the review items
	doc.Find("._4-u2 .userContent").Each(func(i int, s *goquery.Selection) {
		posts := ""
		var comments []string
		s.Find("p").Each(func(i int, s *goquery.Selection) {
			posts += strings.TrimLeft(s.Text(), " ") + "\n"
		})
		logger.Info("idx: ", i, "ret: ", posts)
		ad := &models.ArticleData{
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
