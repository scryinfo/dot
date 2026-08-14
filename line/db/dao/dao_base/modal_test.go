package daobase

import (
	"testing"

	"github.com/scryinfo/dot/lib/kits"
	"github.com/stretchr/testify/assert"
)

type data struct{}

func (d *data) Prefix() []byte {
	return []byte("data")
}

func TestPrefix(t *testing.T) {
	d := (*data)(nil)
	prefix := d.Prefix()
	assert.Equal(t, []byte("data"), prefix)
}

type data2 struct {
	ModalBase
	Name string `json:"name"`
}

// Key implements [Modal].
func (d *data2) Key() []byte {
	return append(d.Prefix(), kits.UnsafeToBytes(string(d.Id))...)
}

// Prefix implements [Modal].
func (d *data2) Prefix() []byte {
	return []byte("data2:")
}

func (d *data2) Value() ([]byte, error) {
	return Value(d)
}

func (d *data2) FromValue(bs []byte) error {
	return FromValue(d, bs)
}

var _ Modal = (*data2)(nil)

func TestModalBase(t *testing.T) {
	d := data2{
		ModalBase: ModalBase{
			Id: "1",
		},
		Name: "data2",
	}
	bs, _ := d.Value()
	str := string(bs)
	assert.Equal(t, `{"id":"1","name":"data2"}`, str)
}
