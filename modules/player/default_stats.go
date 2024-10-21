package player

import (
	types "game/types"
)

func GetDefaultStats() types.PlayerObj {
	playerData := types.PlayerObj{
		Hp:    5,
		MaxHp: 10,
		Eq: types.PlayerEq{
			MaxSlots:  64,
			FreeSlots: 64,
			ItemsList: []types.Item{},
		},
	}

	return playerData
}
