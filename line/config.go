package line

import (
	"encoding/json"

	"github.com/scryInfo/dot/dot"
)

//LiveConfig live config
type LiveConfig struct {
	//LiveId
	LiveId dot.LiveId
	//RelyLives rely lives
	RelyLives []dot.LiveId
	//Json json
	Json *json.RawMessage
}

//DotConfig dot config
type DotConfig struct {
	MetaData dot.Metadata
	Lives    []LiveConfig
}

//Config config
type Config struct {
	Dots []DotConfig
}

//FindConfig find config
func (c *Config) FindConfig(tid dot.TypeId, live dot.LiveId) *LiveConfig {
	var lcon *LiveConfig = nil

OUT_FOR:
	for _, it := range c.Dots {
		if len(tid.String()) > 0 && tid != it.MetaData.TypeId {
			continue
		}

		for _, li := range it.Lives {
			if li.LiveId == live || (len(li.LiveId.String()) < 1 && live.String() == it.MetaData.TypeId.String()) {
				lcon = &li
				break OUT_FOR
			}
		}
	}

	return lcon
}
