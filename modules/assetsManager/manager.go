package assetsmanager

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

/*
	TODO
	funkcja
	.Resize() do nadpisywania wielkości obrazka z textury
	.Color() przyjmowało by coś na styl rl.White czy coś

	ustawianie pozycji na ekranie przy tworzenie defaultowego sprita w init
	dla statycznych spritow

	dodać funkcję "debug"
	która po uruchomieniu na całym ekranie stworzy scrolowalną listę elementów/assetów
	z podpisami np:
	"Heart_Asset_Obj" from "heart_SpriteSheet"
	animated: false
	frames: 0
	frame: 0
	<heart img>

	i jak coś jest animowane to po kolei wszystkie możliwe animacje


	TODO:
	podzielić assesMaanagera
	na 2 części:
	funkcje do zarządzania asetami
	i
	funkcje do rysowania objektów na podstawie danych z am
*/

type SpritesObj struct {
	Name    string       // nazwa
	SrcFile rl.Texture2D // pointer do wczytanej textury / wczytana textura
}

type AssetsObj struct {
	Name string // nazwa asseta
	// albo pointer do SpriteObj zależnie co będzie szybsze
	SrcFileName string            // nazwa kierująca do odpowiedniego SpritesObj.Name
	SrcFileData *SpritesObj       // zamiast używać SrcFileName można przekazać pointer do textury
	IsAnimates  bool              // czy posiada animację
	IdleFrame   AssetsCoordinates // idle Frame - defaultowa pierwsza klatka
	Animations  []Animation       // lista animacji (może zawierać kilka rużnych animacji)
	Rotation    Rotation

	AnimationFrame int // counter - dla obecnej klatki animacji

	// index - która z animacja z zmiennej Animations jest odtwarzana
	AnimationIndex int
}

// add optional animationFrameCooldown w AssetsObj
func (asset *AssetsObj) NextAnimationFrame() {
	// jeżeli jesteśmy na ostatniej klatce
	if len(asset.Animations[asset.AnimationIndex].Frames) >= asset.AnimationFrame {
		asset.AnimationFrame = 0
	} else {
		asset.AnimationFrame++
	}
}

// rysuj texturę na ekranie pobierając ją z asset.SrcFileData
func (asset *AssetsObj) DrawTextureFromData_Idle(onScreenPosition rl.Rectangle) {

	onSpritePosition := rl.Rectangle{
		X:      asset.IdleFrame.X,
		Y:      asset.IdleFrame.Y,
		Width:  asset.IdleFrame.Width,
		Height: asset.IdleFrame.Height,
	}

	rl.DrawTexturePro(
		asset.SrcFileData.SrcFile, // Tekstura
		onSpritePosition,          // Obszar tekstury
		onScreenPosition,          // Pozycja i rozmiar na ekranie
		rl.Vector2{X: asset.Rotation.RotationOrigin.X, Y: asset.Rotation.RotationOrigin.Y}, // Punkt obrotu
		asset.Rotation.RotationValue, // Brak obrotu
		rl.White)
}

// rysuj texturę na ekranie pobierając ją z asset.SrcFileData
func (asset *AssetsObj) DrawTextureFromData_Animation(onScreenPosition rl.Rectangle) {

	onSpritePosition := rl.Rectangle{
		X:      asset.Animations[asset.AnimationIndex].Frames[asset.AnimationFrame].X,
		Y:      asset.Animations[asset.AnimationIndex].Frames[asset.AnimationFrame].Y,
		Width:  asset.Animations[asset.AnimationIndex].Frames[asset.AnimationFrame].Width,
		Height: asset.Animations[asset.AnimationIndex].Frames[asset.AnimationFrame].Height,
	}

	rl.DrawTexturePro(
		asset.SrcFileData.SrcFile, // Tekstura
		onSpritePosition,          // Obszar tekstury
		onScreenPosition,          // Pozycja i rozmiar na ekranie
		rl.Vector2{X: asset.Rotation.RotationOrigin.X, Y: asset.Rotation.RotationOrigin.Y}, // Punkt obrotu
		asset.Rotation.RotationValue, // Brak obrotu
		rl.White)

	// increse frame count
	asset.NextAnimationFrame()
}

// Draw texture from fromData
// todo
// zrobić to tak aby można było podawać pozycję wewnątrz assetu postaci gracza
// tak jak jest liczone PaperBook left top corner tak samo liczyć sprita gracza
// ustalić opcjonalne predefiniowane punkty na spricie gracza (zależne od rotacji lewo, prawo itp)
// których nie trzeba obliczać poza managerem
func (asset *AssetsObj) DrawTFD_ItemInHand(onScreenPOsition rl.Rectangle, rotation Rotation) {

	onSpritePosition = rl.Rectangle{
		X:      asset.Animations[asset.AnimationIndex].Frames[asset.AnimationFrame].X,
		Y:      asset.Animations[asset.AnimationIndex].Frames[asset.AnimationFrame].Y,
		Width:  asset.Animations[asset.AnimationIndex].Frames[asset.AnimationFrame].Width,
		Height: asset.Animations[asset.AnimationIndex].Frames[asset.AnimationFrame].Height,
	}

	rl.DrawTexturePro(
		asset.SrcFileData.SrcFile,
		onSpritePosition,
		onSpritePosition,
		rl.Vector2{X: rotation.RotationOrigin.X, rotation.RotationOrigin.Y},
		rotation.RotationValue,
		rl.White,
	)
}

// rysuj texturę na ekranie pobierając ją na podstawie asset.SrcFileName
func (asset *AssetsObj) DrawTextureFromName() {

}

// zamiast []float32
type AssetsCoordinates struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

type Animation struct {
	Name   string // nazwa animacji
	Frames []AssetsCoordinates
}

type Rotation struct {
	RotationOrigin Point   // opnkt rotacji
	RotationValue  float32 // rotacja
}

type Point struct {
	X float32
	Y float32
}

// AssetsManager - struktura do zarządzania zasobami
type AssetsManager struct {
	textures map[string]*SpritesObj // mapowanie nazw tekstur na SpritesObj
	assets   map[string]*AssetsObj  // mapowanie nazw na AssetsObj
}

func NewAssetManager() *AssetsManager {
	return &AssetsManager{
		textures: make(map[string]*SpritesObj),
		assets:   make(map[string]*AssetsObj),
	}
}

func (am *AssetsManager) LoadTexture(name string, texture *rl.Texture2D) {
	am.textures[name] = &SpritesObj{Name: name, SrcFile: *texture}
}

func (am *AssetsManager) CreateAsset(name string, textureName string, isAnimates bool, idleFrame AssetsCoordinates, animations []Animation) error {
	texture, ok := am.textures[textureName]
	if !ok {
		return fmt.Errorf("texture %s not found", textureName)
	}

	// Tworzymy AssetsObj
	asset := &AssetsObj{
		Name:           name,
		SrcFileData:    texture,
		IsAnimates:     isAnimates,
		IdleFrame:      idleFrame,
		Animations:     animations,
		AnimationFrame: 0,
		AnimationIndex: 0,
	}

	// Dodajemy asset do managera
	am.assets[name] = asset
	return nil
}

func (am *AssetsManager) GetAssetObj(name string) (*AssetsObj, error) {
	asset, ok := am.assets[name]
	if !ok {
		return nil, fmt.Errorf("asset %s not found", name)
	}
	return asset, nil
}

/*
	zastosowanie:
	game.go:

	assetManager := NewAssetManager()

	// load texture
	assetManager.LoadTexture(texture_name, &rl.Texture2D)
	assetManager.LoadTexture(name, &)

	// create Asset with texture
	assetManager.CreateAsset(
	AssetName,
	texture_name,
	isAnimates,
	idleFrame,
	animations,
	)

	draw.drawFrame(&assetManager)

	package draw:

	drawFrame(*assetManager) {

	assetObj, err := assetManager.GetAssetObj(texture_name)
	AssetsObj.DrawTextureFromData_Animation(rectangle)
	}


*/
