package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/cgoder/gsc/config"
	"github.com/cgoder/gsc/ffmpeg"
	"github.com/google/gops/agent"
	log "github.com/sirupsen/logrus"
)

var (
	prefixFlag   = "gsc"
	prefixStart  = "start"
	prefixStop   = "stop"
	prefixCancel = "cancel"

	// progressSignal        = make(chan struct{})
	progressCheckInterval = time.Second * 1
)

func probe(input string) (*ffmpeg.FFProbeResponse, error) {
	var pd *ffmpeg.FFProbeResponse
	//probe source
	probe := ffmpeg.FFProbe{}
	pd, err := probe.Run(input)
	if err != nil {
		log.Errorln("ffprobe fail: ", err.Error())
		// metricMap.MetricsSet(tid, TaskMetric{Statu: taskStatusFail, Err: err.Error()})
		return nil, err
	}
	return pd, nil
}

func runProcess(tid, input, output, payload string) error {
	//probe source
	probe := ffmpeg.FFProbe{}
	probeData, err := probe.Run(input)
	if err != nil {
		log.Errorln("ffprobe fail: ", err.Error())
		metricMap.MetricsSet(tid, TaskMetric{ID: tid, Statu: taskStatusFail, Err: err.Error()})
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	ffmpeg := &ffmpeg.FFmpeg{}

	//update status
	taskMap.TaskFFSet(tid, FFInfo{ffmpeg: ffmpeg, ffProbe: *probeData, ctx: ctx, cancel: cancel})
	metricMap.MetricsSet(tid, TaskMetric{ID: tid, Statu: taskStatusDoing})

	//progress
	go trackProgress(ctx, tid, probeData, ffmpeg)

	//ffmpeg process
	go func() {
		// If we get an error back from ffmpeg, send an error ws message to clients.
		err = ffmpeg.Run(ctx, input, output, payload)
		if err != nil {
			log.Errorln(err.Error())

			metricMap.MetricsSet(tid, TaskMetric{ID: tid, Statu: taskStatusFail, Err: err.Error()})
			return
		}

		if !ffmpeg.BeCancel() {
			metricMap.MetricsSet(tid, TaskMetric{ID: tid, Statu: taskStatusDone, Percent: 100})
		}
	}()

	return nil
}

func stopProcess(tid string) {
	if ff, err := taskMap.TaskFFGet(tid); err == nil {
		ff.cancel()
		ff.ffmpeg.Cancel()
		metricMap.MetricsSet(tid, TaskMetric{ID: tid, Statu: taskStatusCancel})
	} else {
		log.Errorln("stop process fail! ", tid)
	}
}

func trackProgress(ctx context.Context, tid string, p *ffmpeg.FFProbeResponse, f *ffmpeg.FFmpeg) {
	ticker := time.NewTicker(progressCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Debugln("cancel trackProgress. tid: ", tid)
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

			//update status
			metricMap.MetricsSet(tid, TaskMetric{ID: tid, Statu: taskStatusDoing, Percent: pct, Speed: speed, FPS: fps})
		}
	}
}

func checkFFmpeg() error {
	f := &ffmpeg.FFmpeg{}
	version, err := f.Version()
	if err != nil {
		return err
	}
	fmt.Println("Checking FFmpeg version....\u001b[32m" + version + "\u001b[0m")

	probe := &ffmpeg.FFProbe{}
	version, err = probe.Version()
	if err != nil {
		return err
	}
	fmt.Println("Checking FFprobe version...\u001b[32m" + version + "\u001b[0m\n")
	return nil
}

func DebugRuntime() {
	if err := agent.Listen(agent.Options{
		Addr:            ":" + config.Conf.DebugPort,
		ShutdownCleanup: true, // automatically closes on os.Interrupt
	}); err != nil {
		log.Errorln(err)
	}
	time.Sleep(time.Minute)
}
