package server

import (
	"encoding/json"
	"fmt"
	"game/types"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients         = make(map[*websocket.Conn]string)
	clientsMutex    sync.Mutex
	clientIDCounter int
	playerStates    = make(map[string]types.PlayerState) // Mapa przechowująca stan każdego gracza
	upgrader        = websocket.Upgrader{}
)

func StartServer() {
	config, config_err := GetConfig()
	if config_err != nil {
		fmt.Print("Config error")
		os.Exit(1)
	}

	http.HandleFunc("/ws", handleConnections)

	serverAddress := fmt.Sprintf(":%d", config.Port)
	fmt.Printf("WebSocket server started on %s\n", serverAddress)
	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		fmt.Printf("Failed to start WebSocket server: %v\n", err)
	}
}

func GetConfig() (*types.ServerConfig, error) {
	const configFileName = "game_config.json"

	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		// default config
		config := types.ServerConfig{
			Port:       3000,
			Password:   "password",
			MaxPlayers: 10,
		}

		configData, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			fmt.Printf("Błąd podczas serializacji konfiguracji: %v\n", err)
			return nil, fmt.Errorf("Błąd podczas serializacji konfiguracji")
		}

		// Tworzenie pliku i zapisanie do niego konfiguracji
		err = os.WriteFile(configFileName, configData, 0644)
		if err != nil {
			fmt.Printf("Błąd podczas tworzenia pliku config.json: %v\n", err)
			return nil, fmt.Errorf("Błąd podczas tworzenia pliku config.json")
		}

		return &config, nil
	} else {
		data, err := os.ReadFile(configFileName)
		if err != nil {
			return nil, fmt.Errorf("Błąd podczas odczytu pliku config.json: %v\n", err)
		}

		// Deserializacja danych JSON do struktury Config
		var config types.ServerConfig
		err = json.Unmarshal(data, &config)
		if err != nil {
			return nil, fmt.Errorf("Błąd podczas deserializacji danych z config.json: %v\n", err)
		}

		return &config, nil
	}
}
