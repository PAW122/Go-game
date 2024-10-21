package types

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ServerConfig struct {
	Port       int
	Password   string
	MaxPlayers int
}

type PlayerState struct {
	ID          string
	Direction   int
	X           float32
	Y           float32
	PlayerFrame int
	PlayerSrcX  float32
	PlayerSrcY  float32
}

type PlayerObj struct {
	Hp    int32
	MaxHp int32
	Eq    PlayerEq
}

type PlayerEq struct {
	MaxSlots  int
	FreeSlots int
	ItemsList []Item
	// lista int z id ItemList np pierwszy item na pasku to np Itemslist[3] (4 intem w eq list)
	SelectBar []int
}

type Item struct {
	Name string
	Icon string
	Data any
}

type Message struct {
	Action      string  `json:"action"`
	Direction   int     `json:"direction,omitempty"`
	X           float32 `json:"x,omitempty"`
	Y           float32 `json:"y,omitempty"`
	ID          string  `json:"id,omitempty"`
	PlayerFrame int     `json:"player_frame,omitempty"`
	ControlSum  uint32  `json:"control_sum,omitempty"`
	PlayerSrcX  float32 `json:"player_src_x,omitempty"`
	PlayerSrcY  float32 `json:"player_src_y,omitempty"`
}

type ButtonList []Button

type Button struct {
	Rect      rl.Rectangle // Prostokąt reprezentujący pozycję i rozmiar przycisku
	Action    func()       // Funkcja do wywołania po kliknięciu
	Color     rl.Color     // Kolor przycisku
	OverColor rl.Color     // Kolor po najechaniu kursorem
	Key       string
}
