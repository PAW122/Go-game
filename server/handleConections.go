package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	types "game/types"

	"github.com/gorilla/websocket"
)

var lastMessageTime = make(map[string]time.Time)

const messageInterval = 20 * time.Millisecond

func HandlerMessages(clientID string, msg types.Message) {
	// Sprawdź, czy odstęp czasu od ostatniej wiadomości jest wystarczający
	now := time.Now()
	if lastTime, ok := lastMessageTime[clientID]; ok {
		if now.Sub(lastTime) < messageInterval {
			fmt.Printf("Message throttled for client %s\n", clientID)
			return
		}
	}
	lastMessageTime[clientID] = now

	// chceck control summ
	xInt := int32(msg.X * 1000) // Scaling factor to preserve precision
	yInt := int32(msg.Y * 1000)
	sum := uint32(xInt*31 + yInt)
	ctrl_sum := sum & 0xFFFFFFFF
	if msg.ControlSum != ctrl_sum {
		return
	}

	// Aktualizowanie stanu gracza
	playerStates[clientID] = types.PlayerState{
		ID:          clientID,
		Direction:   msg.Direction,
		X:           msg.X,
		Y:           msg.Y,
		PlayerFrame: msg.PlayerFrame,
		PlayerSrcX:  msg.PlayerSrcX,
		PlayerSrcY:  msg.PlayerSrcY,
	}

	// Debugowanie: Wyświetl stan wszystkich graczy
	for id, state := range playerStates {
		fmt.Printf("Player %s state: %+v\n", id, state)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()

	clientsMutex.Lock()
	clientIDCounter++
	clientID := fmt.Sprintf("client-%d", clientIDCounter)
	clients[conn] = clientID
	clientsMutex.Unlock()

	fmt.Printf("New client connected: %s\n", clientID)

	// Send the client their unique ID
	idMessage := types.Message{Action: "set_id", ID: clientID}
	idMessageBytes, _ := json.Marshal(idMessage)
	conn.WriteMessage(websocket.TextMessage, idMessageBytes)

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			clientsMutex.Lock()
			delete(clients, conn)
			clientsMutex.Unlock()
			break
		}

		var msg types.Message
		err = json.Unmarshal(msgBytes, &msg)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			continue
		}

		HandlerMessages(clientID, msg)

		// Logowanie przed rozgłoszeniem
		fmt.Printf("Broadcasting message from %s to all clients\n", clientID)

		// Broadcast the message to all connected clients
		for client, id := range clients {
			if id != clientID {
				err := client.WriteMessage(websocket.TextMessage, msgBytes)
				if err != nil {
					fmt.Println("Error sending message:", err)
					client.Close()
					clientsMutex.Lock()
					delete(clients, client)
					clientsMutex.Unlock()
				}
			}
		}
	}
}

func broadcastState() {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client, clientID := range clients {
		state := playerStates[clientID]
		msg := types.Message{
			Action:    "sync",
			ID:        clientID,
			X:         state.X,
			Y:         state.Y,
			Direction: state.Direction,
		}
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf("Error marshalling state for client %s: %v\n", clientID, err)
			continue
		}

		err = client.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			fmt.Printf("Error sending state to client %s: %v\n", clientID, err)
			client.Close()
			delete(clients, client)
		} else {
			fmt.Printf("Successfully sent state to client %s: %s\n", clientID, string(msgBytes))
		}
	}
}
