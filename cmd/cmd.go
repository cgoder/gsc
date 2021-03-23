package cmd

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	banner = `

	██████╗ ███████╗ ██████╗
	██╔════╝ ██╔════╝██╔════╝
	██║  ███╗███████╗██║     
	██║   ██║╚════██║██║     
	╚██████╔╝███████║╚██████╗
	 ╚═════╝ ╚══════╝ ╚═════╝
													 							 
	`
)

func ParseArgs() {
	log.Println(banner)
	log.Println("version " + time.Now().Local().String() + "\n\n")

	args := os.Args

	// Use defaults if no args are set.
	if len(args) == 1 {
		return
	}

}
