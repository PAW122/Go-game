package draw

import (
	types "game/types"

	rl "github.com/gen2brain/raylib-go/raylib"
)

/*
	TODO:

	scrolowana lista eq
	skalująca się zależnie od max Ilości slotów
	zrobić listę slotów eq
	zrobić listę all przedmiotów

	- może
	renderować tylko widoczne sloty
	( if slot jest poza bookSprite return )

	TODO:
	zrobić jakąś package z wszystkimi spritamia razem z gotowymi kordami
	+docsa (obrazek z podpisanymi co jest czym)
	tak żeby w przypadku reuse można było sobie "importować"
	albo jak się zmieni aset to tylkow 1 miejsco zmieniać wartości

	+ tosamo ale z opcją pobrania całej listy animacji
	np getSprite(<name> <animation bool> <animation_name bool>))
	sword, false
	[x,x,x,x]
	sword, true, attack
	[[x,x,x,x], [x,x,x,x],[x,x,x,x] [x,x,x,x]]
	animacja ataku mieczem składajaca się z 4 klatem
*/

var (
	ItemsSlots []any

	// ramka na item składa się z 9 elementów
	Eq_Slot_Sprite_1 []float32 = []float32{64, 2896, 15, 15}
	Eq_Slot_Sprite_2 []float32 = []float32{96, 2896, 15, 15}
	Eq_Slot_Sprite_3 []float32 = []float32{128, 2896, 15, 16}

	Eq_Slot_Sprite_4 []float32 = []float32{64, 2928, 15, 15}
	Eq_Slot_Sprite_5 []float32 = []float32{96, 2928, 15, 15} // środek - blank
	Eq_Slot_Sprite_6 []float32 = []float32{128, 2928, 15, 16}

	Eq_Slot_Sprite_7 []float32 = []float32{64, 2960, 15, 15}
	Eq_Slot_Sprite_8 []float32 = []float32{96, 2960, 15, 15}
	Eq_Slot_Sprite_9 []float32 = []float32{128, 2960, 15, 15}

	Scroolbar_Slider_Sprite   []float32 = []float32{246, 2994, 5, 11} // scrolbar - move element
	Scroolbar_Holder_Sprite_1 []float32 = []float32{230, 2992, 5, 16} // scroolbar - background sprite
	Scroolbar_Holder_Sprite_2 []float32 = []float32{230, 3024, 5, 16}
	Scroolbar_Holder_Sprite_3 []float32 = []float32{230, 2992, 5, 11}

	Selected_Slot_Sprite []float32 = []float32{119, 3056, 20, 20}
	Bin_Default_Sprite   []float32 = []float32{624, 3085, 30, 63} // ikonka śmiernika - bez animacji

	selected_Slot_Id int = 0
)

func DrawItemsPage(eqBookSprite rl.Texture2D, buttonList *types.ButtonList, spriteCords []float32, settings *types.Settings) {

}
