// Scry Info.  All rights reserved.
// license that can be found in the license file.

package dot

import (
	"encoding/json"
)

//LiveConfig live config
type LiveConfig struct {
	//LiveId
	LiveId LiveId `json:"liveId"`
	//RelyLives rely lives， If cannot confirm key value then key value is livid value
	RelyLives map[string]LiveId `json:"relyLives"`
	//Json json
	JSON *json.RawMessage `json:"json"`
}

//DotConfig dot config
type DotConfig struct {
	MetaData Metadata     `json:"metaData"`
	Lives    []LiveConfig `json:"lives"`
}

//Config config
type Config struct {
	Log  LogConfig   `json:"log"`
	Dots []DotConfig `json:"dots"`
}

//FindConfig find config
func (c *Config) FindConfig(typeId TypeId, live LiveId) *LiveConfig {
	var lcon *LiveConfig = nil

OutFor:
	for _, dotConfig := range c.Dots {
		if len(typeId.String()) > 0 && typeId != dotConfig.MetaData.TypeId {
			continue
		}

		for _, liveConfig := range dotConfig.Lives {
			if liveConfig.LiveId == live || (len(liveConfig.LiveId.String()) < 1 && live.String() == dotConfig.MetaData.TypeId.String()) {
				lcon = &liveConfig
				break OutFor
			}
		}
	}

	return lcon
}
