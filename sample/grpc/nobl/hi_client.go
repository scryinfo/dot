// Scry Info.  All rights reserved.
// license that can be found in the license file.

package nobl

import (
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/grpc/conns"
	"github.com/scryinfo/dot/sample/grpc/go_out/hidot"
)

const (
	HiClientTypeId = "hiclient"
)

type HiClient struct {
	Conn     *conns.ConnName `dot:""`
	hiclient hidot.HiDotClient
}

func (c *HiClient) HiClient() hidot.HiDotClient {
	return c.hiclient
}

func (c *HiClient) AfterAllInject(l dot.Line) {
	c.hiclient = hidot.NewHiDotClient(c.Conn.Conn())
}

//HiClientTypeLives make all type lives
func HiClientTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: HiClientTypeId, NewDoter: func(conf interface{}) (dot.Dot, error) {
			return &HiClient{}, nil
		}},
		Lives: []dot.Live{
			dot.Live{
				LiveId:    HiClientTypeId,
				RelyLives: map[string]dot.LiveId{"Conn": conns.ConnNameTypeId},
			},
		},
	}

	lives := make([]*dot.TypeLives, 0, 3)
	lives = append(lives, conns.ConnNameTypeLives()...)
	lives = append(lives, tl)

	return lives
}
