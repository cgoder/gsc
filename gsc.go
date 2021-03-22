package main

import (
	"github.com/cgoder/gsc/cmd"
	"github.com/cgoder/gsc/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	cmd.ParseArgs()

	if err := service.Run(); err != nil {
		log.Errorln(err.Error())
	}

}
