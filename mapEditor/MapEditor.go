package mapeditor

import (
	"fmt"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth      = 800
	screenHeight     = 600
	gridSize         = 16
	menuWidth        = 300
	tileSize         = 16 // Rozmiar pojedynczego kafelka w teksturze
	tilePadding      = 5  // Odstęp między kafelkami w menu
	inputFieldWidth  = 50
	inputFieldHeight = 20
)

var (
	grassSprite        rl.Texture2D
	tiles              []rl.Rectangle
	selectedTile       int = -1
	tileMap            [][]int
	mapWidth           int   = 10
	mapHeight          int   = 10
	isSettingsMenuOpen bool  = false
	settingsMenuX      int32 = screenWidth - menuWidth + 10
	settingsMenuY      int32 = 50
	settingsMenuWidth  int32 = 200
	settingsMenuHeight int32 = 100
)

func StartMapEditor() {
	rl.InitWindow(screenWidth, screenHeight, "Map Editor")
	defer rl.CloseWindow()

	// Wczytaj teksturę
	grassSprite = rl.LoadTexture("assets/Tilesets/Grass.png")
	defer rl.UnloadTexture(grassSprite)

	// Oblicz prostokąty dla każdego kafelka w teksturze
	tiles = extractTiles(grassSprite)

	// Inicjalizacja mapy kafelków
	resetTileMap()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		handleInput()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		drawEditor()
		drawMenu()

		if isSettingsMenuOpen {
			drawSettingsMenu()
		}

		rl.EndDrawing()
	}
}

func handleInput() {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		mouseX, mouseY := rl.GetMouseX(), rl.GetMouseY()

		if mouseX >= screenWidth-menuWidth {
			// Sprawdzenie kliknięcia w menu kafelków
			for i, _ := range tiles {
				x := screenWidth - menuWidth + (i%5)*(tileSize+tilePadding)
				y := 20 + (i/5)*(tileSize+tilePadding)
				if mouseX >= int32(x) && mouseX <= int32(x)+tileSize && mouseY >= int32(y) && mouseY <= int32(y)+tileSize {
					selectedTile = i
					break
				}
			}

			if mouseX >= screenWidth-menuWidth+10 && mouseX <= screenWidth-menuWidth+100 {
				if mouseY >= screenHeight-40 && mouseY <= screenHeight-20 {
					isSettingsMenuOpen = !isSettingsMenuOpen
				}
			}
		} else if isSettingsMenuOpen {
			if mouseX >= settingsMenuX && mouseX <= settingsMenuX+settingsMenuWidth {
				if mouseY >= settingsMenuY && mouseY <= settingsMenuY+inputFieldHeight {
					widthStr := rl.GetClipboardText()
					mapWidth, _ = strconv.Atoi(widthStr)
					resetTileMap()
				}
				if mouseY >= settingsMenuY+30 && mouseY <= settingsMenuY+30+inputFieldHeight {
					heightStr := rl.GetClipboardText()
					mapHeight, _ = strconv.Atoi(heightStr)
					resetTileMap()
				}
				if mouseY >= settingsMenuY+60 && mouseY <= settingsMenuY+60+inputFieldHeight {
					isSettingsMenuOpen = false
				}
			}
		} else {
			// Sprawdzenie kliknięcia w edytorze
			if selectedTile != -1 {
				gridX := (mouseX / gridSize) * gridSize
				gridY := (mouseY / gridSize) * gridSize
				mapX := int(gridX / gridSize)
				mapY := int(gridY / gridSize)
				if mapX < len(tileMap) && mapY < len(tileMap[0]) {
					tileMap[mapX][mapY] = selectedTile
				}
			}
		}
	}
}

func resetTileMap() {
	tileMap = make([][]int, mapWidth)
	for i := range tileMap {
		tileMap[i] = make([]int, mapHeight)
	}
}

func extractTiles(texture rl.Texture2D) []rl.Rectangle {
	var tiles []rl.Rectangle

	tileWidth := float32(tileSize)
	tileHeight := float32(tileSize)

	tilesAcross := int(texture.Width / int32(tileWidth))
	tilesDown := int(texture.Height / int32(tileHeight))

	for y := 0; y < tilesDown; y++ {
		for x := 0; x < tilesAcross; x++ {
			tiles = append(tiles, rl.Rectangle{
				X:      float32(x) * tileWidth,
				Y:      float32(y) * tileHeight,
				Width:  tileWidth,
				Height: tileHeight,
			})
		}
	}

	return tiles
}

func drawEditor() {
	// Rysowanie siatki w polu edytora
	for x := 0; x < (screenWidth - menuWidth); x += gridSize {
		for y := 0; y < screenHeight; y += gridSize {
			rl.DrawRectangleLines(int32(x), int32(y), gridSize, gridSize, rl.LightGray)

			mapX := x / gridSize
			mapY := y / gridSize

			if mapX < len(tileMap) && mapY < len(tileMap[0]) {
				tileIndex := tileMap[mapX][mapY]
				if tileIndex >= 0 && tileIndex < len(tiles) {
					rl.DrawTextureRec(grassSprite, tiles[tileIndex], rl.Vector2{X: float32(x), Y: float32(y)}, rl.White)
				}
			}
		}
	}

	// Rysowanie separatora
	rl.DrawLine(screenWidth-menuWidth, 0, screenWidth-menuWidth, screenHeight, rl.Black)
}

func drawMenu() {
	// Rysowanie menu po prawej stronie
	rl.DrawText("Menu", screenWidth-menuWidth+10, 10, 20, rl.Black)
	rl.DrawRectangle(screenWidth-menuWidth, 0, menuWidth, screenHeight, rl.LightGray)

	for i, tile := range tiles {
		x := screenWidth - menuWidth + (i%5)*(tileSize+tilePadding)
		y := 20 + (i/5)*(tileSize+tilePadding)

		rl.DrawTextureRec(grassSprite, tile, rl.Vector2{X: float32(x), Y: float32(y)}, rl.White)

		if i == selectedTile {
			rl.DrawRectangleLines(int32(x), int32(y), tileSize, tileSize, rl.Red)
		} else {
			rl.DrawRectangleLines(int32(x), int32(y), tileSize, tileSize, rl.Black)
		}
	}

	// Rysowanie przycisku do ustawień
	settingsButtonX := int32(screenWidth) - menuWidth + 10
	settingsButtonY := int32(screenHeight) - 30
	settingsButtonWidth := int32(100)
	settingsButtonHeight := int32(20)

	rl.DrawRectangle(settingsButtonX, settingsButtonY, settingsButtonWidth, settingsButtonHeight, rl.Gray)
	rl.DrawText("Settings", settingsButtonX+5, settingsButtonY+5, 20, rl.Black)
}

func drawSettingsMenu() {
	// Rysowanie menu ustawień
	rl.DrawRectangle(settingsMenuX, settingsMenuY, settingsMenuWidth, settingsMenuHeight, rl.LightGray)
	rl.DrawRectangle(settingsMenuX, settingsMenuY, settingsMenuWidth, inputFieldHeight, rl.Gray)
	rl.DrawText("Width:", settingsMenuX+10, settingsMenuY+10, 20, rl.Black)
	rl.DrawRectangle(settingsMenuX+60, settingsMenuY+10, inputFieldWidth, inputFieldHeight, rl.White)
	rl.DrawText(fmt.Sprintf("%d", mapWidth), settingsMenuX+65, settingsMenuY+15, 20, rl.Black)

	rl.DrawRectangle(settingsMenuX, settingsMenuY+40, settingsMenuWidth, inputFieldHeight, rl.Gray)
	rl.DrawText("Height:", settingsMenuX+10, settingsMenuY+50, 20, rl.Black)
	rl.DrawRectangle(settingsMenuX+60, settingsMenuY+50, inputFieldWidth, inputFieldHeight, rl.White)
	rl.DrawText(fmt.Sprintf("%d", mapHeight), settingsMenuX+65, settingsMenuY+55, 20, rl.Black)

	rl.DrawRectangle(settingsMenuX+10, settingsMenuY+80, 100, inputFieldHeight, rl.Gray)
	rl.DrawText("Set", settingsMenuX+30, settingsMenuY+85, 20, rl.Black)
}
