package draw

import (
	"fmt"
	types "game/types"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	// Sprites
	// X | Y | Width | Height
	Sprite_progress_bar_empty_1 []float32 = []float32{448, 2901, 15, 6}

	// Ui elements
	ButtonsPadding float32 = 2
	// X | Y | Width | Height
	Audio_MasterVolume_1 []float32 = []float32{310, 62, 12, 7}
	Audio_MasterVolume_2 []float32 = []float32{310 + ButtonsPadding + Audio_MasterVolume_1[2], 62, 9, 7}

	/*
		gor sprita 592
		654
		62

		lewa sprita 1376
		1686
		310
	*/

	trackedItems = make(map[string]bool)
)

func DrawSettings_Audio(eqBookSprite rl.Texture2D, buttonList *types.ButtonList, spriteCords []float32) {

	//audioStreamVolume := 0.5

	//draw buttons
	drawButtons(buttonList, spriteCords)
}

/*
TODO sprity
w tych samych miejscach będzie można sprity walnąć
*/
func drawButtons(buttonList *types.ButtonList, spriteCords []float32) {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	totalButtons := 12 // - 1 initial - 1 kończący

	drawButton(0, screenWidth, screenHeight, Audio_MasterVolume_1, buttonList, func() { AudioButton(0) }, spriteCords)
	for index := range totalButtons {
		index++ // start from index 1 for element index 0
		ButtonCords := []float32{310 + ButtonsPadding + Audio_MasterVolume_1[2]*float32(index), 62, 9, 7}
		if index == 1 {
			ButtonCords[0] = 310 + ButtonsPadding + Audio_MasterVolume_1[2]*float32(index)
		} else {
			ButtonCords[0] = 310 + ButtonsPadding + Audio_MasterVolume_1[2]*float32(index) - float32(index)
			ButtonCords[2] = 10
		}

		drawButton(index, screenWidth, screenHeight, ButtonCords, buttonList, func() { AudioButton(index) }, spriteCords)
	}

	ButtonCords := []float32{310 + ButtonsPadding + Audio_MasterVolume_1[2]*float32(totalButtons), 62, 12, 7}
	drawButton(13, screenWidth, screenHeight, ButtonCords, buttonList, func() { AudioButton(13) }, spriteCords)
}

func drawButton(id int, screenWidth float32, screenHeight float32, localSpriteCords []float32, buttonList *types.ButtonList, exec func(), spriteCords []float32) {
	// oblicz lewy górny róg sprita
	var X_corner = screenWidth/2 - spriteCords[2]
	var Y_corner = screenHeight/2 - spriteCords[3]

	var Button types.Button = types.Button{
		Rect: rl.Rectangle{
			X:      X_corner + localSpriteCords[0]*size_multiplayer,
			Y:      Y_corner + localSpriteCords[1]*size_multiplayer,
			Width:  localSpriteCords[2] * size_multiplayer, // Szerokość przycisku
			Height: localSpriteCords[3] * size_multiplayer, // Wysokość przycisku
		},
		Action:    exec,
		Color:     transparentColor, // Ustawienie przezroczystości
		OverColor: transparentColor, // Ustawienie przezroczystości
		Key:       fmt.Sprintf("UI_SETTINGS_AUDIO_BUTTON_%v", id),
	}

	if !trackedItems[Button.Key] {
		trackedItems[Button.Key] = true
		*buttonList = append(*buttonList, Button)
	}
}

/*
narazie same przyciski
odpowiednie sprity doda się póżnej
*/
func DrawSprites() {

}

// Buttons Click Functions ===============

// id 0 - 13
func AudioButton(id int) {
	fmt.Println("Volume %v", id)
}

/*
	TODO
	w rl głośność jest od 0.1 do 1.0
	więc zrobić licznie ile w % głosność to jest i /100
*/
