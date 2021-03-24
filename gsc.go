package main

import (
	"github.com/cgoder/gsc/common"
	"github.com/cgoder/gsc/config"
	"github.com/cgoder/gsc/rpc"
	"github.com/cgoder/gsc/service"
)

func main() {
	config.LoadConfig()

	common.ParseArgs()

	//regist and server rpc service
	go rpc.ServiceRegist()

	//pprof
	go service.DebugRuntime()
	//init server
	service.Init()

}
