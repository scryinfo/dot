// Scry Info.  All rights reserved.
// license that can be found in the license file.

package line

import (
	"sync"

	"github.com/scryinfo/dot/dot"
)

var (
	_ dot.Eventer = (*dotEventerImp)(nil)
)

type dotEventerImp struct {
	mutex      sync.Mutex
	typeEvents map[dot.TypeID][]dot.TypeEvents
	liveEvents map[dot.LiveID][]dot.LiveEvents
}

func (c *dotEventerImp) Init() {
	c.mutex.Lock()
	c.typeEvents = make(map[dot.TypeID][]dot.TypeEvents)
	c.liveEvents = make(map[dot.LiveID][]dot.LiveEvents)
	c.mutex.Unlock()
}

//
func (c *dotEventerImp) ReSetLiveEvents(lid dot.LiveID, liveEvents *dot.LiveEvents) {
	c.mutex.Lock()
	c.liveEvents[lid] = []dot.LiveEvents{*liveEvents}
	c.mutex.Unlock()
}

func (c *dotEventerImp) AddLiveEvents(lid dot.LiveID, liveEvents *dot.LiveEvents) {
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

func (c *dotEventerImp) LiveEvents(lid dot.LiveID) []dot.LiveEvents {
	c.mutex.Lock()
	f := c.liveEvents[lid]
	c.mutex.Unlock()
	return f
}

func (c *dotEventerImp) ReSetTypeEvents(tid dot.TypeID, typeEvents *dot.TypeEvents) {
	c.mutex.Lock()
	c.typeEvents[tid] = []dot.TypeEvents{*typeEvents}
	c.mutex.Unlock()
}

func (c *dotEventerImp) AddTypeEvents(tid dot.TypeID, typeEvents *dot.TypeEvents) {
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

func (c *dotEventerImp) TypeEvents(tid dot.TypeID) []dot.TypeEvents {
	c.mutex.Lock()
	f := c.typeEvents[tid]
	c.mutex.Unlock()
	return f
}
