package main

import (
	"fmt"
	"math"
	"strconv"

	assetsmanager "game/modules/AssetsManager"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func drawScene(AssetsManager *assetsmanager.AssetsManager) {

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
	//fmt.Printf("Drawing Local player at (%f, %f)\n", playerDest.X, playerDest.Y)

	rl.DrawTexturePro(playerSprite,
		playerSrc, playerDest,
		rl.NewVector2(playerDest.Width/2, playerDest.Height/2),
		0, rl.White,
	)

	// TODO
	/*
		przenieść do do modules/draw
		i tutaj tylko wykonać funkcję np drawWeapons()

			zrobić funkcję w am do obsługi broni + z funkcją attack animation
			+ mają brać poprawkę na obrut gracza, jego pozycję itp
			(rysować razem z graczem w draw)
			(dla multika musi się zgadzać)

			albo użyć animacja atakowania z graczem w 1 spricie

			można podpatrzeć jak to np albion czy coś robi
	*/
	// item in hand ============================================

	// narysuj item w dłoni gracza
	axe_obj, err := AssetsManager.GetAssetObj("items_basic_axe_obj")
	if err != nil {
		fmt.Println("ERROR Drawscene.go : 86 axe_obj")
		return
	}

	cameraX := int32(cam.Target.X)
	cameraY := int32(cam.Target.Y)
	x := cameraX + int32(40)
	y := cameraY + int32(55)

	// oblicz kąt obrotu przedmiotu tak aby był skierowany
	// w stronę kursora
	// Pozycja kamery (celowanie w środek ekranu)
	sc_mid_X := float32(rl.GetScreenWidth()) / 2
	sc_mid_Y := float32(rl.GetScreenHeight()) / 2
	// Pozycja kursora
	mouseX := float32(rl.GetMouseX())
	mouseY := float32(rl.GetMouseY())
	// Oblicz różnicę współrzędnych między środkiem ekranu a kursorem
	deltaX := mouseX - sc_mid_X
	deltaY := mouseY - sc_mid_Y
	// Oblicz kąt w radianach i przekonwertuj na stopnie
	angleRad := math.Atan2(float64(deltaY), float64(deltaX))
	angleDeg := float32(angleRad * (180 / math.Pi))

	// axe_obj.DrawTextureFromData_Idle(rl.Rectangle{})

	item_position := rl.Rectangle{
		X:      float32(x), // Pozycja X na ekranie (obok liczby HP)
		Y:      float32(y), // Pozycja Y na ekranie
		Width:  16 * 1.5,   // Szerokość ikony na ekranie (skaluje do większych wymiarów)
		Height: 16 * 1.5,   // Wysokość ikony na ekranie (skaluje do większych wymiarów)
	}
	axe_obj.Rotation.RotationValue = angleDeg + 100
	axe_obj.Rotation.RotationOrigin.X = item_position.Width / 2
	axe_obj.Rotation.RotationOrigin.Y = item_position.Height / 2
	axe_obj.DrawTextureFromData_Idle(item_position)

	// item in hand ============================================

	// Debugowanie: Sprawdź liczba graczy i ich stan
	//fmt.Printf("Number of players: %d\n", len(playersMap))

	for _, state := range playersMap {
		//fmt.Printf("Player %s: SrcX: %f, SrcY: %f, DestX: %f, DestY: %f\n", id, state.PlayerSrcX, state.PlayerSrcY, state.X, state.Y)

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
