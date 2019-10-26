package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/satori/go.uuid"
	"gitlab.azbit.cn/web/facebook-spider/conf"
	"gitlab.azbit.cn/web/facebook-spider/consts"
	"gitlab.azbit.cn/web/facebook-spider/library/log"
	"net/http"
	"os"
	"time"
)

var logger = logging.MustGetLogger("middleware")

func Access(c *gin.Context) {
	if conf.Config.IsProduction() {
		start := time.Now()
		defer func() {
			end := time.Now()
			duration := float64(end.Sub(start).Nanoseconds()) / 1000000 //ms

			var accessLog = map[string]interface{}{
				"idc":       os.Getenv("AZBIT_KUBERNETES_IDC"),
				"index":     "access",
				"timestamp": end.UnixNano() / int64(time.Millisecond),
				"status":    c.Writer.Status(),
				"duration":  duration,
				"clientIp":  c.ClientIP(),
				"method":    c.Request.Method,
				"path":      c.Request.URL.Path,
			}
			accessLogJson, _ := json.Marshal(accessLog)
			_, _ = os.Stdout.Write(append(accessLogJson, '\n'))
		}()
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, conf.Config.Server.MaxHttpRequestBody*1024*1024)
	requestID := uuid.NewV4().String()
	c.Set(consts.REQUEST_ID_KEY, log.RequestID(requestID))
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Next()
}
