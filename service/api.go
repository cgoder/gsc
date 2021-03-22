package service

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

func Init() error {
	var err error
	// Check if FFmpeg/FFprobe are available.
	err = checkFFmpeg()
	if err != nil {
		log.Errorln(err.Error())
		log.Errorln("Please install FFmpeg and FFprobe on $PATH.")
		return err
	}

	// Handles incoming WS messages from client.
	go handleMessages()

	// HTTP/WS Server.
	err = startServer()

	return err
}

func Start(msg GscMsg) error {
	if msg.Cmd == prefixStart {
		broadcast <- msg.Msg
	} else {
		err := errors.New("gsc unsupport cmd: " + msg.Cmd)
		log.Errorln(err.Error(), JsonFormat(msg))
		return err
	}
	return nil
}

func Stop(msg GscMsg) error {
	if msg.Cmd == prefixStop {
		broadcast <- msg.Msg
	} else {
		err := errors.New("gsc unsupport cmd: " + msg.Cmd)
		log.Errorln(err.Error(), JsonFormat(msg))
		return err
	}
	return nil
}
