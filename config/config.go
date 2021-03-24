package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cgoder/gsc/common"
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

var (
	Conf      Config
	configMD5 string
)

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
	viper.AddConfigPath("/etc/")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("read config file error:", err)
		fmt.Println("load default config.")
		// panic("init config error")
	}

	err = viper.Unmarshal(&Conf)
	if err != nil {
		fmt.Println("init config unmarshal error:", err)
		// panic("init config unmarshal error")
	} else {
		fmt.Println("load config ok", Conf)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		tmpMD5 := common.GetFileMd5("/etc/config.json")
		if tmpMD5 == configMD5 {
			fmt.Println("config file changed, but MD5 same.")
			return
		}

		var tmpConf Config
		if err := viper.Unmarshal(&tmpConf); err != nil {
			fmt.Println("Config file reload parse error:", err)
			return
		}

		configMD5 = tmpMD5
		Conf = tmpConf
		fmt.Println("Config file reload ok:", Conf)

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
