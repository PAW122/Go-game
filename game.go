package main

import (
	"fmt"
	"game/types"
	"time"

	assetsManager "game/modules/AssetsManager"
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
	runing = true

	AssetsManager *assetsManager.AssetsManager = assetsManager.NewAssetManager()

	bkgColor = rl.NewColor(147, 211, 196, 255)

	texture         rl.Texture2D
	grassSprite     rl.Texture2D
	chestSprite     rl.Texture2D
	playerSprite    rl.Texture2D
	hartSprite      rl.Texture2D
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

	Settings      types.Settings
	Defaultvolume float32 = 0.8
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
		AssetsManager,
		playerObj,
		cam,
		eqBookSprite,
		eqOpen,
		&buttonList,
		&Settings,
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

	Settings.AudioVolume = &Defaultvolume
	Settings.ItemListMaxSlots = 64
	rl.InitWindow(screenWidth, screenHeight, "Pawiu Game")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	playerSprite = rl.LoadTexture("assets/Characters/Basic Charakter Spritesheet.png")

	grassSprite = rl.LoadTexture("assets/Tilesets/Grass.png")
	chestSprite = rl.LoadTexture("assets/Objects/Chest.png")
	hartSprite = rl.LoadTexture("assets/Objects/hart.png")
	eqSprite = rl.LoadTexture("assets/2 SpriteSheet/Png/Paper UI/Folding & Cutout/1.png")
	eqBookSprite = rl.LoadTexture("assets/Pixel_Paper_v1.0/2 Spritesheet/1_v2.png")
	eqBookSideIcons = rl.LoadTexture("assets/Pixel_Paper_v1.0/2 Spritesheet/22.png")

	// load textures to assetManager==========================
	// TODO przenieść to do oddzielnego package
	// w modules
	// package createAssetsObj i tam będą te wszsytkie rzeczy ogarniane

	/*
		do zrobienia z next textur:
		zamienić postać na tego spritesheeta

		chyba wymaga to totalnego przerobienia jak gracz/jego rysowanie działa.

		obecnie w update() zą zmieniane klatki i ich pozycje z spriteSheet

		TODO:
		zmapować wszystkie animacje pod ospowiednie nazwy
		dodać do assetManagera wybieranie animacji po nazwie
		zależnie od tego w jaką stronę 1/2/3/4 gracz jest skierowany
		odpalać odpowiednią animację
	*/

	AssetsManager.LoadTexture("basic_player_SpriteSheet", &playerSprite)
	AssetsManager.CreateAsset(
		"basic_player_Asset_obj",
		"basic_player_SpriteSheet",
		true,
		assetsManager.AssetsCoordinates{
			X:      0,
			Y:      0,
			Width:  48,
			Height: 48,
		},
		[]assetsManager.Animation{
			{
				Name: "Idle_down",
				Frames: []assetsManager.AssetsCoordinates{
					{ // klatka 0
						X:      0,
						Y:      0,
						Width:  48,
						Height: 48,
					},
					{ //chyba dobrze odczytałem
						X:      48,
						Y:      0,
						Width:  48,
						Height: 48,
					},
				},
			},
			{
				Name: "Idle_up",
				Frames: []assetsManager.AssetsCoordinates{
					{
						X:      0,
						Y:      48,
						Width:  48,
						Height: 48,
					},
					{
						X:      48,
						Y:      48,
						Width:  48,
						Height: 48,
					},
				},
			},
			{
				Name: "Idle_left",
				Frames: []assetsManager.AssetsCoordinates{
					{
						X:      0,
						Y:      96,
						Width:  48,
						Height: 48,
					},
					{
						X:      48,
						Y:      96,
						Width:  48,
						Height: 48,
					},
				},
			},
			{
				Name: "Idle_right",
				Frames: []assetsManager.AssetsCoordinates{
					{
						X:      0,
						Y:      144,
						Width:  48,
						Height: 48,
					},
					{
						X:      48,
						Y:      144,
						Width:  48,
						Height: 48,
					},
				},
			},
		},
	)

	AssetsManager.LoadTexture("grass_SpriteSheet", &grassSprite)

	AssetsManager.LoadTexture("heart_SpriteSheet", &hartSprite)
	AssetsManager.CreateAsset(
		"Heart_Asset_Obj",
		"heart_SpriteSheet",
		false,
		assetsManager.AssetsCoordinates{
			X:      0,
			Y:      0,
			Width:  7,
			Height: 5,
		},
		nil,
	)

	AssetsManager.LoadTexture("eqSprite_SpriteSheet", &eqBookSprite)
	AssetsManager.CreateAsset( // jeden z 8 assetow progress bara
		"EqBook_progress_bar_empty_1_Asset_Obj",
		"eqSprite_SpriteSheet",
		false,
		assetsManager.AssetsCoordinates{
			X:      448,
			Y:      2981,
			Width:  12,
			Height: 7,
		},
		nil,
	)

	// ==============================================
	tileDest = rl.NewRectangle(0, 0, 16, 16)

	grassTileSrc = rl.NewRectangle(0, 0, floorTileWidth, floorTileHeight)
	chestTileSrc = rl.NewRectangle(0, 0, chestTileWidth, chestTileHeight)

	// grafikaGracza
	playerSrc = rl.NewRectangle(0, 0, frameWidth, frameHeight)
	// pozycja na ekranie
	playerDest = rl.NewRectangle(110, 100, 100, 100)

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("assets/Avery's Farm.mp3")
	musicPasued = false
	rl.PlayMusicStream(music)

	Settings.MusicStream = &music

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
	rl.UnloadTexture(hartSprite)
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
