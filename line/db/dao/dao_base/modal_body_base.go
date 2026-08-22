package daobase

import (
	"encoding/json/jsontext"

	"github.com/scryinfo/dot/lib/kits"
)

type ModalBodySample struct {
	ModalBase
	Body jsontext.Value `json:"body"`
}

func NewModalBodySample() ModalBodySample {
	return ModalBodySample{ModalBase: NewModalBase()}
}

func NewModalBodySampleId(id IdType) ModalBodySample {
	return ModalBodySample{ModalBase: NewModalBaseId(id)}
}

// FromValueBody implements [ModalBody].
func (m *ModalBodySample) FromValueBody(bs []byte) error {
	return FromValue(&m.Body, bs)
}

// KeyBody implements [ModalBody].
func (m *ModalBodySample) KeyBody() []byte {
	return append(m.PrefixBody(), kits.StringToBytes(string(m.Id))...)
}

// PrefixBody implements [ModalBody].
func (m *ModalBodySample) PrefixBody() []byte {
	return ([]byte)("modelbodybase")
}

// ValueBody implements [ModalBody].
func (m *ModalBodySample) ValueBody() ([]byte, error) {
	return Value(&m.Body)
}

// FromValue implements [Modal].
func (m *ModalBodySample) FromValue(bs []byte) error {
	return FromValue(m, bs)
}

// Key implements [Modal].
func (m *ModalBodySample) Key() []byte {
	return append(m.Prefix(), kits.StringToBytes(string(m.Id))...)
}

// Prefix implements [Modal].
func (m *ModalBodySample) Prefix() []byte {
	return ([]byte)("modelbase")
}

// Value implements [Modal].
func (m *ModalBodySample) Value() ([]byte, error) {
	return Value(m)
}

var _ ModalBody = (*ModalBodySample)(nil)
var _ Modal = (*ModalBodySample)(nil)
