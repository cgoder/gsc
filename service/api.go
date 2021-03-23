package service

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
)

type gsc struct {
}

type Args struct {
	msg GscMsg
}

type Reply struct {
	ret interface{}
}

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
	go HandleTaskMessages()

	// HTTP/WS Server.
	err = startServer()

	return err
}

func (c *gsc) Start(ctx context.Context, args *Args, reply *Reply) error {
	if args.msg.Flag == prefixFlag {
		taskCh <- args.msg.Msg
	} else {
		err := errors.New("gsc unsupport cmd: " + args.msg.Flag)
		log.Errorln(err.Error(), JsonFormat(args.msg))
		return err
	}
	return nil
}

func (c *gsc) Stop(ctx context.Context, args *Args, reply *Reply) error {
	if args.msg.Flag == prefixFlag {
		taskCh <- args.msg.Msg
	} else {
		err := errors.New("gsc unsupport cmd: " + args.msg.Flag)
		log.Errorln(err.Error(), JsonFormat(args.msg))
		return err
	}
	return nil
}

func GetInfo(ctx context.Context, args *Args, reply *Reply) error {
	pb, err := probe(args.msg.Msg.Input)
	if err != nil {
		return err
	}

	reply.ret = pb
	return nil
}
