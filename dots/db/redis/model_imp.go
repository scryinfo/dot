package redis

import (
    "encoding/json"
    "github.com/albrow/zoom"
)

type ModelImp struct {
    Id   string
    Name string
}

const defaultModelName = "defaultModelName"

var _ Model = (*ModelImp)(nil)

func (m *ModelImp) Marshal(v interface{}) ([]byte, error) {
    return json.Marshal(v)
}

func (m *ModelImp) Unmarshal(data []byte, v interface{}) error {
    return json.Unmarshal(data, v)
}

func (m *ModelImp) ModelID() string {
    if m.Id == "" {
        m.Id = (&zoom.RandomID{}).ModelID()
    }

    return m.Id
}

func (m *ModelImp) SetModelID(id string) {
    m.Id = id
}

func (m *ModelImp) ModelName() string {
    if m.Name == "" {
        m.Name = defaultModelName
    }

    return m.Name
}

func (m *ModelImp) SetModelName(name string) {
    m.Name = name
}
