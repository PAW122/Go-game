package main

import (
	"fmt"
	"game/types"
	"time"

	draw "game/modules/draw"
	playerHandler "game/modules/player"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth     = 1620
	screenHeight    = 860
	animationSpeed  = 8  // co ile klatek ma się zmienić animacja
	animationFrames = 4  // ile jest klatek animacji
	frameWidth      = 48 // wielkość klatki
	frameHeight     = 48 // wielkość klatki

	floorTileWidth  = 16 // wielkość klatki (img) (grass)
	floorTileHeight = 16 // wielkość klatki (img) (grass)

	chestTileWidth  = 48 // wielkość klatki (img) (chest)
	chestTileHeight = 48 // wielkość klatki (img) (chest)

	eqFramesCooldown = 15 // cooldown użycia przycisku E
)

var (
	runing   = true
	bkgColor = rl.NewColor(147, 211, 196, 255)

	texture         rl.Texture2D
	grassSprite     rl.Texture2D
	chestSprite     rl.Texture2D
	playerSprite    rl.Texture2D
	hareSprite      rl.Texture2D
	eqSprite        rl.Texture2D
	eqBookSprite    rl.Texture2D
	eqBookSideIcons rl.Texture2D

	// grafika gracza
	playerSrc rl.Rectangle
	// pozycja gracza
	playerDest   rl.Rectangle
	playerMoving bool
	// 0 Down | 1 Up | 2 Left | 3 Right
	playerDir                                     int
	playerUp, playerDown, playerRight, playerLeft bool

	// player animation frame
	playerFrame int

	// game frame count
	frameCount int
	// game frames count for idle animation
	frameCountIdle int

	tileDest     rl.Rectangle
	tileSrc      rl.Rectangle
	grassTileSrc rl.Rectangle
	chestTileSrc rl.Rectangle

	tileMap        []string
	srcMap         []string
	mapW, mapH     int
	layers         [][]string
	collisionLayer []string

	playerSpeed float32 = 3

	musicPasued bool
	music       rl.Music

	cam         rl.Camera2D
	camZoom     float32 = 1.5
	camRotation float32 = 0.0

	// multiplayer
	playersMap   = make(map[string]types.PlayerState)
	playerInputs = make(map[string]types.PlayerState)
	ticker       *time.Ticker
	stopTicker   chan struct{}

	// buttons
	eqOpen bool = false

	// stats
	playerObj types.PlayerObj

	buttonList types.ButtonList
)

func isCollision(x, y int) bool {
	// Sprawdź, czy na danym polu występuje kolizja (zakładamy, że '#' oznacza kolizję)
	index := y*mapW + x
	return collisionLayer[index] == "#"
}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(bkgColor)
	rl.BeginMode2D(cam)

	drawScene()
	draw.DrawUI(
		playerObj,
		cam,
		hareSprite,
		eqBookSprite,
		eqOpen,
		&buttonList,
	)

	rl.EndMode2D()
	rl.EndDrawing()
}

func HandleButtons(buttons *types.ButtonList) {
	mousePosition := rl.GetMousePosition()

	// TODO click cooldown 15klatek
	// efekt najechania na przycisk ma być spritem z assetow

	for _, button := range *buttons {
		isMouseOverButton := rl.CheckCollisionPointRec(mousePosition, button.Rect)

		// Zmiana koloru, jeśli kursor nad przyciskiem
		if isMouseOverButton {
			rl.DrawRectangleRec(button.Rect, button.OverColor)

			// Jeśli lewy przycisk myszy został wciśnięty
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				button.Action() // Wywołanie przypisanej akcji
			}
		} else {
			// Normalny kolor przycisku
			rl.DrawRectangleRec(button.Rect, button.Color)
		}
	}
}

func game_init() {
	rl.InitWindow(screenWidth, screenHeight, "Pawiu Game")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	grassSprite = rl.LoadTexture("assets/Tilesets/Grass.png")
	chestSprite = rl.LoadTexture("assets/Objects/Chest.png")
	hareSprite = rl.LoadTexture("assets/Objects/hart.png")
	eqSprite = rl.LoadTexture("assets/2 SpriteSheet/Png/Paper UI/Folding & Cutout/1.png")
	eqBookSprite = rl.LoadTexture("assets/Pixel_Paper_v1.0/2 Spritesheet/1.png")
	eqBookSideIcons = rl.LoadTexture("assets/Pixel_Paper_v1.0/2 Spritesheet/22.png")

	tileDest = rl.NewRectangle(0, 0, 16, 16)

	grassTileSrc = rl.NewRectangle(0, 0, floorTileWidth, floorTileHeight)
	chestTileSrc = rl.NewRectangle(0, 0, chestTileWidth, chestTileHeight)

	playerSprite = rl.LoadTexture("assets/Characters/Basic Charakter Spritesheet.png")

	playerSrc = rl.NewRectangle(0, 0, frameWidth, frameHeight)
	playerDest = rl.NewRectangle(110, 100, 100, 100)

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("assets/Avery's Farm.mp3")
	musicPasued = false
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(
		rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		rl.NewVector2(
			float32(playerDest.X-(playerDest.Width/2)),
			float32(playerDest.Y-(playerDest.Height/2)),
		),
		camRotation, camZoom,
	)

	// init player stats
	// todo sync with server init stats
	playerObj = playerHandler.GetDefaultStats()

	loadMap("one.map")

}

func quit() {
	rl.UnloadTexture(grassSprite)
	rl.UnloadTexture(chestSprite)
	rl.UnloadTexture(eqSprite)
	rl.UnloadTexture(eqBookSprite)
	rl.UnloadTexture(eqBookSideIcons)
	rl.UnloadTexture(hareSprite)
	rl.UnloadTexture(playerSprite)
	rl.UnloadMusicStream(music)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}

/*
TODO
informacja o hp gracza
informacja czy gracz ma otworzone eq

	> jeżeli tak dodawać jakąś ikonkę np 3 kropek nad graczem
*/
func sendPosition(force bool) {
	if force {
		fmt.Print("force data update")
	}
	// Tylko jeśli pozycja lub kierunek się zmieniły
	if playerDest.X != lastX || playerDest.Y != lastY || playerDir != previousDir || force {
		fmt.Println(playerSrc)

		xInt := int32(playerDest.X * 1000) // Scaling factor to preserve precision
		yInt := int32(playerDest.Y * 1000)
		sum := uint32(xInt*31 + yInt)
		ctrl_sum := sum & 0xFFFFFFFF

		go SendToServer(types.Message{
			Action:      "loop",
			ID:          clientID,
			X:           playerDest.X,
			Y:           playerDest.Y,
			Direction:   playerDir,
			PlayerFrame: playerFrame,
			ControlSum:  ctrl_sum,
			PlayerSrcX:  playerSrc.X,
			PlayerSrcY:  playerSrc.Y,
		})
		lastX = playerDest.X
		lastY = playerDest.Y
		previousDir = playerDir
	}
}

var lastX, lastY float32
var previousDir int

func StartGame() {

	ticker = time.NewTicker(time.Second)
	stopTicker = make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				sendPosition(true)
			case <-stopTicker:
				ticker.Stop()
				return
			}
		}
	}()

	game_init()

	for runing {
		input()
		HandleButtons(&buttonList)
		update()
		sendPosition(false)
		render()

	}

	defer func() { stopTicker <- struct{}{} }()
	quit()
}