package service

import (
	"net/http"

	"github.com/cgoder/gsc/ffmpeg"
	log "github.com/sirupsen/logrus"
)

var (
	port           = "8080"
	allowedOrigins = []string{
		"http://localhost:" + port,
	}
)

// Message payload from client.
type Message struct {
	Type    string `json:"type"`
	Input   string `json:"input"`
	Output  string `json:"output"`
	Payload string `json:"payload"`
}

// Status response to client.
type Status struct {
	Percent float64 `json:"percent"`
	Speed   string  `json:"speed"`
	FPS     float64 `json:"fps"`
	Err     string  `json:"err,omitempty"`
}

// FilesResponse http response for files endpoint.
type FilesResponse struct {
	Cwd     string   `json:"cwd"`
	Folders []string `json:"folders"`
	Files   []file   `json:"files"`
}

type file struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func Run() error {
	// Check if FFmpeg/FFprobe are available.
	err := checkFFmpeg()
	if err != nil {
		log.Errorln(err.Error())
		log.Errorln("Please install FFmpeg and FFprobe on $PATH.")
		return err
	}

	// HTTP/WS Server.
	return startServer()
}

func startServer() error {
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/files", handleFiles)
	http.Handle("/", http.FileServer(http.Dir("./")))

	// Handles incoming WS messages from client.
	go handleMessages()

	log.Println("Wait websocket to connect on port: ", port)
	log.Println("Waiting for connection...")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Errorln("ListenAndServe: ", err)
		return err
	}
	return nil
}

func checkFFmpeg() error {
	f := &ffmpeg.FFmpeg{}
	version, err := f.Version()
	if err != nil {
		return err
	}
	log.Println("  Checking FFmpeg version....\u001b[32m" + version + "\u001b[0m")

	probe := &ffmpeg.FFProbe{}
	version, err = probe.Version()
	if err != nil {
		return err
	}
	log.Println("  Checking FFprobe version...\u001b[32m" + version + "\u001b[0m\n")
	return nil
}
