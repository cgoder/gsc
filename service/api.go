package service

import (
	"errors"

	"github.com/cgoder/gsc/ffmpeg"
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
	if msg.Flag == prefixFlag {
		taskCh <- msg.Msg
	} else {
		err := errors.New("gsc unsupport cmd: " + msg.Flag)
		log.Errorln(err.Error(), JsonFormat(msg))
		return err
	}
	return nil
}

func Stop(msg GscMsg) error {
	if msg.Flag == prefixFlag {
		taskCh <- msg.Msg
	} else {
		err := errors.New("gsc unsupport cmd: " + msg.Flag)
		log.Errorln(err.Error(), JsonFormat(msg))
		return err
	}
	return nil
}

func GetInfo(src string) (*ffmpeg.FFProbeResponse, error) {
	return probe(src)
}
