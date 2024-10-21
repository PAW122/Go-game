package draw

import (
	"fmt"
	types "game/types"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	// Pozycja paska życia
	hp_bar_cords = []int{0, 0}
)

/*
draw all ui elements
*/
func DrawUI(
	playerObj types.PlayerObj,
	cam rl.Camera2D,
	heartTexture rl.Texture2D,
	eqSprite rl.Texture2D,
	eqOpen bool,
	buttonList *types.ButtonList,
) {
	// Najpierw rysujemy mapę, a potem UI
	drawHpBar(playerObj, cam, heartTexture)
	drawEq(playerObj, cam, eqSprite, eqOpen, buttonList)
}

/*
w przyciskać dodać E - eq
if eq = true
to rysować eq z tych nowych asetow (ten książkowy styl)

tak samo można będzie dodać dla ustawień (esc)

+ pasek z przedmiotami do scrolowania
+ będzie zintegrowany np z jednym z pasków eq
*/
var (
	// Zmienne pozycji książek
	buttonsWidth float32 = 25
	// X | Y | width | height
	blandBook             []float32 = []float32{96, 144, 517 + buttonsWidth, 320}
	profileAndMapBook     []float32 = []float32{735, 144, 517 + buttonsWidth, 320}
	inventoryAndItemsBook []float32 = []float32{1376, 144, 517 + buttonsWidth, 320}
	statusAndEq           []float32 = []float32{2016, 144, 517 + buttonsWidth, 320}
	questsAndTasksBook    []float32 = []float32{96, 592, 517 + buttonsWidth, 320}
	saveSysBook           []float32 = []float32{736, 592, 517 + buttonsWidth, 320}
	settingsBook          []float32 = []float32{1376, 592, 517 + buttonsWidth, 320}

	drawEqBook []float32 = blandBook

	size_multiplayer float32 = 2.0

	// Zmienne dla przycisków Width | Height
	button_obj          []float32 = []float32{19, 18}
	button_top_offset   float32   = 31 // przerwa pomiędzy sprite.Top a pierwszym pryciskiem
	buttons_padding     float32   = 8  // przerwa pomiędzy przyciskami
	clicked_size_change float32   = 5  // o ile zwiększa się przycisk po kliknięciu
	clickedButtonId     int8
	addSizeList         []float32 = []float32{0, 0, 0, 0, 0, 0}

	// add buttons to ButtonsHandler
	buttonsAdded   = false
	buttonsDeleted = false
	reloadButtons  = false

	// Zmienna trybu debugowania
	debugMode bool = false

	transparentColor rl.Color = rl.NewColor(0, 0, 255, 0)
)

// Funkcja rysująca ekwipunek i dodająca przyciski do listy
func drawEq(playerObj types.PlayerObj, cam rl.Camera2D, eqBookSprite rl.Texture2D, eqOpen bool, buttonList *types.ButtonList) {
	if !eqOpen {
		if !buttonsDeleted {
			// Usunięcie przycisków po zamknięciu
			for i := len(*buttonList) - 1; i >= 0; i-- {
				if (*buttonList)[i].Key == "EQ_profile_button" ||
					(*buttonList)[i].Key == "EQ_items_button" ||
					(*buttonList)[i].Key == "EQ_inventory_button" ||
					(*buttonList)[i].Key == "EQ_quests_button" || //
					(*buttonList)[i].Key == "EQ_saves_button" ||
					(*buttonList)[i].Key == "EQ_settings_button" { // Sprawdzenie kluczy
					*buttonList = append((*buttonList)[:i], (*buttonList)[i+1:]...) // Usuwanie przycisku
				}
			}
			buttonsDeleted = true
			buttonsAdded = false // Resetowanie flagi
		}
		return
	}

	buttonsDeleted = false
	spriteCords := drawEqBook

	// Zatrzymanie trybu kamery
	rl.EndMode2D()

	// Ustawienie prostokąta dla UI (stała pozycja na ekranie)
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	// Tworzenie przycisku profilu
	buttonNumber := float32(0)
	var profileButton types.Button = types.Button{
		Rect: rl.Rectangle{
			X:      screenWidth/2 + spriteCords[2]*size_multiplayer/2 - buttonsWidth*size_multiplayer,
			Y:      screenHeight/2 - spriteCords[3]*size_multiplayer/2 + button_top_offset*size_multiplayer,
			Width:  button_obj[0]*size_multiplayer + (addSizeList[int(buttonNumber)] * size_multiplayer),
			Height: button_obj[1] * size_multiplayer,
		},
		Action:    OpenProfileBook_Button,
		Color:     transparentColor,
		OverColor: transparentColor,
		Key:       "EQ_profile_button",
	}

	buttonNumber = 1
	var itemsButton = types.Button{
		Rect: rl.Rectangle{
			X:      screenWidth/2 + spriteCords[2]*size_multiplayer/2 - buttonsWidth*size_multiplayer,
			Y:      screenHeight/2 - spriteCords[3]*size_multiplayer/2 + button_top_offset*size_multiplayer + (buttons_padding*size_multiplayer+button_obj[0]*size_multiplayer)*buttonNumber,
			Width:  button_obj[0]*size_multiplayer + (addSizeList[int(buttonNumber)] * size_multiplayer),
			Height: button_obj[1] * size_multiplayer,
		},
		Action:    OpenItemsBook_Button,
		Color:     transparentColor,
		OverColor: transparentColor,
		Key:       "EQ_items_button",
	}

	buttonNumber = 2
	var inventoryAndItems_Button types.Button = types.Button{
		Rect: rl.Rectangle{
			X:      screenWidth/2 + spriteCords[2]*size_multiplayer/2 - buttonsWidth*size_multiplayer,
			Y:      screenHeight/2 - spriteCords[3]*size_multiplayer/2 + button_top_offset*size_multiplayer + (buttons_padding*size_multiplayer+button_obj[0]*size_multiplayer)*buttonNumber,
			Width:  button_obj[0]*size_multiplayer + (addSizeList[int(buttonNumber)] * size_multiplayer),
			Height: button_obj[1] * size_multiplayer,
		},
		Action:    OpenInventoryAndItems_Button,
		Color:     transparentColor,
		OverColor: transparentColor,
		Key:       "EQ_inventory_button",
	}

	buttonNumber = 3
	var questsAndTasks_Button types.Button = types.Button{
		Rect: rl.Rectangle{
			X:      screenWidth/2 + spriteCords[2]*size_multiplayer/2 - buttonsWidth*size_multiplayer,
			Y:      screenHeight/2 - spriteCords[3]*size_multiplayer/2 + button_top_offset*size_multiplayer + (buttons_padding*size_multiplayer+button_obj[0]*size_multiplayer)*buttonNumber,
			Width:  button_obj[0]*size_multiplayer + (addSizeList[int(buttonNumber)] * size_multiplayer),
			Height: button_obj[1] * size_multiplayer,
		},
		Action:    OpenQuestsBook_Button,
		Color:     transparentColor,
		OverColor: transparentColor,
		Key:       "EQ_quests_button",
	}

	buttonNumber = 4
	var saveSys_Button types.Button = types.Button{
		Rect: rl.Rectangle{
			X:      screenWidth/2 + spriteCords[2]*size_multiplayer/2 - buttonsWidth*size_multiplayer,
			Y:      screenHeight/2 - spriteCords[3]*size_multiplayer/2 + button_top_offset*size_multiplayer + (buttons_padding*size_multiplayer+button_obj[0]*size_multiplayer)*buttonNumber,
			Width:  button_obj[0]*size_multiplayer + (addSizeList[int(buttonNumber)] * size_multiplayer),
			Height: button_obj[1] * size_multiplayer,
		},
		Action:    OpenSavesBook_Button,
		Color:     transparentColor,
		OverColor: transparentColor,
		Key:       "EQ_saves_button",
	}

	buttonNumber = 5
	var settings_Button types.Button = types.Button{
		Rect: rl.Rectangle{
			X:      screenWidth/2 + spriteCords[2]*size_multiplayer/2 - buttonsWidth*size_multiplayer,
			Y:      screenHeight/2 - spriteCords[3]*size_multiplayer/2 + button_top_offset*size_multiplayer + (buttons_padding*size_multiplayer+button_obj[0]*size_multiplayer)*buttonNumber,
			Width:  button_obj[0]*size_multiplayer + (addSizeList[int(buttonNumber)] * size_multiplayer),
			Height: button_obj[1] * size_multiplayer,
		},
		Action:    OpenSettingsBook_Button,
		Color:     transparentColor,
		OverColor: transparentColor,
		Key:       "EQ_settings_button",
	}

	// Dodanie przycisku do listy
	if !buttonsAdded {
		*buttonList = append(*buttonList, profileButton)
		*buttonList = append(*buttonList, itemsButton)
		*buttonList = append(*buttonList, inventoryAndItems_Button)
		*buttonList = append(*buttonList, questsAndTasks_Button)
		*buttonList = append(*buttonList, saveSys_Button)
		*buttonList = append(*buttonList, settings_Button)

		buttonsAdded = true
	}

	if reloadButtons {
		// delete this scene buttons
		for i := len(*buttonList) - 1; i >= 0; i-- {
			if (*buttonList)[i].Key == "EQ_profile_button" ||
				(*buttonList)[i].Key == "EQ_items_button" ||
				(*buttonList)[i].Key == "EQ_inventory_button" ||
				(*buttonList)[i].Key == "EQ_quests_button" || //
				(*buttonList)[i].Key == "EQ_saves_button" ||
				(*buttonList)[i].Key == "EQ_settings_button" { // Sprawdzenie kluczy
				*buttonList = append((*buttonList)[:i], (*buttonList)[i+1:]...) // Usuwanie przycisku
			}
		}

		//register new buttons
		*buttonList = append(*buttonList, profileButton)
		*buttonList = append(*buttonList, itemsButton)
		*buttonList = append(*buttonList, inventoryAndItems_Button)
		*buttonList = append(*buttonList, questsAndTasks_Button)
		*buttonList = append(*buttonList, saveSys_Button)
		*buttonList = append(*buttonList, settings_Button)

		reloadButtons = false
	}

	// Rysowanie tekstury UI
	rect := rl.Rectangle{
		X:      spriteCords[0],
		Y:      spriteCords[1],
		Width:  spriteCords[2],
		Height: spriteCords[3],
	}

	posRect := rl.Rectangle{
		X:      screenWidth/2 - spriteCords[2]*size_multiplayer/2,
		Y:      screenHeight/2 - spriteCords[3]*size_multiplayer/2,
		Width:  spriteCords[2] * size_multiplayer,
		Height: spriteCords[3] * size_multiplayer,
	}

	rl.DrawTexturePro(
		eqBookSprite,           // Tekstura
		rect,                   // Obszar tekstury
		posRect,                // Pozycja i rozmiar na ekranie
		rl.Vector2{X: 0, Y: 0}, // Punkt obrotu
		0,                      // Brak obrotu
		rl.White)

	// Tryb debugowania: rysowanie prostokątów przycisków
	if debugMode {
		for _, btn := range *buttonList {
			rl.DrawRectangleLinesEx(btn.Rect, 2, rl.Red) // Czerwony obrys przycisku w trybie debugowania
		}
	}

	// Powrót do trybu kamery
	rl.BeginMode2D(cam)
}

// draw Hp bar
// TODO zamienić na element z spritesheets do lewego gornego rogu
func drawHpBar(playerObj types.PlayerObj, cam rl.Camera2D, heartTexture rl.Texture2D) {
	// Pobieramy HP gracza
	hp := playerObj.Hp

	// Pozycja kamery, aby pasek życia rysować w odniesieniu do niej
	cameraX := int32(cam.Target.X)
	cameraY := int32(cam.Target.Y)

	// Pozycja, w której będzie rysowana liczba HP i ikona serca
	x := cameraX + int32(hp_bar_cords[0])
	y := cameraY + int32(hp_bar_cords[1])

	// Rysujemy liczbę HP
	rl.DrawText(fmt.Sprintf("%d", hp), x+20, y, 10, rl.White)

	// Zdefiniuj obszar tekstury, który odpowiada ikonie serca
	heartRect := rl.Rectangle{
		X:      0, // 5 pikseli od lewej strony
		Y:      0, // 30 pikseli od góry
		Width:  7, // Szerokość 7 pikseli
		Height: 5, // Wysokość 5 pikseli
	}

	// Definiujemy, gdzie na ekranie wyświetlić ikonę serca
	heartPosRect := rl.Rectangle{
		X:      float32(x + 30), // Pozycja X na ekranie (obok liczby HP)
		Y:      float32(y),      // Pozycja Y na ekranie
		Width:  14,              // Szerokość ikony na ekranie (skaluje do większych wymiarów)
		Height: 10,              // Wysokość ikony na ekranie (skaluje do większych wymiarów)
	}

	// Rysujemy wyciętą część tekstury (ikonę serca)
	rl.DrawTexturePro(
		heartTexture,           // Tekstura
		heartRect,              // Obszar tekstury (ikonka serca)
		heartPosRect,           // Pozycja i rozmiar na ekranie
		rl.Vector2{X: 0, Y: 0}, // Punkt obrotu (nie ma obrotu, więc to 0,0)
		0,                      // Brak obrotu
		rl.White)               // Kolor (biały, ponieważ nie chcemy zmieniać koloru tekstury)
}

func OpenBlankBook_Button() {
	drawEqBook = blandBook
	// nie liczy się jako kliknięty przycisk
}

// EQ_profile_button
func OpenProfileBook_Button() {
	fmt.Println("Profile Button click ****") //
	drawEqBook = profileAndMapBook
	clickedButtonId = 0
	addSizeList = []float32{0, 0, 0, 0, 0, 0}
	addSizeList[clickedButtonId] = clicked_size_change

	//replace button
	reloadButtons = true
}

// EQ_items_button
func OpenItemsBook_Button() {
	fmt.Println("items_Button click ****")
	drawEqBook = inventoryAndItemsBook
	clickedButtonId = 1
	addSizeList = []float32{0, 0, 0, 0, 0, 0}
	addSizeList[clickedButtonId] = clicked_size_change

	//replace button
	reloadButtons = true
}

func OpenInventoryAndItems_Button() {
	fmt.Println("Eq & Items Button click ****") //
	drawEqBook = statusAndEq
	clickedButtonId = 2
	addSizeList = []float32{0, 0, 0, 0, 0, 0}
	addSizeList[clickedButtonId] = clicked_size_change

	//replace button
	reloadButtons = true
}

func OpenQuestsBook_Button() {
	fmt.Println("Quests Button click ****") //
	drawEqBook = questsAndTasksBook
	clickedButtonId = 3
	addSizeList = []float32{0, 0, 0, 0, 0, 0}
	addSizeList[clickedButtonId] = clicked_size_change

	//replace button
	reloadButtons = true
}

func OpenSavesBook_Button() {
	fmt.Println("saves Button click ****") //
	drawEqBook = saveSysBook
	clickedButtonId = 4
	addSizeList = []float32{0, 0, 0, 0, 0, 0}
	addSizeList[clickedButtonId] = clicked_size_change

	//replace button
	reloadButtons = true
}

func OpenSettingsBook_Button() {
	fmt.Println("settings Button click ****") //
	drawEqBook = settingsBook
	clickedButtonId = 5
	addSizeList = []float32{0, 0, 0, 0, 0, 0}
	addSizeList[clickedButtonId] = clicked_size_change

	//replace button
	reloadButtons = true
}
