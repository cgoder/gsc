package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/cgoder/gsc/config"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (

	//websocket
	allowedOrigins = []string{
		"http://localhost:" + config.Conf.HTTPPort,
	}
	upgrader = websocket.Upgrader{
		// CheckOrigin: func(r *http.Request) bool {
		// 	for _, origin := range allowedOrigins {
		// 		if r.Header.Get("Origin") == origin {
		// 			return true
		// 		}
		// 	}
		// 	return false
		// },
	}
)

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

func startServer() error {
	http.HandleFunc("/ws", handleConnections)

	http.HandleFunc("/files", handleFiles)
	http.Handle("/", http.FileServer(http.Dir("./")))

	log.Println("Wait websocket to connect on Port: ", config.Conf.HTTPPort)
	log.Println("Waiting for connection...")

	err := http.ListenAndServe(":"+config.Conf.HTTPPort, nil)
	if err != nil {
		log.Errorln("ListenAndServe: ", err)
		return err
	}
	return nil
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorln("websocket connection failed!", err.Error())
		return
	}
	defer ws.Close()

	// Register client.
	cid, err := ClientsAdd(ws)
	if err != nil {
		log.Debugln("client regist fail! ", err.Error())
		return
	}
	defer ClientsRemove(cid)

	log.Debugln("Client connected! cid: ", cid)
	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object.
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Debugln("Client disconnect! cid: ", cid)
			break
		}
		// log.Debugln(JsonFormat(msg))
		// Send the newly received message to the task channel.
		taskCh <- msg
	}
}

func handleFiles(w http.ResponseWriter, r *http.Request) {
	log.Debugln(r.URL)
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
