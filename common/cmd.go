package common

import (
	"fmt"
	"os"
	"time"
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
	fmt.Println(banner)
	fmt.Println("version " + time.Now().Local().String() + "\n\n")

	args := os.Args

	// Use defaults if no args are set.
	if len(args) == 1 {
		return
	}

}
