// Scry Info.  All rights reserved.
// license that can be found in the license file.

package dot

import (
	"encoding/json"
)

//LiveConfig live config
type LiveConfig struct {
	//LiveID
	LiveID LiveID `json:"liveId"`
	//RelyLives rely lives， If cannot confirm key value then key value is livid value
	RelyLives map[string]LiveID `json:"relyLives"`
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
func (c *Config) FindConfig(typeID TypeID, liveID LiveID) *LiveConfig {
	var lcon *LiveConfig = nil

OutFor:
	for _, dotConfig := range c.Dots {
		if len(typeID.String()) > 0 && typeID != dotConfig.MetaData.TypeID {
			continue
		}

		for _, liveConfig := range dotConfig.Lives {
			if liveConfig.LiveID == liveID || (len(liveConfig.LiveID.String()) < 1 && liveID.String() == dotConfig.MetaData.TypeID.String()) {
				lcon = &liveConfig
				break OutFor
			}
		}
	}

	return lcon
}
