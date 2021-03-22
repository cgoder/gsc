package main

import (
	"github.com/cgoder/gsc/cmd"
	"github.com/cgoder/gsc/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	cmd.ParseArgs()

	log.SetLevel(log.DebugLevel)

	//init service
	if err := service.Init(); err != nil {
		log.Errorln(err.Error())
	}

	// //send task cmd
	// start service.Message
	// if _, err := service.Start(start); err != nil {
	// 	log.Errorln(err.Error())
	// }

	// //stop/cancel task cmd
	// stop service.Message
	// if err := service.Stop(stop); err != nil {
	// 	log.Errorln(err.Error())
	// }

}
