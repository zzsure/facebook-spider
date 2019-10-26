package response

import (
	"github.com/gin-gonic/gin"
	"gitlab.azbit.cn/web/facebook-spider/consts"
)

func Response(c *gin.Context, code int, msg string, data interface{}) {
	requestID := c.MustGet(consts.REQUEST_ID_KEY)
	c.JSON(200, map[string]interface{}{
		"data":                data,
		"error_no":            code,
		"error_msg":           msg,
		consts.REQUEST_ID_KEY: requestID,
	})
}

func ClientErr(c *gin.Context, msg string) {
	Response(c, 400, msg, nil)
}

func ServerErr(c *gin.Context, msg string) {
	Response(c, 500, msg, nil)
}

func Success(c *gin.Context) {
	Response(c, 0, "success", nil)
}
