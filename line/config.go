package line

import (
	"github.com/scryinfo/dot"
	"encoding/json"
)

type DotConfig struct {
	TypeId dot.TypeId
	InstanceId dot.InstanceId
	RelyInsts []dot.InstanceId
	Json *json.RawMessage
}

type LineConfig struct {

	DotsConfigPath string
	Dots []DotConfig

}