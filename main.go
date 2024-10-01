package main

import (
	"flag"
	"fmt"
	"os"

	"game/mapeditor"
	"game/server"
)

func main() {
	helpFlag := flag.Bool("help", false, "Hepl message")
	editorFlag := flag.Bool("editor", false, "Map Editor")
	serverFlag := flag.Bool("server", false, "Start game server")
	joinFlag := flag.Bool("join", false, "Join a game server")
	ip := flag.String("ip", "", "Server IP address")
	port := flag.Int("port", 0, "Server port")
	password := flag.String("password", "", "Server password")

	flag.Parse()

	if *helpFlag {
		fmt.Println("Help Message:")
		fmt.Println("-help -> show help message")
		fmt.Println("-editor -> open map editor")
		fmt.Println("Hot multiplayer game:")
		fmt.Println("-server")
		fmt.Println("Join multiplayer game:")
		fmt.Println("game.exe -join -ip 127.0.0.1 -port 3000 -password password")
		os.Exit(0)
	}

	if *joinFlag {
		if *ip == "" || *port == 0 {
			fmt.Println("Error: -ip and -port must be specified for joining a server.")
			flag.Usage()
			os.Exit(1)
		}
		// Add logic to connect to the server
		fmt.Printf("Joining server at %s:%d with password %s...\n", *ip, *port, *password)
		go JoinServer(*ip, *port, *password)
		StartGame()
		return
	}

	if *editorFlag {
		mapeditor.StartMapEditor()
		return
	}

	if *serverFlag {
		server.StartServer()
		return
	}

	fmt.Print("Starting game...\n")
	StartGame()
}
