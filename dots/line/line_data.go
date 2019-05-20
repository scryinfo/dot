// Scry Info.  All rights reserved.
// license that can be found in the license file.

package line

import (
	"github.com/scryinfo/dot/dot"
)

//Metas
type Metas struct {
	metas map[dot.TypeId]*dot.Metadata
}

//NewMetas
func NewMetas() *Metas {
	m := &Metas{}
	m.metas = make(map[dot.TypeId]*dot.Metadata, 0)
	return m
}

//Lives
type Lives struct {
	// typeIdMap map[dot.TypeId][]*dot.Live
	LiveIdMap map[dot.LiveId]*dot.Live
}

func NewLives() *Lives {
	l := &Lives{}
	l.LiveIdMap = make(map[dot.LiveId]*dot.Live)
	return l
}

func (ms *Metas) Add(m *dot.Metadata) error {
	if m == nil || m.TypeId.String() == "" {
		return dot.SError.NilParameter
	}

	_, ok := ms.metas[m.TypeId]
	if ok {
		return dot.SError.Existed.AddNewError(m.TypeId.String())
	}

	ms.metas[m.TypeId] = m.Clone()

	return nil
}

func (ms *Metas) UpdateOrAdd(m *dot.Metadata) error {
	if m == nil || m.TypeId.String() == "" {
		return dot.SError.NilParameter
	}

	old, ok := ms.metas[m.TypeId]
	if ok {
		old.NewDoter = m.NewDoter
		old.RefType = m.RefType
	} else {
		ms.metas[m.TypeId] = m.Clone()
	}
	return nil
}

func (ms *Metas) Remove(m *dot.Metadata) error {
	delete(ms.metas, m.TypeId)
	return nil
}

func (ms *Metas) RemoveById(typeid dot.TypeId) error {
	delete(ms.metas, typeid)
	return nil
}

func (ms *Metas) Get(typeId dot.TypeId) (meta *dot.Metadata, err error) {
	meta = nil
	err = nil

	meta, ok := ms.metas[typeId]
	if !ok {
		err = dot.SError.NotExisted.AddNewError(typeId.String())
	}
	return
}

func (ms *Metas) NewDot(t dot.TypeId) (dot dot.Dot, err error) {
	dot = nil
	err = nil

	m, err := ms.Get(t)
	if err == nil {
		dot, err = m.NewDot(nil)
	}
	return
}

func (ms *Lives) Add(m *dot.Live) error {
	if m == nil || m.LiveId.String() == "" {
		return dot.SError.NilParameter
	}

	_, ok := ms.LiveIdMap[m.LiveId]
	if ok {
		return dot.SError.Existed.AddNewError(m.LiveId.String())
	}
	ms.LiveIdMap[m.LiveId] = m

	return nil
}

func (ms *Lives) UpdateOrAdd(m *dot.Live) error {
	if m == nil || m.LiveId.String() == "" {
		return dot.SError.NilParameter
	}

	old, ok := ms.LiveIdMap[m.LiveId]
	if ok {
		old.Dot = m.Dot
		old.TypeId = m.TypeId
		if len(m.RelyLives) > 0 { //merge the rely lives
			if old.RelyLives == nil {
				old.RelyLives = make(map[string]dot.LiveId, len(m.RelyLives))
			}
			for k, v := range m.RelyLives {
				old.RelyLives[k] = v
			}
		}
	} else {
		ms.LiveIdMap[m.LiveId] = m
	}

	return nil
}

func (ms *Lives) Remove(m *dot.Live) error {
	delete(ms.LiveIdMap, m.LiveId)
	return nil
}

func (ms *Lives) RemoveById(id dot.LiveId) error {
	delete(ms.LiveIdMap, id)
	return nil
}

func (ms *Lives) Get(liveId dot.LiveId) (meta *dot.Live, err error) {
	meta = nil
	err = nil

	meta, ok := ms.LiveIdMap[liveId]
	if !ok {
		err = dot.SError.NotExisted.AddNewError(liveId.String())
	}
	return
}

func CloneRelyLiveId(old map[string]dot.LiveId) map[string]dot.LiveId {
	re := make(map[string]dot.LiveId, len(old))
	for k, v := range old {
		re[k] = v
	}
	return re
}
