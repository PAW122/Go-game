package draw

import (
	"fmt"
	types "game/types"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	volume int = 13
	// Sprites
	// X | Y | Width | Height
	Sprite_progress_bar_empty_1 []float32 = []float32{448, 2981, 12, 7}
	Sprite_progress_bar_empty_2 []float32 = []float32{468, 2981, 9, 7}
	Sprite_progress_bar_empty_3 []float32 = []float32{484, 2981, 9, 7}
	Sprite_progress_bar_empty_4 []float32 = []float32{500, 2981, 12, 7}

	Sprite_progress_bar_full_1 []float32 = []float32{448, 2965, 12, 7}
	Sprite_progress_bar_full_2 []float32 = []float32{468, 2965, 9, 7}
	Sprite_progress_bar_full_3 []float32 = []float32{484, 2965, 9, 7}
	Sprite_progress_bar_full_4 []float32 = []float32{500, 2965, 12, 7}

	Sprite_progress_bar_arrow []float32 = []float32{181, 3091, 7, 9}

	spriteList [][]float32 = [][]float32{
		Sprite_progress_bar_empty_1,
		Sprite_progress_bar_empty_2,
		Sprite_progress_bar_empty_4,
		Sprite_progress_bar_full_1,
		Sprite_progress_bar_full_2,
		Sprite_progress_bar_full_4,
		Sprite_progress_bar_arrow,
	}

	// Ui elements
	ButtonsPadding float32 = 2
	// X | Y | Width | Height
	Audio_MasterVolume_1 []float32 = []float32{310, 62, 12, 7}
	Audio_MasterVolume_2 []float32 = []float32{310 + ButtonsPadding + Audio_MasterVolume_1[2], 62, 9, 7}

	trackedItems = make(map[string]bool)
)

func DrawSettings_Audio(eqBookSprite rl.Texture2D, buttonList *types.ButtonList, spriteCords []float32, settings *types.Settings) {

	//audioStreamVolume := 0.5

	//draw buttons
	drawButtons(buttonList, spriteCords, settings, spriteList, eqBookSprite)
}

/*
TODO sprity
w tych samych miejscach będzie można sprity walnąć
*/
func drawButtons(buttonList *types.ButtonList, spriteCords []float32, settings *types.Settings, spriteList [][]float32, eqBookSprite rl.Texture2D) {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	full_bars := volume
	totalButtons := 11 // Zmieniamy liczbę przycisków na 11 (bez ostatniego)

	// Rysowanie pierwszego przycisku (otwierającego)
	drawButton(1, screenWidth, screenHeight, Audio_MasterVolume_1, buttonList, func() { AudioButton(1, settings) }, spriteCords, spriteList, full_bars, eqBookSprite)

	// Rysowanie środkowych przycisków (od 2 do 10 zamiast 11)
	for index := 2; index <= totalButtons; index++ {
		// Oblicz współrzędne X dla każdego przycisku
		ButtonCords := []float32{
			310 + ButtonsPadding + Audio_MasterVolume_1[2]*float32(index-1), // Poprawione przeliczanie X
			62, // Współrzędna Y
			9,  // Szerokość
			7,  // Wysokość
		}

		// Rysowanie każdego przycisku z odpowiednią pozycją
		drawButton(index, screenWidth, screenHeight, ButtonCords, buttonList, func() { AudioButton(index, settings) }, spriteCords, spriteList, full_bars, eqBookSprite)
	}

	// Rysowanie ostatniego przycisku (zamykającego) - teraz będzie dwunasty
	ButtonCords := []float32{
		310 + ButtonsPadding + Audio_MasterVolume_1[2]*float32(totalButtons), // Poprawione przeliczanie X dla ostatniego przycisku
		62, // Współrzędna Y
		12, // Szerokość
		7,  // Wysokość
	}
	drawButton(totalButtons+1, screenWidth, screenHeight, ButtonCords, buttonList, func() { AudioButton(totalButtons+1, settings) }, spriteCords, spriteList, full_bars, eqBookSprite)
}

func drawButton(id int, screenWidth float32, screenHeight float32, localSpriteCords []float32, buttonList *types.ButtonList, exec func(), spriteCords []float32, spriteList [][]float32, full_bars int, eqBookSprite rl.Texture2D) {
	// Oblicz lewy górny róg sprita
	var X_corner = screenWidth/2 - spriteCords[2]
	var Y_corner = screenHeight/2 - spriteCords[3]

	// Definiujemy button z odpowiednią pozycją i rozmiarem
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

	// Dodajemy przycisk do listy, jeśli jeszcze go nie śledzimy
	if !trackedItems[Button.Key] {
		trackedItems[Button.Key] = true
		*buttonList = append(*buttonList, Button)
	}

	// Zmienna przechowująca sprite dla aktualnego elementu
	var DrawSpriteCords []float32
	var lastFullButton = false

	// Logika wyboru sprita
	if id >= 1 && id <= 11 {
		if full_bars >= id {
			if id == full_bars { // Sprawdzenie, czy to ostatni pełny bar
				lastFullButton = true
			}

			if id == 1 {
				DrawSpriteCords = spriteList[3] // Pełny sprite dla 1
			} else if id > 1 && id < 12 {
				DrawSpriteCords = spriteList[4] // Pełny sprite dla pozostałych
			}
		} else {
			if id == 1 {
				DrawSpriteCords = spriteList[0] // Pusty sprite dla 1
			} else if id > 1 && id < 12 {
				DrawSpriteCords = spriteList[1] // Pusty sprite dla pozostałych
			}
		}
	} else {
		if full_bars >= id && len(spriteList) > 5 && len(spriteList[5]) >= 4 {
			DrawSpriteCords = spriteList[5] // Pełny sprite dla ostatnich elementów
		} else if len(spriteList) > 2 && len(spriteList[2]) >= 4 {
			DrawSpriteCords = spriteList[2] // Pusty sprite dla ostatnich elementów
		} else {
			fmt.Println("Error: spriteList[5] or spriteList[2] is empty or has insufficient elements")
		}
	}

	// Rysowanie sprita z odpowiednimi koordynatami
	drawTexture(eqBookSprite, Button.Rect, DrawSpriteCords)

	if lastFullButton {
		newRect := Button.Rect
		newRect.Y = newRect.Y + 20
		newRect.Width = spriteList[6][2] * size_multiplayer
		newRect.Height = spriteList[6][3] * size_multiplayer

		drawTexture(eqBookSprite, newRect, spriteList[6])
		lastFullButton = false
	}
}

func drawTexture(eqBookSprite rl.Texture2D, ButtonRect rl.Rectangle, DrawSpriteCords []float32) {
	// Sprawdzenie, czy DrawSpriteCords ma wystarczającą liczbę elementów
	if len(DrawSpriteCords) < 4 {
		fmt.Println("Error: DrawSpriteCords has less than 4 elements.")
		return
	}

	// Tworzenie InTextureRect z wartości w DrawSpriteCords
	InTextureRect := rl.Rectangle{
		X:      DrawSpriteCords[0],
		Y:      DrawSpriteCords[1],
		Width:  DrawSpriteCords[2],
		Height: DrawSpriteCords[3],
	}

	rl.DrawTexturePro(
		eqBookSprite,           // Tekstura
		InTextureRect,          // Obszar tekstury
		ButtonRect,             // Pozycja i rozmiar na ekranie
		rl.Vector2{X: 0, Y: 0}, // Punkt obrotu
		0,                      // Brak obrotu
		rl.White)
}

// Buttons Click Functions ===============

// id 0 - 13
func AudioButton(id int, settings *types.Settings) {
	maxVolume := 13
	volume = id
	var newVolume float32 = float32((volume*100)/maxVolume) / 100
	settings.AudioVolume = &newVolume

	rl.SetAudioStreamVolume(settings.MusicStream.Stream, newVolume)
}
