package service

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cgoder/gsc/ffmpeg"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan Message)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			for _, origin := range allowedOrigins {
				if r.Header.Get("Origin") == origin {
					return true
				}
			}
			return false
		},
	}
	progressCh chan struct{}

	progressCheckInterval = time.Second * 1
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorln("websocket connection failed!")
		return
	}
	defer ws.Close()

	// Register client.
	clients[ws] = true

	for {
		log.Debugln("connected!")
		var msg Message
		// Read in a new message as JSON and map it to a Message object.
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Debugln("disconnected!")
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel.
		broadcast <- msg
	}
}

func handleFiles(w http.ResponseWriter, r *http.Request) {
	prefix := r.URL.Query().Get("prefix")
	if prefix == "" {
		prefix = "."
	}
	prefix = strings.TrimSuffix(prefix, "/")

	wd, _ := os.Getwd()
	resp := &FilesResponse{
		Cwd:     wd,
		Folders: []string{},
		Files:   []file{},
	}

	files, _ := ioutil.ReadDir(prefix)
	for _, f := range files {
		if f.IsDir() {
			if prefix == "." {
				resp.Folders = append(resp.Folders, f.Name()+"/")
			} else {
				resp.Folders = append(resp.Folders, prefix+"/"+f.Name()+"/")
			}
		} else {
			var obj file
			if prefix == "./" {
				obj.Name = prefix + f.Name()
			} else {
				obj.Name = prefix + "/" + f.Name()
			}
			obj.Size = f.Size()
			resp.Files = append(resp.Files, obj)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	cors(&w, r)
	json.NewEncoder(w).Encode(resp)
}

func cors(w *http.ResponseWriter, r *http.Request) {
	for _, origin := range allowedOrigins {
		if r.Header.Get("Origin") == origin {
			(*w).Header().Set("Access-Control-Allow-Origin", origin)
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		if msg.Type == "encode" {
			runEncode(msg.Input, msg.Output, msg.Payload)
		}
	}
}

func runEncode(input, output, payload string) {
	probe := ffmpeg.FFProbe{}
	probeData, err := probe.Run(input)
	if err != nil {
		sendError(err)
		return
	}

	ffmpeg := &ffmpeg.FFmpeg{}
	go trackEncodeProgress(probeData, ffmpeg)
	err = ffmpeg.Run(input, output, payload)

	// If we get an error back from ffmpeg, send an error ws message to clients.
	if err != nil {
		close(progressCh)
		sendError(err)
		return
	}
	close(progressCh)

	for client := range clients {
		p := &Status{
			Percent: 100,
		}
		err := client.WriteJSON(p)
		if err != nil {
			log.Errorln("error: %w", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func trackEncodeProgress(p *ffmpeg.FFProbeResponse, f *ffmpeg.FFmpeg) {
	progressCh = make(chan struct{})
	ticker := time.NewTicker(progressCheckInterval)

	for {
		select {
		case <-progressCh:
			ticker.Stop()
			log.Debugln("Waiting for next job...")
			return
		case <-ticker.C:
			currentFrame := f.Progress.Frame
			totalFrames, _ := strconv.Atoi(p.Streams[0].NbFrames)
			speed := f.Progress.Speed
			fps := f.Progress.FPS

			// Only track progress if we know the total frames.
			if totalFrames != 0 {
				pct := (float64(currentFrame) / float64(totalFrames)) * 100
				pct = math.Round(pct*100) / 100

				log.Debugf("Encoding... %d / %d (%0.2f%%) %s @ %0.2f fps", currentFrame, totalFrames, pct, speed, fps)

				for client := range clients {
					p := &Status{
						Percent: pct,
						Speed:   speed,
						FPS:     fps,
					}
					err := client.WriteJSON(p)
					if err != nil {
						log.Errorln("error: %w", err)
						client.Close()
						delete(clients, client)
					}
				}
			}
		}
	}
}

func sendError(err error) {
	for client := range clients {
		p := &Status{
			Err: err.Error(),
		}
		err := client.WriteJSON(p)
		if err != nil {
			log.Errorln("error: %w", err)
			client.Close()
			delete(clients, client)
		}
	}
}
