package line

import (
	"encoding/json"

	"github.com/scryinfo/dot"
)

type DotConfig struct {
	TypeId    dot.TypeId
	LiveId    dot.LiveId
	RelyLives []dot.LiveId
	Json      *json.RawMessage
}

type Config struct {
	Dots []DotConfig
}
