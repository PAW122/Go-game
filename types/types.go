package types

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
