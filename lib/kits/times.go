package kits

import "time"

type _time struct {
}

var Times = _time{}

func (t _time) SocondTs() int64 {
	return time.Now().Unix()
}
