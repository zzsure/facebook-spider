package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"gitlab.azbit.cn/web/facebook-spider/controller/response"
)

func Spider(c *gin.Context) {
	co := colly.NewCollector()
	co.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	co.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	co.Visit("http://go-colly.org/")

	response.Response(c, 0, "", nil)
}
