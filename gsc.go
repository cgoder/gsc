package main

import (
	"github.com/cgoder/gsc/cmd"
	"github.com/cgoder/gsc/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	cmd.ParseArgs()

	//pprof
	go service.DebugRuntime()

	//init server
	go service.Init()

	//regist and server rpc service
	service.ServiceRegist()

}
