package line

import (
	"github.com/scryinfo/dot/dot"
	"sync"
)

type dotEvent struct {
	mutex sync.Mutex
	//Before the dot create
	beforeCreate map[dot.LiveId]dot.BeforeCreate
	//after the dot create
	afterCreate map[dot.LiveId]dot.AfterCreate
	//Before the dot start
	beforeStart map[dot.LiveId]dot.BeforeStart
	//After the dot start
	afterStart map[dot.LiveId]dot.AfterStart
	//Before the dot stop
	beforeStop map[dot.LiveId]dot.BeforeStop
	//After the dot stop
	afterStop map[dot.LiveId]dot.AfterStop
	//Before the dot destroy
	beforeDestroy map[dot.LiveId]dot.BeforeDestroy
	//After the dot destroy
	afterDestroy map[dot.LiveId]dot.AfterDestroy
}

func (c *dotEvent) Init() {
	//Before the dot create
	c.beforeCreate = make(map[dot.LiveId]dot.BeforeCreate, 0)
	//after the dot create
	c.afterCreate = make(map[dot.LiveId]dot.AfterCreate, 0)
	//Before the dot start
	c.beforeStart = make(map[dot.LiveId]dot.BeforeStart, 0)
	//After the dot start
	c.afterStart = make(map[dot.LiveId]dot.AfterStart, 0)
	//Before the dot stop
	c.beforeStop = make(map[dot.LiveId]dot.BeforeStop, 0)
	//After the dot stop
	c.afterStop = make(map[dot.LiveId]dot.AfterStop, 0)
	//Before the dot destroy
	c.beforeDestroy = make(map[dot.LiveId]dot.BeforeDestroy, 0)
	//After the dot destroy
	c.afterDestroy = make(map[dot.LiveId]dot.AfterDestroy, 0)
}

//
func (c *dotEvent) SetBeforeCreate(lid dot.LiveId, f dot.BeforeCreate) {
	c.mutex.Lock()
	c.beforeCreate[lid] = f
	c.mutex.Unlock()
}
func (c *dotEvent) BeforeCreate(lid dot.LiveId) dot.BeforeCreate {
	c.mutex.Lock()
	f, _ := c.beforeCreate[lid]
	c.mutex.Unlock()
	return f
}

func (c *dotEvent) SetAfterCreate(lid dot.LiveId, f dot.AfterCreate) {
	c.mutex.Lock()
	c.afterCreate[lid] = f
	c.mutex.Unlock()
}

func (c *dotEvent) AfterCreate(lid dot.LiveId) dot.AfterCreate {
	c.mutex.Lock()
	f, _ := c.afterCreate[lid]
	c.mutex.Unlock()
	return f
}

func (c *dotEvent) SetBeforeStart(lid dot.LiveId, f dot.BeforeStart) {
	c.mutex.Lock()
	c.beforeStart[lid] = f
	c.mutex.Unlock()
}
func (c *dotEvent) BeforeStart(lid dot.LiveId) dot.BeforeStart {
	c.mutex.Lock()
	f, _ := c.beforeStart[lid]
	c.mutex.Unlock()
	return f
}

func (c *dotEvent) SetAfterStart(lid dot.LiveId, f dot.AfterStart) {
	c.mutex.Lock()
	c.afterStart[lid] = f
	c.mutex.Unlock()
}
func (c *dotEvent) AfterStart(lid dot.LiveId) dot.AfterStart {
	c.mutex.Lock()
	f, _ := c.afterStart[lid]
	c.mutex.Unlock()
	return f
}

func (c *dotEvent) SetBeforeStop(lid dot.LiveId, f dot.BeforeStop) {
	c.mutex.Lock()
	c.beforeStop[lid] = f
	c.mutex.Unlock()
}
func (c *dotEvent) BeforeStop(lid dot.LiveId) dot.BeforeStop {
	c.mutex.Lock()
	f, _ := c.beforeStop[lid]
	c.mutex.Unlock()
	return f
}

func (c *dotEvent) SetAfterStop(lid dot.LiveId, f dot.AfterStop) {
	c.mutex.Lock()
	c.afterStop[lid] = f
	c.mutex.Unlock()
}

func (c *dotEvent) AfterStop(lid dot.LiveId) dot.AfterStop {
	c.mutex.Lock()
	f, _ := c.afterStop[lid]
	c.mutex.Unlock()
	return f
}

func (c *dotEvent) SetBeforeDestroy(lid dot.LiveId, f dot.BeforeDestroy) {
	c.mutex.Lock()
	c.beforeDestroy[lid] = f
	c.mutex.Unlock()
}
func (c *dotEvent) BeforeDestroy(lid dot.LiveId) dot.BeforeDestroy {
	c.mutex.Lock()
	f, _ := c.beforeDestroy[lid]
	c.mutex.Unlock()
	return f
}

func (c *dotEvent) SetAfterDestroy(lid dot.LiveId, f dot.AfterDestroy) {
	c.mutex.Lock()
	c.afterDestroy[lid] = f
	c.mutex.Unlock()
}
func (c *dotEvent) AfterDestroy(lid dot.LiveId) dot.AfterDestroy {
	c.mutex.Lock()
	f, _ := c.afterDestroy[lid]
	c.mutex.Unlock()
	return f
}

//clean
func (c *dotEvent) CleanBeforeCreate() {
	c.mutex.Lock()
	c.beforeCreate = nil
	c.mutex.Unlock()
}
func (c *dotEvent) CleanAfterCreate() {
	c.mutex.Lock()
	c.afterCreate = nil
	c.mutex.Unlock()
}
func (c *dotEvent) CleanBeforeStart() {
	c.mutex.Lock()
	c.beforeStart = nil
	c.mutex.Unlock()
}
func (c *dotEvent) CleanAfterStart() {
	c.mutex.Lock()
	c.afterStart = nil
	c.mutex.Unlock()
}
func (c *dotEvent) CleanBeforeStop() {
	c.mutex.Lock()
	c.beforeStop = nil
	c.mutex.Unlock()
}
func (c *dotEvent) CleanAfterStop() {
	c.mutex.Lock()
	c.afterStop = nil
	c.mutex.Unlock()
}
func (c *dotEvent) CleanBeforeDestroy() {
	c.mutex.Lock()
	c.beforeDestroy = nil
	c.mutex.Unlock()
}
func (c *dotEvent) CleanAfterDestroy() {
	c.mutex.Lock()
	c.afterDestroy = nil
	c.mutex.Unlock()
}
