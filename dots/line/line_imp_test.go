// Scry Info.  All rights reserved.
// license that can be found in the license file.

package line

import (
	"github.com/scryinfo/dot/dot"
	"reflect"
	"testing"
)

func TestLineImp_Start(t *testing.T) {
	l, _ := BuildAndStart(func(l dot.Line) error {
		_ = l.PreAdd(&dot.TypeLives{
			Meta: dot.Metadata{
				TypeID: "one1",
				NewDoter: func(conf []byte) (dot dot.Dot, err error) {
					t := 1
					return &t, nil
				},
			},
			Lives: []dot.Live{
				{
					LiveID: "one1-1",
				},
				{
					LiveID: "one1-2",
				},
			},
		})

		_ = l.PreAdd(&dot.TypeLives{
			Meta: dot.Metadata{
				TypeID: "one3",
				NewDoter: func(conf []byte) (dot dot.Dot, err error) {
					t := 1
					return &t, nil
				},
			},
			Lives: []dot.Live{
				{
					LiveID: "one3-1",
				},
			},
		})

		_ = l.PreAdd(&dot.TypeLives{
			Meta: dot.Metadata{
				TypeID: "two",
				NewDoter: func(conf []byte) (dot dot.Dot, err error) {
					t := 1
					return &t, nil
				},
			},
			Lives: []dot.Live{
				{
					LiveID:    "two-1",
					RelyLives: map[string]dot.LiveID{"_": "one1-1"},
				},
				{
					LiveID:    "two-2",
					RelyLives: map[string]dot.LiveID{"_": "one1-2"},
				},
			},
		})

		_ = l.PreAdd(&dot.TypeLives{
			Meta: dot.Metadata{
				TypeID: "three",
				NewDoter: func(conf []byte) (dot dot.Dot, err error) {
					t := 1
					return &t, nil
				},
			},
			Lives: []dot.Live{
				{
					LiveID:    "three-1",
					RelyLives: map[string]dot.LiveID{"_": "one1-1"},
				},
				{
					LiveID:    "three-2",
					RelyLives: map[string]dot.LiveID{"_": "two-2"},
				},
			},
		})

		{ //circle
			_ = l.PreAdd(&dot.TypeLives{
				Meta: dot.Metadata{
					TypeID: "circle",
					NewDoter: func(conf []byte) (dot dot.Dot, err error) {
						t := 1
						return &t, nil
					},
				},
				Lives: []dot.Live{
					{
						LiveID:    "circle-1",
						RelyLives: map[string]dot.LiveID{"_": "circle2-1"},
					},
					{
						LiveID: "circle-2",
					},
				},
			})

			_ = l.PreAdd(&dot.TypeLives{
				Meta: dot.Metadata{
					TypeID: "circle2",
					NewDoter: func(conf []byte) (dot dot.Dot, err error) {
						t := 1
						return &t, nil
					},
				},
				Lives: []dot.Live{
					{
						LiveID:    "circle2-1",
						RelyLives: map[string]dot.LiveID{"_": "circle3-1"},
					},
				},
			})

			_ = l.PreAdd(&dot.TypeLives{
				Meta: dot.Metadata{
					TypeID: "circle3",
					NewDoter: func(conf []byte) (dot dot.Dot, err error) {
						t := 1
						return &t, nil
					},
					RelyTypeIDs: []dot.TypeID{"circle"},
				},
				Lives: []dot.Live{
					{
						LiveID:    "circle3-1",
						RelyLives: map[string]dot.LiveID{"_": "circle-1"},
					},
				},
			})

			_ = l.PreAdd(&dot.TypeLives{ //duplication
				Meta: dot.Metadata{
					TypeID: "circle3",
					NewDoter: func(conf []byte) (dot dot.Dot, err error) {
						t := 1
						return &t, nil
					},
					RelyTypeIDs: []dot.TypeID{"circle"},
				},
				Lives: []dot.Live{
					{
						LiveID:    "circle3-1",
						RelyLives: map[string]dot.LiveID{"_": "circle-1"},
					},
				},
			})
		}

		return nil
	})

	lImp, _ := l.(*lineImp)

	order, err := lImp.relyOrder()
	if len(order) != 11 {
		t.Error("len(order) != 11")
	}

	if err == nil {
		t.Error("error == nil, circle dependency")
	}

	_ = l.ToLifer().Stop(true)
	_ = l.ToLifer().Destroy(true)
}

func TestLineImp_InterfaceOrPoint(t *testing.T) {
	type Dot1 struct {
	}

	type Dot2 struct {
		Dot1 Dot1 `dot:""`
	}

	l, err := BuildAndStart(func(l dot.Line) error {
		_ = l.PreAdd(&dot.TypeLives{
			Meta: dot.Metadata{
				TypeID: "dot1",
				NewDoter: func(conf []byte) (dot dot.Dot, err error) {
					return &Dot1{}, nil
				},
			},
			Lives: []dot.Live{
				{
					LiveID: "dot1",
				},
			},
		}, &dot.TypeLives{
			Meta: dot.Metadata{
				TypeID:  "dot2",
				RefType: reflect.TypeOf((*Dot2)(nil)).Elem(),
			},
		})

		return nil
	})

	if err == nil {
		t.Error("the field is not pointer")
	}

	_ = l.ToLifer().Stop(true)
	_ = l.ToLifer().Destroy(true)
}
