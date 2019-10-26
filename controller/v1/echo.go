package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/op/go-logging"
	"gitlab.azbit.cn/web/facebook-spider/consts"
	"gitlab.azbit.cn/web/facebook-spider/controller/request"
	"gitlab.azbit.cn/web/facebook-spider/controller/response"
	"math/rand"
	"time"
)

var logger = logging.MustGetLogger("controller/v1")

func Echo(c *gin.Context) {
	requestID := c.MustGet(consts.REQUEST_ID_KEY)
	var input request.Echo
	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		logger.Error(requestID, err)
		response.ClientErr(c, err.Error())
		return
	}

	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	response.Response(c, 0, "", input.Data)
}
