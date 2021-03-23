package service

import (
	"sync"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	cid  string
	conn *websocket.Conn
}

type GscClients struct {
	m       sync.RWMutex
	clients map[string]*Client
}

var (
	//clients list
	clientsMap = GscClients{clients: make(map[string]*Client)}
)

func ClientsAdd(ws *websocket.Conn) (string, error) {
	clientsMap.m.Lock()
	defer clientsMap.m.Unlock()

	for _, client := range clientsMap.clients {
		if client.conn == ws {
			return "", ErrorClientExisit
		}
	}

	var err error
	cid := uuid.Must(uuid.NewV4(), err).String()
	clientsMap.clients[cid] = &Client{cid: cid, conn: ws}

	return cid, nil
}

func ClientsRemove(cid string) error {
	clientsMap.m.Lock()
	defer clientsMap.m.Unlock()

	if c, ok := clientsMap.clients[cid]; ok {
		c.conn.Close()
		delete(clientsMap.clients, cid)
		return nil
	}

	return ErrorClientNotFound
}

func ClientsGet(cid string) *Client {
	clientsMap.m.Lock()
	defer clientsMap.m.Unlock()

	if c, ok := clientsMap.clients[cid]; ok {
		return c
	}

	return nil
}

func sendInfoClients(stats Status) error {
	clientsMap.m.Lock()
	defer clientsMap.m.Unlock()

	p := &stats

	for cid, client := range clientsMap.clients {
		err := client.conn.WriteJSON(p)
		if err != nil {
			log.Errorln("error: %w", err)
			ClientsRemove(cid)
		}
	}
	return nil
}
