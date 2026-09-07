package oidc_storage

import (
	"slices"
	"testing"
)

func BenchmarkScopeMap(b *testing.B) {
	l := len(_arrayScope)
	i := 0
	b.ResetTimer()
	for b.Loop() {
		index := i % l
		_findScopeMap(_arrayScope[index])
		i++
	}
}
func BenchmarkScopeBinarySearch(b *testing.B) {
	l := len(_arrayScope)
	i := 0
	b.ResetTimer()
	for b.Loop() {
		index := i % l
		_findScopeBinary(_arrayScope[index])
		i++
	}
}

func BenchmarkScopeArray(b *testing.B) {
	l := len(_arrayScope)
	i := 0
	b.ResetTimer()
	for b.Loop() {
		index := i % l
		_findScopeArray(_arrayScope[index])
		i++
	}
}

var (
	_mapScope   = map[string]struct{}{}
	_arrayScope = ScopeEx.EnumsGo()
)

func init() {
	slices.Sort(_arrayScope)
	for _, s := range _arrayScope {
		_mapScope[s] = struct{}{}
	}
}

func _findScopeMap(s string) bool {
	_, ok := _mapScope[s]
	return ok
}

func _findScopeArray(s string) bool {
	return slices.Contains(_arrayScope, s)
}

func _findScopeBinary(s string) bool {
	_, ok := slices.BinarySearch(_arrayScope, s)
	return ok
}
