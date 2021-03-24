package rpc

import (
	"time"

	"github.com/cgoder/gsc/config"
	"github.com/cgoder/gsc/service"
	"github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"
	rpcx "github.com/smallnest/rpcx/server"
	rpcxPlugin "github.com/smallnest/rpcx/serverplugin"
)

var (
	serviceName = "service.gsf.gsc"
	basePath    = "/rpc"

	// registerAddr = []string{"localhost:2181"}
)

func ServiceRegist() {
	s := rpcx.NewServer()

	r := &rpcxPlugin.ZooKeeperRegisterPlugin{
		ServiceAddress:   "tcp@" + config.Conf.RPCPort,
		ZooKeeperServers: []string{config.Conf.RegisterAddr},
		BasePath:         basePath,
		Metrics:          metrics.NewRegistry(),
		UpdateInterval:   5 * time.Second,
	}
	err := r.Start()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	s.Plugins.Add(r)

	s.RegisterName(serviceName, new(service.Gsc), "")
	err = s.Serve("tcp", ":"+config.Conf.RPCPort)
	if err != nil {
		panic(err)
	}
}
