package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func loadMap(mapFile string) {
	file, err := os.ReadFile(mapFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lines := strings.Split(string(file), "\n")

	header := strings.Fields(lines[0])
	if len(header) < 2 {
		fmt.Println("Invalid map file header")
		os.Exit(1)
	}
	mapW, _ = strconv.Atoi(header[0])
	mapH, _ = strconv.Atoi(header[1])

	fmt.Println("Map dimensions:", mapW, mapH)

	layers = [][]string{}
	layerCount := (len(lines)-1)/mapH - 1 // Liczba warstw wizualnych, ostatnia jest warstwą kolizji

	// Wczytywanie warstw spritów
	for l := 0; l < layerCount; l++ {
		layer := []string{}
		for y := 0; y < mapH; y++ {
			line := lines[1+y+(l*mapH)]
			tiles := strings.Fields(line)
			layer = append(layer, tiles...)
		}
		layers = append(layers, layer)
	}

	// Wczytywanie warstwy kolizji
	collisionLayer = []string{}
	collisionStartLine := 1 + layerCount*mapH
	for y := 0; y < mapH; y++ {
		line := lines[collisionStartLine+y]
		tiles := strings.Fields(line)
		collisionLayer = append(collisionLayer, tiles...)
	}

	fmt.Println("Layers loaded:", layers)
	fmt.Println("Collision layer loaded:", collisionLayer)
}
