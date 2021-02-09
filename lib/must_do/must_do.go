// Scry Info.  All rights reserved.
// license that can be found in the license file.

package must_do

import (
	"github.com/scryinfo/dot/dot"
	"time"
)

/*
	必须完成某个动作
	按组进行尝试，组与组之间有间隔时间
	每组固定次数与间隔时间
*/
/*
	第一版初步完成
*/
type MustDoConfig struct {
	LargeInterval time.Duration //每组间隔执行时间
	SmallInterval time.Duration //每次尝试的间隔时间
	CycleIndex    int           //每组循环次数
}

type MustDo struct {
	Thing func() bool
	Conf  *MustDoConfig
}

func (c *MustDo) Start() {

	for true {
		for i := 0; i < c.Conf.CycleIndex; i++ {
			if c.Thing() {
				//todo 日志怎么知道具体做了什么事情
				dot.Logger().Infoln("try successful")
				return
			}
			dot.Logger().Warnln("try fail")
			time.Sleep(c.Conf.SmallInterval)
		}
		time.Sleep(c.Conf.LargeInterval)
	}
}
