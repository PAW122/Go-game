package main

import (
	"fmt"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func drawScene() {

	// Iterowanie po warstwach
	for l := 0; l < len(layers); l++ {
		layer := layers[l]
		for i := 0; i < len(layer); i++ {
			spriteID := layer[i]
			if spriteID == "00" {
				continue
			}

			var texture rl.Texture2D
			var spriteNumber int

			if len(spriteID) > 0 {
				switch spriteID[0] {
				case 'g':
					texture = grassSprite
					tileSrc = grassTileSrc
				case 'c':
					texture = chestSprite
					tileSrc = chestTileSrc
				default:
					texture = grassSprite
					tileSrc = grassTileSrc
				}

				spriteNumber, _ = strconv.Atoi(spriteID[1:3])
			}

			tileDest.X = tileDest.Width * float32(i%mapW)
			tileDest.Y = tileDest.Height * float32(i/mapW)

			tilesPerRow := int(texture.Width / int32(tileSrc.Width))
			tileSrc.X = tileSrc.Width * float32((spriteNumber-1)%tilesPerRow)
			tileSrc.Y = tileSrc.Height * float32((spriteNumber-1)/tilesPerRow)

			rl.DrawTexturePro(texture,
				tileSrc, tileDest,
				rl.NewVector2(tileDest.Width/2, tileDest.Height/2),
				0, rl.White,
			)
		}
	}

	// narysuj postać gracza
	fmt.Printf("Drawing Local player at (%f, %f)\n", playerDest.X, playerDest.Y)
	rl.DrawTexturePro(playerSprite,
		playerSrc, playerDest,
		rl.NewVector2(playerDest.Width/2, playerDest.Height/2),
		0, rl.White,
	)

	// Debugowanie: Sprawdź liczba graczy i ich stan
	fmt.Printf("Number of players: %d\n", len(playersMap))

	for id, state := range playersMap {
		fmt.Printf("Player %s: SrcX: %f, SrcY: %f, DestX: %f, DestY: %f\n", id, state.PlayerSrcX, state.PlayerSrcY, state.X, state.Y)

		// Tworzenie nowego prostokąta dla pozycji gracza z serwera
		newPlayerDest := rl.NewRectangle(float32(state.X), float32(state.Y), playerDest.Width, playerDest.Height)

		newplayerSrc := rl.Rectangle{
			X:      state.PlayerSrcX,
			Y:      state.PlayerSrcY,
			Width:  playerSrc.Width,  // Zdefiniuj szerokość klatki tekstury
			Height: playerSrc.Height, // Zdefiniuj wysokość klatki tekstury
		}

		// Upewnij się, że używasz prawidłowych współrzędnych
		rl.DrawTexturePro(playerSprite,
			newplayerSrc, newPlayerDest,
			rl.NewVector2(newPlayerDest.Width/2, newPlayerDest.Height/2),
			0, rl.White,
		)
	}
}
