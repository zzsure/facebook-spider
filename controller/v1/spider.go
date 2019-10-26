package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

func Spider(c *gin.Context) {
    colly := colly.NewCollector()
    colly.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	colly.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://go-colly.org/")

	response.Response(c, 0, "", nil)
}
