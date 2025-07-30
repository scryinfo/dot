// Scry Info.  All rights reserved.
// license that can be found in the license file.

package dot

// LiveConfig live config
type LiveConfig struct {
	//LiveID
	LiveID LiveID `json:"liveId" toml:"liveId" yaml:"liveId"`
	//RelyLives rely lives， If cannot confirm key value then key value is livid value
	RelyLives map[string]LiveID `json:"relyLives" toml:"relyLives" yaml:"relyLives"`
	//dot config
	//todo tag
	Config interface{} `json:"json" toml:"json" yaml:"json" mapstructure:"json"`
}

// MetaLivesConfig dot config
type MetaLivesConfig struct {
	MetaData Metadata     `json:"metaData" toml:"metaData" yaml:"metaData"`
	Lives    []LiveConfig `json:"lives" toml:"lives" yaml:"lives"`
}

// Config config
type Config struct {
	Log  LogConfig         `json:"log" toml:"log" yaml:"log"`
	Dots []MetaLivesConfig `json:"dots" toml:"dots" yaml:"dots"`
}

// FindConfig find config
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
