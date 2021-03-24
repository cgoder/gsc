package config

import (
	"flag"
	"io/ioutil"
	"strings"

	"github.com/fsnotify/fsnotify"
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

var Conf Config

func LoadConfig() {
	flag.Parse()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("Debug", "false")
	viper.SetDefault("LogLevel", "debug")
	viper.SetDefault("HTTPPort", "8084")
	viper.SetDefault("RPCPort", "8085")
	viper.SetDefault("DebugPort", "8086")
	viper.SetDefault("RegisterAddr", "0.0.0.0:2181")

	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		log.Errorln("read config file error:", err)
		log.Errorln("load default config.")
		// panic("init config error")
	}

	err = viper.Unmarshal(&Conf)
	if err != nil {
		log.Errorln("init config unmarshal error:", err)
		// panic("init config unmarshal error")
	} else {
		log.Println("load config ok", Conf)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err != nil {
			log.Errorln("init config unmarshal error:", err)
			// panic("init config unmarshal error")
		}
		log.Println("Config file reload:", e)
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
