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
	Debug        bool   `mapstructure:"debug"`
	LogLevel     string `mapstructure:"log_level"`
	HTTPPort     string `mapstructure:"http_port"`
	RPCPort      string `mapstructure:"rpc_port"`
	DebugPort    string `mapstructure:"debug_port"`
	RegisterAddr string `mapstructure:"register_addr"`
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
		log.Errorln("init config error:", err)
		// panic("init config error")
	}

	err = viper.Unmarshal(&Conf)
	if err != nil {
		log.Errorln("init config unmarshal error:", err)
		// panic("init config unmarshal error")
	}
	log.Println("load config ok", Conf)

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
