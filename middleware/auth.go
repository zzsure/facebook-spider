package middleware

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.azbit.cn/web/facebook-spider/conf"
	"gitlab.azbit.cn/web/facebook-spider/consts"
	"gitlab.azbit.cn/web/facebook-spider/controller/response"
	"io/ioutil"
)

var whitePaths = map[string]struct{}{
	"/health": {},
}

func Auth(c *gin.Context) {
	requestID := c.MustGet(consts.REQUEST_ID_KEY)
	_, ok := whitePaths[c.Request.URL.Path]
	if ok {
		return
	}

	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Warning(requestID, err)
		c.Abort()
		response.ClientErr(c, "鉴权失败")
		return
	}

	if secret, ok := c.GetQuery("secret"); !ok || secret != conf.Config.Auth.Secret {
		apiKey := c.Request.Header.Get("Azbit-Auth-ApiKey")
		sign := c.Request.Header.Get("Azbit-Auth-Sign")
		ts := c.Request.Header.Get("Azbit-Auth-Timestamp")

		// 需要的话，可加上对ts值的校验，例如判断是否为最近1分钟请求

		apiSecret, ok := conf.Config.Auth.Account[apiKey]
		if apiKey == "" || sign == "" || !ok {
			logger.Warning(requestID, fmt.Sprintf("apiKey: %s, sign: %s", apiKey, sign))
			c.Abort()
			response.ClientErr(c, "鉴权失败")
			return
		}

		// 计算签名
		genSign := fmt.Sprintf("%x", sha1.Sum(append(reqBody, []byte(ts+apiSecret)...)))
		if genSign != sign {
			logger.Warning(requestID, "want:", genSign, "get:", sign, string(reqBody), ts)
			c.Abort()
			response.ClientErr(c, "鉴权失败")
			return
		}
	}

	// restore Body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	c.Next()
}
