// Scry Info.  All rights reserved.
// license that can be found in the license file.

package must_do

import (
	"errors"
	"log"
	"math/rand"
	"testing"
	"time"
)

func One(msg ...interface{}) error {
	//随机一个10
	i := rand.Intn(11)
	log.Print("当前值为：", i)
	if i != 10 {
		return errors.New("err")
	}
	return nil
}

type two struct {
}

func (t *two) GetTen() bool {
	i := rand.Intn(11)
	log.Print("当前值为：", i)
	if i != 10 {
		return false
	}
	return true
}

func TestMustDo_Start(t *testing.T) {

	do := &MustDo{
		Conf: &MustDoConfig{
			LargeInterval: time.Second * 20,
			SmallInterval: time.Second * 5,
			CycleIndex:    3,
		},
	}
	//第一个,函数测试
	do.Thing = func() bool {
		if One() != nil {
			return false
		} else {
			return true
		}
	}
	do.Start()

	//第二个,方法测试
	tt := &two{}
	do.Thing = func() bool {
		return tt.GetTen()
	}
	do.Start()
	//go do.Start()

}
