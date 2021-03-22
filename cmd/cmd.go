package cmd

import (
	"os"

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
	version = "gsc version 0.0.1"
	usage   = `
Usage:
  gsc            Run server.
  gsc version    Print version.
  gsc help       This help text.
	`
)

func ParseArgs() {
	log.Println(banner)

	args := os.Args

	// Use defaults if no args are set.
	if len(args) == 1 {
		return
	}

	// Print version, help or set port.
	if args[1] == "version" || args[1] == "-v" {
		log.Println(version)
		os.Exit(1)
	} else if args[1] == "help" || args[1] == "-h" {
		log.Println(usage)
		os.Exit(1)
	}
}
