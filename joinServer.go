package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	types "game/types"

	"github.com/gorilla/websocket"
)

var conn *websocket.Conn
var connected bool = false
var clientID string

var lastSentTime time.Time

const sendInterval = 20 * time.Millisecond // minimalny odstęp między wysyłaniem wiadomości

func SendToServer(msg types.Message) {
	msg.ID = clientID
	if !connected {
		return
	}
	if conn == nil {
		fmt.Println("Connection is not established")
		return
	}

	// Sprawdź, czy odstęp czasu od ostatniej wiadomości jest wystarczający
	now := time.Now()
	if now.Sub(lastSentTime) < sendInterval {
		fmt.Println("Message throttled, not sending")
		return
	}
	lastSentTime = now

	fmt.Printf("Sending message to server: %+v\n", msg)

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error marshalling message: %v\n", err)
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, msgBytes)
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
	}
}

func HandlerMessages(msgBytes []byte) {
	var msg types.Message
	err := json.Unmarshal(msgBytes, &msg)
	if err != nil {
		fmt.Printf("Error unmarshalling message: %v\n", err)
		return
	}

	if msg.Action == "set_id" {
		clientID = msg.ID
		connected = true
		fmt.Printf("Assigned ID: %s\n", clientID)
	} else if msg.Action == "sync" {
		// Synchronizacja stanu gry
		fmt.Printf("Synchronizing state for player %s: %+v\n", msg.ID, msg)
		playersMap[msg.ID] = types.PlayerState{
			ID:          msg.ID,
			Direction:   msg.Direction,
			X:           msg.X,
			Y:           msg.Y,
			PlayerFrame: msg.PlayerFrame,
			PlayerSrcX:  msg.PlayerSrcX,
			PlayerSrcY:  msg.PlayerSrcY,
		}
	} else if msg.ID != clientID {
		// Debugowanie: Sprawdź, czy dane są prawidłowe
		fmt.Printf("Updating player %s: %+v\n", msg.ID, msg)
		playersMap[msg.ID] = types.PlayerState{
			ID:          msg.ID,
			Direction:   msg.Direction,
			X:           msg.X,
			Y:           msg.Y,
			PlayerFrame: msg.PlayerFrame,
			PlayerSrcX:  msg.PlayerSrcX,
			PlayerSrcY:  msg.PlayerSrcY,
		}
	} else {
		fmt.Printf("Message received: %+v\n", msg)
	}
}

// JoinServer connects to the WebSocket server and starts listening for messages.
func JoinServer(ip string, port int, password string) {
	url := fmt.Sprintf("ws://%s:%d/ws", ip, port)
	fmt.Printf("Connecting to WebSocket server at %s\n", url)

	var err error
	conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to server")
	connected = true

	// Start a goroutine to listen for incoming messages
	go func() {
		for {
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("Error reading message: %v\n", err)
				return
			}
			HandlerMessages(msgBytes)
		}
	}()

	// Send a message to the server every 10 seconds
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			SendToServer(types.Message{Action: "heartbeat"})
		}
	}
}
