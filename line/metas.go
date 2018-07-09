package line

import (
	"sync"

	"github.com/scryinfo/dot"
)

type Metas struct {
	metas sync.Map
}

//
func NewMetas() *Metas {
	return &Metas{}
}

func (ms *Metas) Add(m *dot.MetaData) error {
	if m == nil || m.TypeId.String() == "" {
		return dot.Error.ErrNullParameter
	}

	_, ok := ms.metas.Load(m.TypeId)
	if ok {
		return dot.Error.ErrExited
	}

	ms.metas.LoadOrStore(m.TypeId, m)

	return nil
}

func (ms *Metas) Remove(m *dot.MetaData) error {

	ms.metas.Delete(m.TypeId)

	return nil
}

func (ms *Metas) Get(typeId dot.TypeId) (meta *dot.MetaData, err error) {
	meta = nil
	err = nil

	t, ok := ms.metas.Load(typeId)
	if !ok {
		err = dot.Error.ErrExited
	} else {
		meta, ok = t.(*dot.MetaData)
	}
	return
}

func (ms *Metas) NewDot(t dot.TypeId) (dot dot.Dot, err error) {
	dot = nil
	err = nil

	m,err := ms.Get(t)
	if err != nil {
		dot = m.NewDot(nil)
	}
	return
}

