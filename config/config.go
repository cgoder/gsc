package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cgoder/gsc/common"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Debug        bool   `json:"debug"`
	LogLevel     string `json:"log_level"`
	HTTPPort     string `json:"http_port"`
	RPCPort      string `json:"rpc_port"`
	DebugPort    string `json:"debug_port"`
	RegisterAddr string `json:"register_addr"`
}

var (
	Conf      Config
	configMD5 string
)

func LoadConfig() {
	flag.Parse()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("debug", "true")
	viper.SetDefault("log_level", "debug")
	viper.SetDefault("http_port", "8080")
	viper.SetDefault("rpc_port", "8081")
	viper.SetDefault("debug_port", "8082")
	viper.SetDefault("register_addr", "register")

	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("load default config.", err)
		// panic("init config error")
	}

	err = viper.Unmarshal(&Conf, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "json"
	})
	if err != nil {
		fmt.Println("init config unmarshal error:", err)
		// panic("init config unmarshal error")
	} else {
		configMD5 = common.GetFileMd5("config.json")
		fmt.Println("load config ok", Conf, configMD5)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		tmpMD5 := common.GetFileMd5("config.json")
		if tmpMD5 == configMD5 {
			fmt.Println("config file changed, but MD5 same.", tmpMD5)
			return
		}

		var tmpConf Config
		if err := viper.Unmarshal(&tmpConf, func(dc *mapstructure.DecoderConfig) {
			dc.TagName = "json"
		}); err != nil {
			fmt.Println("Config file reload parse error:", err)
			return
		}

		configMD5 = tmpMD5
		Conf = tmpConf
		fmt.Println("Config file reload ok:", Conf, configMD5)

		ReloadConfig(Conf)
	})

	// init log
	level, _ := log.ParseLevel(Conf.LogLevel)
	log.SetLevel(level)
	if !Conf.Debug {
		log.SetOutput(ioutil.Discard)
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: false,
		})
		log.SetReportCaller(false)
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
		log.SetReportCaller(true)
	}

}

func ReloadConfig(conf Config) {
	level, _ := log.ParseLevel(conf.LogLevel)
	log.SetLevel(level)
	if !conf.Debug {
		log.SetOutput(ioutil.Discard)
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: false,
		})
		log.SetReportCaller(false)
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
		log.SetReportCaller(true)
	}
}
