package mgconf

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/sadlil/gologger"
	"strings"
	"time"
)

var conf *koanf.Koanf
var logger = gologger.GetLogger()

const config_file = "./application.yml"
const AUTO_CHECK_MINUTES = 30 //自动检查连接间隔时间，单位为分钟

func InitConfig(cf string) {
	if cf == "" {
		cf = config_file
	}
	logger.Debug("读取配置文件:" + cf)
	conf = koanf.New(".")
	f := file.Provider(cf)
	err := conf.Load(f, yaml.Parser())
	if err != nil {
		logger.Error("读取配置文件错误:" + err.Error())
	}

	configs := conf.String("go.config.used")

	if strings.Contains(configs, "mysql") {
		logger.Info("正在连接MySQL")
		mySqlInit()
	}
	if strings.Contains(configs, "mongodb") {
		logger.Info("正在连接MongoDB")
		mgoInit()
	}
	if strings.Contains(configs, "redis") {
		logger.Info("正在连接Redis")
		redisInit()
	}
	if strings.Contains(configs, "nacos") {
		logger.Info("正在注册到Nacos")
		registerNacos()
	}
	if strings.Contains(configs, "rabbitmq") {
		logger.Info("正在连接RabbitMQ")
		rabbitMQInit()
	}

	//设置定时任务自动检查
	ticker := time.NewTicker(time.Minute * AUTO_CHECK_MINUTES)
	go func() {
		for _ = range ticker.C {
			checkAll()
		}
	}()

}

func GetConfigString(name string) string {
	if conf == nil {
		return ""
	}
	if conf.Exists(name) {
		return conf.String(name)
	} else {
		return ""
	}
}

func GetConfigInt(name string) int {
	if conf == nil {
		return 0
	}
	if conf.Exists(name) {
		return conf.Int(name)
	} else {
		return 0
	}
}

func SafeExit() {

	configs := conf.String("go.config.used")

	if strings.Contains(configs, "mysql") {
		logger.Info("正在关闭MySQL连接")
		mySqlClose()
	}
	if strings.Contains(configs, "mongodb") {
		logger.Info("正在关闭MongoDB连接")
		mgoClose()
	}
	if strings.Contains(configs, "redis") {
		logger.Info("正在关闭Redis连接")
		redisClose()
	}
	if strings.Contains(configs, "nacos") {
		logger.Info("正在注销Nacos")
		deRegisterNacos()
	}
	if strings.Contains(configs, "rabbitmq") {
		logger.Info("正在关闭RabbitMQ连接")
		rabbitMQClose()
	}

}

func checkAll() {

	configs := conf.String("go.config.used")

	if strings.Contains(configs, "mysql") {
		logger.Debug("正在检查MySQL")
		mySqlCheck()
	}
	if strings.Contains(configs, "mongodb") {
		logger.Debug("正在检查MongoDB")
		MgoCheck()
	}
	if strings.Contains(configs, "redis") {
		logger.Debug("正在检查Redis")
		RedisCheck()
	}
}

func getConfigUrl(prefix string) string {
	serverType := conf.String("go.config.server_type")
	configUrl := conf.String("go.config.server")
	switch serverType {
	case "nacos":
		configUrl = configUrl + "nacos/v1/cs/configs?group=DEFAULT_GROUP&dataId=" + prefix + conf.String("go.config.mid") + conf.String("go.config.env") + conf.String("go.config.type")
	case "consul":
		configUrl = configUrl + "v1/kv/" + prefix + conf.String("go.config.mid") + conf.String("go.config.env") + conf.String("go.config.type") + "?dc=dc1&raw=true"
	case "springconfig":
		configUrl = configUrl + prefix + conf.String("go.config.mid") + conf.String("go.config.env") + conf.String("go.config.type")
	default:
		configUrl = configUrl + prefix + conf.String("go.config.mid") + conf.String("go.config.env") + conf.String("go.config.type")
	}
	return configUrl
}
