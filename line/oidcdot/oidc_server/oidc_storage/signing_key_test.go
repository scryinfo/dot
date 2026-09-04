package oidc_storage

import (
	"slices"
	"testing"

	"github.com/go-jose/go-jose/v4"
	"github.com/stretchr/testify/assert"
	"github.com/zitadel/oidc/v4/pkg/op"
)

var SigningKeys _SigningKeys = NewSigningKeys()

func TestGetSigningKey(t *testing.T) {
	for _, alg := range SigningKeys.mapKeys {
		it := SigningKeys.GetMapSigningKey(alg.SignatureAlgorithm())
		assert.NotNil(t, it)
	}
}

func BenchmarkMap(b *testing.B) {
	l := len(SigningKeys.mapKeys)
	arrayAlg := _makeArrayAlg()
	mapKeys := _makeMapKeys(arrayAlg)
	i := 0
	b.ResetTimer()
	for b.Loop() {
		index := i % l
		_findMap(mapKeys, arrayAlg[index])
		i++
	}
}
func BenchmarkBinarySearch(b *testing.B) {
	l := len(SigningKeys.mapKeys)
	arrayAlg := _makeArrayAlg()
	arrayKeys := _makeArrayKeys(arrayAlg)
	i := 0
	b.ResetTimer()
	for b.Loop() {
		index := i % l
		_findBinary(arrayKeys, arrayAlg[index])
		i++
	}
}

func BenchmarkArray(b *testing.B) {
	l := len(SigningKeys.mapKeys)

	arrayAlg := _makeArrayAlg()
	arrayKeys := _makeArrayKeys(arrayAlg)
	i := 0
	b.ResetTimer()
	for b.Loop() {
		index := i % l
		_findArray(arrayKeys, arrayAlg[index])
		i++
	}
}

// inline
func _makeArrayAlg() []jose.SignatureAlgorithm {
	l := len(SigningKeys.mapKeys)

	arrayAlg := make([]jose.SignatureAlgorithm, 0, l)
	for _, alg := range SigningKeys.mapKeys {
		arrayAlg = append(arrayAlg, alg.SignatureAlgorithm())
	}
	slices.SortFunc(arrayAlg, func(it, it2 jose.SignatureAlgorithm) int {
		if it == it2 {
			return 0
		} else if it < it2 {
			return -1
		} else {
			return 1
		}
	})
	return arrayAlg
}

// inline
func _makeArrayKeys(arrayAlg []jose.SignatureAlgorithm) []op.SigningKey {
	l := len(arrayAlg)
	arrayKeys := make([]op.SigningKey, 0, l)
	for _, alg := range arrayAlg {
		arrayKeys = append(arrayKeys, SigningKeys.mapKeys[alg])
	}
	return arrayKeys
}

// inline
func _makeMapKeys(arrayAlg []jose.SignatureAlgorithm) map[jose.SignatureAlgorithm]op.SigningKey {
	l := len(arrayAlg)
	mapKeys := make(map[jose.SignatureAlgorithm]op.SigningKey, l)
	for _, alg := range arrayAlg {
		mapKeys[alg] = SigningKeys.mapKeys[alg]
	}
	return mapKeys
}

// inline
func _compare(it op.SigningKey, t jose.SignatureAlgorithm) int {
	itt := it.SignatureAlgorithm()
	if itt == t {
		return 0
	} else if itt < t {
		return -1
	} else {
		return 1
	}
}

// inline
func _findMap(mapKeys map[jose.SignatureAlgorithm]op.SigningKey, alg jose.SignatureAlgorithm) op.SigningKey {
	return mapKeys[alg]
}

// inline
func _findArray(arrayKeys []op.SigningKey, alg jose.SignatureAlgorithm) op.SigningKey {
	for _, key := range arrayKeys {
		if key.SignatureAlgorithm() == alg {
			return key
		}
	}
	return nil
}

// inline
func _findBinary(arrayKeys []op.SigningKey, alg jose.SignatureAlgorithm) op.SigningKey {
	index, find := slices.BinarySearchFunc(arrayKeys, alg, _compare)
	if find {
		return arrayKeys[index]
	} else {
		return nil
	}
}
