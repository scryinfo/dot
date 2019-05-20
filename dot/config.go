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
	//RelyLives rely lives， 如果不能确定key的值那么 key的值为 livid的值
	RelyLives map[string]LiveId `json:"relyLives"`
	//Json json
	Json *json.RawMessage `json:"json"`
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
func (c *Config) FindConfig(tid TypeId, live LiveId) *LiveConfig {
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
