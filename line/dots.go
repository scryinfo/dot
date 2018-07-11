package line

import (
	"github.com/scryinfo/dot"
)

//Metas
type Metas struct {
	metas map[dot.TypeId]*dot.MetaData
}

//NewMetas
func NewMetas() *Metas {
	m := &Metas{}
	m.metas = make(map[dot.TypeId]*dot.MetaData, 0)
	return m
}

//Lives
type Lives struct {
	// typeIdMap map[dot.TypeId][]*dot.Live
	liveIdMap map[dot.LiveId]*dot.Live
}

func NewLives() *Lives {
	l := &Lives{}
	l.liveIdMap = make(map[dot.LiveId]*dot.Live)
	return l
}

func (ms *Metas) Add(m *dot.MetaData) error {
	if m == nil || m.TypeId.String() == "" {
		return dot.SError.ErrNullParameter
	}

	_, ok := ms.metas[m.TypeId]
	if ok {
		return dot.NewError(dot.SError.ErrExisted.Code(), dot.SError.ErrExisted.Error()+m.TypeId.String())
	}

	ms.metas[m.TypeId] = m.Clone()

	return nil
}

func (ms *Metas) Remove(m *dot.MetaData) error {

	delete(ms.metas, m.TypeId)
	return nil
}

func (ms *Metas) Get(typeId dot.TypeId) (meta *dot.MetaData, err error) {
	meta = nil
	err = nil

	meta, ok := ms.metas[typeId]
	if !ok {
		err = dot.NewError(dot.SError.ErrNotExisted.Code(), dot.SError.ErrNotExisted.Error()+typeId.String())
	}
	return
}

func (ms *Metas) NewDot(t dot.TypeId) (dot dot.Dot, err error) {
	dot = nil
	err = nil

	m, err := ms.Get(t)
	if err != nil {
		dot = m.NewDot(nil)
	}
	return
}

func (ms *Lives) Add(m *dot.Live) error {
	if m == nil || m.TypeId.String() == "" {
		return dot.SError.ErrNullParameter
	}

	_, ok := ms.liveIdMap[m.LiveId]
	if ok {
		return dot.NewError(dot.SError.ErrExisted.Code(), dot.SError.ErrExisted.Error()+m.LiveId.String())
	}
	ms.liveIdMap[m.LiveId] = m

	return nil
}

func (ms *Lives) Remove(m *dot.Live) error {

	delete(ms.liveIdMap, m.LiveId)
	return nil
}

func (ms *Lives) Get(liveId dot.LiveId) (meta *dot.Live, err error) {
	meta = nil
	err = nil

	meta, ok := ms.liveIdMap[liveId]
	if !ok {
		err = dot.NewError(dot.SError.ErrNotExisted.Code(), dot.SError.ErrNotExisted.Error()+liveId.String())
	}
	return
}
