package line

import (
	"github.com/scryinfo/dot/dot"
	"sync"
)

type dotEventerImp struct {
	mutex sync.Mutex
	typeEvents map[dot.TypeId]*dot.TypeEvents
	liveEvents map[dot.LiveId]*dot.LiveEvents
}

func (c *dotEventerImp) Init() {
	c.mutex.Lock()
	c.typeEvents = make(map[dot.TypeId]*dot.TypeEvents,0)
	c.liveEvents = make(map[dot.LiveId]*dot.LiveEvents,0)
	c.mutex.Unlock()
}

//
func (c *dotEventerImp) SetLiveEvents(lid dot.LiveId, liveEvents *dot.LiveEvents) {
	c.mutex.Lock()
	c.liveEvents[lid] = liveEvents
	c.mutex.Unlock()
}
func (c *dotEventerImp) LiveEvents(lid dot.LiveId) *dot.LiveEvents {
	c.mutex.Lock()
	f, _ := c.liveEvents[lid]
	c.mutex.Unlock()
	return f
}

func (c *dotEventerImp) SetTypeEvents(lid dot.TypeId, typeEvents *dot.TypeEvents) {
	c.mutex.Lock()
	c.typeEvents[lid] = typeEvents
	c.mutex.Unlock()
}

func (c *dotEventerImp) TypeEvents(lid dot.TypeId) *dot.TypeEvents {
	c.mutex.Lock()
	f, _ := c.typeEvents[lid]
	c.mutex.Unlock()
	return f
}