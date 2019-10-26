package conf

import (
	"github.com/koding/multiconfig"
	"gitlab.azbit.cn/web/facebook-spider/library/util"
	"strings"
	"time"
)

type ConfigTOML struct {
	Server struct {
		Listen             string         `required:"true" flagUsage:"服务监听地址"`
		Env                string         `default:"Pro" flagUsage:"服务运行时环境"`
		MaxHttpRequestBody int64          `default:"4" flagUsage:"最大允许的http请求body，单位M"`
		TimeLocation       *time.Location `flagUsage:"用于time.ParseInLocation"`
	}

	Auth struct {
		Secret  string            `flagUsage:"跳过鉴权的hack"`
		Account map[string]string `flagUsage:"复杂验证，apiKey=>apiSecret"`
	}

	Database struct {
		HostPort     string `required:"true" flagUsage:"数据库连接，eg：tcp(127.0.0.1:3306)"`
		UserPassword string `required:"true" flagUsage:"数据库账号密码"`
		DB           string `required:"true" flagUsage:"数据库"`
		Conn         struct {
			MaxLifeTime int `default:"600" flagUsage:"连接最长存活时间，单位s"`
			MaxIdle     int `default:"10" flagUsage:"最多空闲连接数"`
			MaxOpen     int `default:"80" flagUsage:"最多打开连接数"`
		}
	}

	Log struct {
		Type  string `default:"json" flagUsage:"日志格式，json|raw"`
		Level int    `default:"5" flagUsage:"日志级别：0 CRITICAL, 1 ERROR, 2 WARNING, 3 NOTICE, 4 INFO, 5 DEBUG"`
	} `flagUsage:"服务日志配置"`
}

func (c *ConfigTOML) IsProduction() bool {
	return strings.ToLower(c.Server.Env) == "pro"
}

var Config *ConfigTOML

func Init(tomlPath, args string) {
	var err error
	var loaders = []multiconfig.Loader{
		&multiconfig.TagLoader{},
		&multiconfig.TOMLLoader{Path: tomlPath},
		&multiconfig.EnvironmentLoader{},
	}
	m := multiconfig.DefaultLoader{
		Loader:    multiconfig.MultiLoader(loaders...),
		Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
	}
	Config = new(ConfigTOML)
	m.MustLoad(&Config)

	Config.Server.TimeLocation, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}

	_ = util.PrettyPrint(Config)
}
