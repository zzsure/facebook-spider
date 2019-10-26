package server

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"gitlab.azbit.cn/web/facebook-spider/conf"
	"gitlab.azbit.cn/web/facebook-spider/controller/v1"
	"gitlab.azbit.cn/web/facebook-spider/library/log"
	"gitlab.azbit.cn/web/facebook-spider/middleware"
)

var Server = cli.Command{
	Name:  "server",
	Usage: "golang_framework http server",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "conf, c",
			Value: "config.toml",
			Usage: "toml配置文件",
		},
		cli.StringFlag{
			Name:  "args",
			Value: "",
			Usage: "multi config cmd line args",
		},
	},
	Action: run,
}

func run(c *cli.Context) {
	conf.Init(c.String("conf"), c.String("args"))
	log.Init()
	//db.Init()

	_ = GinEngine().Run(conf.Config.Server.Listen)
}

func GinEngine() *gin.Engine {
	var r *gin.Engine
	if conf.Config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
		r.Use(middleware.Recovery)
	} else {
		r = gin.Default()
	}
    r.Use(middleware.Access)
	//r.Use(middleware.Auth)
	r.GET("/health")
	V1(r)

	return r
}

func V1(r *gin.Engine) {
	g := r.Group("/v1")
	{
		g.POST("/echo", v1.Echo)
        g.GET("/spider", v1.Spider)
	}
}
