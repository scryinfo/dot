package daobase

import "github.com/scryinfo/dot/lib/kits"

type ModalSample struct {
	ModalBase
	// add common fields if needed
}

func NewModalSample() ModalSample {
	return ModalSample{
		ModalBase: NewModalBase(),
	}
}

func NewModalSampleId(id IdType) ModalSample {
	return ModalSample{ModalBase: NewModalBaseId(id)}
}

// FromValue implements [Modal].
func (m *ModalSample) FromValue(bs []byte) error {
	return FromValue(m, bs)
	// return FromValue(m, bs)
}

// Key implements [Modal].
func (m *ModalSample) Key() []byte {
	return append(m.Prefix(), kits.UnsafeToBytes(string(m.Id))...)
}

// Prefix implements [Modal].
func (m *ModalSample) Prefix() []byte {
	return ([]byte)("modelbase")
}

// Value implements [Modal].
func (m *ModalSample) Value() ([]byte, error) {
	return Value(m)
}

var _ Modal = (*ModalSample)(nil)
