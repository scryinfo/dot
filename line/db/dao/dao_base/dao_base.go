package daobase

type Dao[T any, PT ModalPtr[T]] interface {
	Add(m PT) error
	Remove(m PT) error
	RemoveBy(id IdType) error
	Find(id IdType) (T, error)
	PrevPage(page *PageMeta) ([]T, error)
	NextPage(page *PageMeta) ([]T, error)
}

type DaoBody[T any, PT ModalBodyPtr[T]] interface {
	Add(m PT) error
	Remove(m PT) error
	RemoveBy(id IdType) error
	Find(id IdType) (T, error)
	PrevPage(page *PageMeta) ([]T, error)
	NextPage(page *PageMeta) ([]T, error)
}

type DaoKey[T any] interface {
	Set(m *T) error
	Get() (*T, error)
	Remove() error
	Key() []byte
}
