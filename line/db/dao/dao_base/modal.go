package daobase

import (
	"encoding/json"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/kits"
)

type IdType string

type Modal interface {
	GetId() IdType
	SetId(id IdType)
	Prefix() []byte //dont use any field in this methos, it is a static method
	Key() []byte
	Value() ([]byte, error)
	FromValue(bs []byte) error
}

type ModalBody interface {
	PrefixBody() []byte //dont use any field in this methos, it is a static method
	KeyBody() []byte
	ValueBody() ([]byte, error)
	FromValueBody(bs []byte) error
}

type ModalPtr[T any] interface {
	*T
	Modal
}
type ModalBodyPtr[T any] interface {
	*T
	Modal
	ModalBody
}

func ModalPrefix[T any, PT ModalPtr[T]]() []byte {
	return ((PT)(nil)).Prefix()
}

func ModalKey[T any, PT ModalPtr[T]](id IdType) []byte {
	return append(ModalPrefix[T, PT](), kits.UnsafeToBytes(string(id))...)
}

func ModalPrefixBody[T any, PT interface {
	*T
	ModalBody
}]() []byte {
	return ((PT)(nil)).PrefixBody()
}

func ModalKeyBody[T any, PT interface {
	*T
	ModalBody
}](id IdType) []byte {
	return append(ModalPrefixBody[T, PT](), kits.UnsafeToBytes(string(id))...)
}

// inline
func FromValue[T any](m *T, bs []byte) error {
	err := json.Unmarshal(bs, m)
	if err != nil {
		dot.Logger.Error().AnErr("unmarshal error: ", err).Send()
	}
	return err
}

// inline
func Value[T any](m *T) ([]byte, error) {
	return json.Marshal(m)
}

type ModalBase struct {
	Id IdType `json:"id"`
}

func NewModalBase() ModalBase {
	return ModalBase{Id: IdType(kits.Ids.NewXId())}
}

func NewModalBaseId(id IdType) ModalBase {
	return ModalBase{Id: id}
}

// GetId implements [Modal].
func (m *ModalBase) GetId() IdType {
	return m.Id
}

// SetId implements [Modal].
func (m *ModalBase) SetId(id IdType) {
	m.Id = id
}
