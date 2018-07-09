package line

import (
	"sync"

	"github.com/scryinfo/dot"
	"github.com/scryinfo/dot/dots"
)

type DotActive struct {
	Meta  dot.MetaData
	Relys []dot.RelyInstance
}

func NewDotActive() *DotActive {
	inst := &DotActive{}
	inst.Relys = make([]dot.RelyInstance, 0)
	return inst
}

type Line struct {
	dot.Lifer
	SLog    dots.SLoger
	SConfig dot.SConfiger
	Metas   Metas
	dots    sync.Map
}

func NewLine() (c *Line) {
	return &Line{}
}

func (ms *Line) Add(m *dot.MetaData) error {
	return ms.Metas.Add(m)
}

///////////////

func (c *Line) Create(conf dot.SConfiger) {

}

//Start
func (c *Line) Start() {

}

//Stop
func (c *Line) Stop() {

}

//Destroy 销毁 Dot
func (c *Line) Destroy() {

}

///////////////
