// Scry Info.  All rights reserved.
// license that can be found in the license file.

package line

import (
	"github.com/scryinfo/dot/dot"
	"sync"
)

var (
	_ dot.Eventer = (*dotEventerImp)(nil)
)

type dotEventerImp struct {
	mutex      sync.Mutex
	typeEvents map[dot.TypeId][]dot.TypeEvents
	liveEvents map[dot.LiveId][]dot.LiveEvents
}

func (c *dotEventerImp) Init() {
	c.mutex.Lock()
	c.typeEvents = make(map[dot.TypeId][]dot.TypeEvents, 0)
	c.liveEvents = make(map[dot.LiveId][]dot.LiveEvents, 0)
	c.mutex.Unlock()
}

//
func (c *dotEventerImp) ReSetLiveEvents(lid dot.LiveId, liveEvents *dot.LiveEvents) {
	c.mutex.Lock()
	c.liveEvents[lid] = []dot.LiveEvents{*liveEvents}
	c.mutex.Unlock()
}

func (c *dotEventerImp) AddLiveEvents(lid dot.LiveId, liveEvents *dot.LiveEvents) {
	c.mutex.Lock()
	f, ok := c.liveEvents[lid]
	if ok {
		f = append(f, *liveEvents)
	} else {
		f = []dot.LiveEvents{*liveEvents}
	}
	c.liveEvents[lid] = f
	c.mutex.Unlock()
}

func (c *dotEventerImp) LiveEvents(lid dot.LiveId) []dot.LiveEvents {
	c.mutex.Lock()
	f, _ := c.liveEvents[lid]
	c.mutex.Unlock()
	return f
}

func (c *dotEventerImp) ReSetTypeEvents(tid dot.TypeId, typeEvents *dot.TypeEvents) {
	c.mutex.Lock()
	c.typeEvents[tid] = []dot.TypeEvents{*typeEvents}
	c.mutex.Unlock()
}

func (c *dotEventerImp) AddTypeEvents(tid dot.TypeId, typeEvents *dot.TypeEvents) {
	c.mutex.Lock()
	f, ok := c.typeEvents[tid]
	if ok {
		f = append(f, *typeEvents)
	} else {
		f = []dot.TypeEvents{*typeEvents}
	}
	c.typeEvents[tid] = f
	c.mutex.Unlock()
}

func (c *dotEventerImp) TypeEvents(tid dot.TypeId) []dot.TypeEvents {
	c.mutex.Lock()
	f, _ := c.typeEvents[tid]
	c.mutex.Unlock()
	return f
}
