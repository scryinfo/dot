// Scry Info.  All rights reserved.
// license that can be found in the license file.

package line

import (
	"github.com/scryinfo/dot/dot"
)

// Metas metas of dot
type Metas struct {
	metas map[dot.TypeID]*dot.Metadata
}

// NewMetas new metas
func NewMetas() *Metas {
	m := &Metas{}
	m.metas = make(map[dot.TypeID]*dot.Metadata)
	return m
}

// Lives lives of dot
type Lives struct {
	// typeIDMap map[dot.TypeID][]*dot.Live
	LiveIDMap map[dot.LiveID]*dot.Live
}

// NewLives new lives
func NewLives() *Lives {
	l := &Lives{}
	l.LiveIDMap = make(map[dot.LiveID]*dot.Live)
	return l
}

// Add add meta data to metas
func (ms *Metas) Add(m *dot.Metadata) error {
	if m == nil || m.TypeID.String() == "" {
		return dot.SError.NilParameter
	}

	if _, ok := ms.metas[m.TypeID]; ok {
		return dot.SError.Existed.AddNewError(m.TypeID.String())
	}

	ms.metas[m.TypeID] = m.Clone()

	return nil
}

// UpdateOrAdd if the meta exist then update, or add it
func (ms *Metas) UpdateOrAdd(m *dot.Metadata) error {
	if m == nil || m.TypeID.String() == "" {
		return dot.SError.NilParameter
	}

	old, ok := ms.metas[m.TypeID]
	if ok {
		old.Merge(m)
	} else {
		ms.metas[m.TypeID] = m.Clone()
	}
	return nil
}

// RemoveByID remove meta data from metas
func (ms *Metas) RemoveByID(typeID dot.TypeID) error {
	delete(ms.metas, typeID)
	return nil
}

// Get get meta by type id
func (ms *Metas) Get(typeID dot.TypeID) (meta *dot.Metadata, err error) {
	meta = nil
	err = nil

	meta, ok := ms.metas[typeID]
	if !ok {
		err = dot.SError.NotExisted.AddNewError(typeID.String())
	}
	return
}

// NewDot new dot
func (ms *Metas) NewDot(t dot.TypeID) (dot dot.Dot, err error) {
	dot = nil
	err = nil

	m, err := ms.Get(t)
	if err == nil {
		dot, err = m.NewDot(nil)
	}
	return
}

// Add add live
func (ms *Lives) Add(m *dot.Live) error {
	if m == nil || m.LiveID.String() == "" {
		return dot.SError.NilParameter
	}

	_, ok := ms.LiveIDMap[m.LiveID]
	if ok {
		return dot.SError.Existed.AddNewError(m.LiveID.String())
	}
	ms.LiveIDMap[m.LiveID] = m

	return nil
}

// UpdateOrAdd if the live exist then update, or add it
func (ms *Lives) UpdateOrAdd(m *dot.Live) error {
	if m == nil || m.LiveID.String() == "" {
		return dot.SError.NilParameter
	}

	old, ok := ms.LiveIDMap[m.LiveID]
	if ok {
		old.Dot = m.Dot
		old.TypeID = m.TypeID
		if len(m.RelyLives) > 0 { //merge the rely lives
			if old.RelyLives == nil {
				old.RelyLives = make(map[string]dot.LiveID, len(m.RelyLives))
			}
			for k, v := range m.RelyLives {
				old.RelyLives[k] = v
			}
		}
	} else {
		ms.LiveIDMap[m.LiveID] = m
	}

	return nil
}

// RemoveByID remove the dot by live id
func (ms *Lives) RemoveByID(id dot.LiveID) error {
	delete(ms.LiveIDMap, id)
	return nil
}

// Get get dot by live id
func (ms *Lives) Get(liveID dot.LiveID) (meta *dot.Live, err error) {
	var ok bool
	meta, ok = ms.LiveIDMap[liveID]
	if !ok {
		err = dot.SError.NotExisted.AddNewError(liveID.String())
	}
	return
}
