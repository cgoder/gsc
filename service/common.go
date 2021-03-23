package service

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/google/gops/agent"
	log "github.com/sirupsen/logrus"
)

//JsonFormat Json outupt.
func JsonFormat(v interface{}) string {
	// if msg == "" {
	// 	return ""
	// }
	var out bytes.Buffer

	bs, _ := json.Marshal(v)
	json.Indent(&out, bs, "", "\t")

	return out.String()
}

var pprofAddr = "0.0.0.0:8048"

func DebugRuntime() {
	if err := agent.Listen(agent.Options{
		Addr:            pprofAddr,
		ShutdownCleanup: true, // automatically closes on os.Interrupt
	}); err != nil {
		log.Errorln(err)
	}
	time.Sleep(time.Minute)
}
