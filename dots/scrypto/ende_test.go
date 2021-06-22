package scrypto

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEndeData(t *testing.T) {
	data := EndeData{
		PublicKey: []byte{1},
		EndeType:  "test",
		Hash:      []byte{2},
		Signature: []byte{3},
		EnData:    false,
		Body:      []byte{10},
	}

	{
		bytes, err := json.Marshal(&data)
		assert.Equal(t, nil, err)
		assert.NotNil(t, bytes)

		var jData EndeData
		err = json.Unmarshal(bytes, &jData)
		assert.Equal(t, nil, err)

		assert.Equal(t, data, jData)
	}
	{
		bytes, err := json.Marshal(data)
		assert.Equal(t, nil, err)
		assert.NotNil(t, bytes)

		var jData EndeData
		err = json.Unmarshal(bytes, &jData)
		assert.Equal(t, nil, err)

		assert.Equal(t, data, jData)
	}
}
