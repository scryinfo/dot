package line

import (
	"sync"

	"github.com/scryinfo/dot"
	"github.com/scryinfo/dot/dots"
)

var (
	_ dot.Lifer = (*Line)(nil)
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

func (c *Line) Create(conf dot.SConfiger) error {

	//first create config
	c.SConfig = dot.NewSConfig()
	c.SConfig.Create(nil)

	//create log

	//create others

	return nil
}

//Start
func (c *Line) Start() error {

	c.SConfig.Start()

	return nil
}

//Stop
func (c *Line) Stop() error {

	//stop others

	//stop log

	//stop config
	c.SConfig.Stop()

	return nil
}

//Destroy 销毁 Dot
func (c *Line) Destroy() error {

	//Destroy others

	//Destroy log

	//Destroy config
	c.SConfig.Destroy()
	return nil
}

///////////////
