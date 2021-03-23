package service

import (
	"flag"
	"time"

	"github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"
	rpc "github.com/smallnest/rpcx/server"
	rpcPlugin "github.com/smallnest/rpcx/serverplugin"
)

var (
	serviceName = "service.gsf.gsc"
	basePath    = "/rpc"

	serviceAddr  = flag.String("addr", "localhost:9527", "server address")
	registerAddr = []string{"localhost:2181"}
)

func ServiceRegist() {
	s := rpc.NewServer()

	r := &rpcPlugin.ZooKeeperRegisterPlugin{
		ServiceAddress:   "tcp@" + *serviceAddr,
		ZooKeeperServers: registerAddr,
		BasePath:         basePath,
		Metrics:          metrics.NewRegistry(),
		UpdateInterval:   5 * time.Second,
	}
	err := r.Start()
	if err != nil {
		log.Errorln(err.Error())
	}
	s.Plugins.Add(r)

	s.RegisterName(serviceName, new(gsc), "")
	err = s.Serve("tcp", *serviceAddr)
	if err != nil {
		panic(err)
	}
}
