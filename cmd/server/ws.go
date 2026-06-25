package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// writeWait caps how long a single client write may block, so one slow/dead
// client can't stall broadcasts to everyone else.
const writeWait = 3 * time.Second

// wsHub fans out live messages to all connected WebSocket clients.
type wsHub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]bool
	writeMu sync.Mutex // gorilla forbids concurrent writes to a conn
}

func newWSHub() *wsHub {
	return &wsHub{clients: make(map[*websocket.Conn]bool)}
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // dev: allow all origins
}

// handle upgrades an HTTP request to a WebSocket and registers the client.
func (h *wsHub) handle(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	h.mu.Lock()
	h.clients[conn] = true
	n := len(h.clients)
	h.mu.Unlock()
	log.Printf("[WS] client connected (%d total)", n)

	// Greet, then read until the client disconnects.
	h.writeOne(conn, []byte(`{"type":"hello","data":"connected"}`))
	go func() {
		defer func() {
			h.mu.Lock()
			delete(h.clients, conn)
			h.mu.Unlock()
			conn.Close()
			log.Printf("[WS] client disconnected")
		}()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
}

func (h *wsHub) writeOne(conn *websocket.Conn, payload []byte) {
	h.writeMu.Lock()
	defer h.writeMu.Unlock()
	_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		conn.Close()
	}
}

// broadcast sends a {type,data} message to every connected client.
func (h *wsHub) broadcast(msgType string, data interface{}) {
	payload, err := json.Marshal(map[string]interface{}{"type": msgType, "data": data})
	if err != nil {
		return
	}
	h.mu.RLock()
	conns := make([]*websocket.Conn, 0, len(h.clients))
	for c := range h.clients {
		conns = append(conns, c)
	}
	h.mu.RUnlock()

	h.writeMu.Lock()
	defer h.writeMu.Unlock()
	for _, c := range conns {
		_ = c.SetWriteDeadline(time.Now().Add(writeWait))
		if err := c.WriteMessage(websocket.TextMessage, payload); err != nil {
			c.Close()
			h.mu.Lock()
			delete(h.clients, c)
			h.mu.Unlock()
		}
	}
}
