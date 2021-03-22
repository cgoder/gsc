package service

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/cgoder/gsc/ffmpeg"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	prefixCmd    = "gsc"
	prefixStart  = "start"
	prefixStop   = "stop"
	prefixCancel = "cancel"

	//clients list
	clients = make(map[*websocket.Conn]bool)
	//task chan
	broadcast = make(chan Message)

	progressCheckInterval = time.Second * 1
)

func runProcess(input, output, payload string) (string, error) {
	var tid string

	//probe source
	probe := ffmpeg.FFProbe{}
	probeData, err := probe.Run(input)
	if err != nil {
		log.Errorln("ffprobe fail: ", err.Error())
		sendInfoClients(Status{Err: err.Error()})
		return "", err
	}
	log.Debugln(JsonFormat(probeData))

	ffmpeg := &ffmpeg.FFmpeg{}

	//progress
	progressCh := make(chan struct{})
	go trackProgress(context.TODO(), progressCh, probeData, ffmpeg)
	defer close(progressCh)

	// If we get an error back from ffmpeg, send an error ws message to clients.
	err = ffmpeg.Run(context.TODO(), input, output, payload)
	if err != nil {
		log.Errorln(err.Error())
		sendInfoClients(Status{Err: err.Error()})
		return "", err
	}

	sendInfoClients(Status{Percent: 100})
	return tid, nil
}

func trackProgress(ctx context.Context, progress <-chan struct{}, p *ffmpeg.FFProbeResponse, f *ffmpeg.FFmpeg) {
	ticker := time.NewTicker(progressCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Debugln("cancel trackProgress.")
			return
		case <-progress:
			log.Debugln("Waiting for next job...")
			return
		case <-ticker.C:
			currentFrame := f.Progress.Frame
			totalFrames, _ := strconv.Atoi(p.Streams[0].NbFrames)
			speed := f.Progress.Speed
			fps := f.Progress.FPS

			// Only track progress if we know the total frames.
			var pct float64
			if totalFrames != 0 {
				pct = (float64(currentFrame) / float64(totalFrames)) * 100
				pct = math.Round(pct*100) / 100

				log.Debugf("Encoding... %d / %d (%0.2f%%) %s @ %0.2f fps", currentFrame, totalFrames, pct, speed, fps)

			} else {
				pct = f.Progress.Progress
			}
			//write progress to clients
			sendInfoClients(Status{Percent: pct, Speed: speed, FPS: fps})
		}
	}
}

func checkFFmpeg() error {
	f := &ffmpeg.FFmpeg{}
	version, err := f.Version()
	if err != nil {
		return err
	}
	log.Debugln("Checking FFmpeg version....\u001b[32m" + version + "\u001b[0m")

	probe := &ffmpeg.FFProbe{}
	version, err = probe.Version()
	if err != nil {
		return err
	}
	log.Debugln("Checking FFprobe version...\u001b[32m" + version + "\u001b[0m\n")
	return nil
}

func handleMessages() {
	for {
		msg := <-broadcast
		// log.Infoln(JsonFormat(msg))

		if msg.Type == prefixCmd {
			tid, err := runProcess(msg.Input, msg.Output, msg.Payload)
			if err != nil {
				log.Errorln("process fail! ", tid, err.Error())
			}
		} else {
			log.Errorln("unsupport cmd: ", JsonFormat(msg))
		}
	}
}

func sendInfoClients(stas Status) error {
	p := &stas

	for client := range clients {
		err := client.WriteJSON(p)
		if err != nil {
			log.Errorln("error: %w", err)
			client.Close()
			delete(clients, client)
		}
	}
	return nil
}
