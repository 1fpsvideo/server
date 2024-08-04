package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/debug"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// clients is a nested map that stores WebSocket connections for each session.
	// The outer map uses session IDs as keys, and the inner map uses WebSocket connections as keys.
	// This structure allows for efficient management of multiple clients across different sessions.
	//
	// Ruby equivalent:
	// clients = {
	//   "session1" => {
	//     websocket_conn1 => true,
	//     websocket_conn2 => true
	//   },
	//   "session2" => {
	//     websocket_conn3 => true
	//   }
	// }
	clients = make(map[string]map[*websocket.Conn]bool)
)

func CursorWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionID"]

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		debug.PrintDebug(fmt.Sprintf("WebSocket upgrade failed: %v", err))
		return
	}
	defer conn.Close()

	if clients[sessionID] == nil {
		clients[sessionID] = make(map[*websocket.Conn]bool)
	}
	clients[sessionID][conn] = true

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			debug.PrintDebug(fmt.Sprintf("WebSocket read failed: %v", err))
			delete(clients[sessionID], conn)
			if len(clients[sessionID]) == 0 {
				delete(clients, sessionID)
			}
			break
		}

		var data struct {
			X int `json:"x"`
			Y int `json:"y"`
		}
		err = json.Unmarshal(message, &data)
		if err != nil {
			debug.PrintDebug(fmt.Sprintf("JSON unmarshal failed: %v", err))
			continue
		}

		// Broadcast the cursor position to all connected clients in the same session
		for client := range clients[sessionID] {
			err := client.WriteJSON(data)
			if err != nil {
				debug.PrintDebug(fmt.Sprintf("WebSocket write failed: %v", err))
				client.Close()
				delete(clients[sessionID], client)
				if len(clients[sessionID]) == 0 {
					delete(clients, sessionID)
				}
			}
		}
	}
}
