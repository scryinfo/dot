// Scry Info.  All rights reserved.
// license that can be found in the license file.

package line

import (
	"github.com/scryinfo/dot/dot"
	"testing"
)

func TestLineImp_Start(t *testing.T) {
	l, _ := BuildAndStart(func(l dot.Line) error {
		_ = l.PreAdd(&dot.TypeLives{
			Meta: dot.Metadata{
				TypeId: "one1",
				NewDoter: func(args interface{}) (dot dot.Dot, err error) {
					t := 1
					return &t, nil
				},
			},
			Lives: []dot.Live{
				{
					LiveId: "one1-1",
				},
				{
					LiveId: "one1-2",
				},
			},
		})

		_ = l.PreAdd(&dot.TypeLives{
			Meta: dot.Metadata{
				TypeId: "one3",
				NewDoter: func(args interface{}) (dot dot.Dot, err error) {
					t := 1
					return &t, nil
				},
			},
			Lives: []dot.Live{
				{
					LiveId: "one3-1",
				},
			},
		})

		_ = l.PreAdd(&dot.TypeLives{
			Meta: dot.Metadata{
				TypeId: "two",
				NewDoter: func(args interface{}) (dot dot.Dot, err error) {
					t := 1
					return &t, nil
				},
			},
			Lives: []dot.Live{
				{
					LiveId:    "two-1",
					RelyLives: map[string]dot.LiveId{"_": "one1-1"},
				},
				{
					LiveId:    "two-2",
					RelyLives: map[string]dot.LiveId{"_": "one1-2"},
				},
			},
		})

		_ = l.PreAdd(&dot.TypeLives{
			Meta: dot.Metadata{
				TypeId: "three",
				NewDoter: func(args interface{}) (dot dot.Dot, err error) {
					t := 1
					return &t, nil
				},
			},
			Lives: []dot.Live{
				{
					LiveId:    "three-1",
					RelyLives: map[string]dot.LiveId{"_": "one1-1"},
				},
				{
					LiveId:    "three-2",
					RelyLives: map[string]dot.LiveId{"_": "two-2"},
				},
			},
		})

		{ //circle
			_ = l.PreAdd(&dot.TypeLives{
				Meta: dot.Metadata{
					TypeId: "circle",
					NewDoter: func(args interface{}) (dot dot.Dot, err error) {
						t := 1
						return &t, nil
					},
				},
				Lives: []dot.Live{
					{
						LiveId:    "circle-1",
						RelyLives: map[string]dot.LiveId{"_": "circle2-1"},
					},
					{
						LiveId: "circle-2",
					},
				},
			})

			_ = l.PreAdd(&dot.TypeLives{
				Meta: dot.Metadata{
					TypeId: "circle2",
					NewDoter: func(args interface{}) (dot dot.Dot, err error) {
						t := 1
						return &t, nil
					},
				},
				Lives: []dot.Live{
					{
						LiveId:    "circle2-1",
						RelyLives: map[string]dot.LiveId{"_": "circle3-1"},
					},
				},
			})

			_ = l.PreAdd(&dot.TypeLives{
				Meta: dot.Metadata{
					TypeId: "circle3",
					NewDoter: func(args interface{}) (dot dot.Dot, err error) {
						t := 1
						return &t, nil
					},
					RelyTypeIds: []dot.TypeId{dot.TypeId("circle")},
				},
				Lives: []dot.Live{
					{
						LiveId:    "circle3-1",
						RelyLives: map[string]dot.LiveId{"_": "circle-1"},
					},
				},
			})

			_ = l.PreAdd(&dot.TypeLives{ //duplication
				Meta: dot.Metadata{
					TypeId: "circle3",
					NewDoter: func(args interface{}) (dot dot.Dot, err error) {
						t := 1
						return &t, nil
					},
					RelyTypeIds: []dot.TypeId{dot.TypeId("circle")},
				},
				Lives: []dot.Live{
					{
						LiveId:    "circle3-1",
						RelyLives: map[string]dot.LiveId{"_": "circle-1"},
					},
				},
			})
		}

		return nil
	})

	lImp, _ := l.(*lineImp)

	order, err := lImp.RelyOrder()
	if len(order) != 11 {
		t.Error("len(order) != 11")
	}

	if err == nil {
		t.Error("error == nil, circle dependency")
	}

	_ = l.ToLifer().Stop(true)
	_ = l.ToLifer().Destroy(true)
}
