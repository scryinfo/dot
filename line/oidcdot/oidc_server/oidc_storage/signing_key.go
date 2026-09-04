package oidc_storage

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"uuid"

	"github.com/go-jose/go-jose/v4"
	"github.com/scryinfo/dot/dot"
	"github.com/zitadel/oidc/v4/pkg/op"
)

// the map is performance optimized for binary search
type _SigningKeys struct {
	// arrayKeys []op.SigningKey
	mapKeys map[jose.SignatureAlgorithm]op.SigningKey
}

// func (s *_SigningKeys) GetArraySigningKey(alg jose.SignatureAlgorithm) op.SigningKey {
// 	index, find := slices.BinarySearchFunc(s.arrayKeys, alg, func(it op.SigningKey, t jose.SignatureAlgorithm) int {
// 		itt := it.SignatureAlgorithm()
// 		if itt == t {
// 			return 0
// 		} else if itt < t {
// 			return -1
// 		} else {
// 			return 1
// 		}
// 	})
// 	if find {
// 		return s.arrayKeys[index]
// 	} else {
// 		dot.Logger.Debug().Msgf("no signing key found for algorithm %s", alg)
// 		return nil
// 	}
// }

func (s *_SigningKeys) GetMapSigningKey(alg jose.SignatureAlgorithm) op.SigningKey {
	v, find := s.mapKeys[alg]
	if find {
		return v
	} else {
		dot.Logger.Debug().Msgf("no signing key found for algorithm %s", alg)
		return nil
	}
}

var SigningKeys _SigningKeys = _SigningKeys{
	// arrayKeys: []op.SigningKey{
	// 	SigningKeyES256,
	// 	SigningKeyES384,
	// 	SigningKeyES512,
	// 	SigningKeyEdDSA,
	// 	SigningKeyHS256,
	// 	SigningKeyHS384,
	// 	SigningKeyHS512,
	// 	SigningKeyPS256,
	// 	SigningKeyPS384,
	// 	SigningKeyPS512,
	// 	SigningKeyRS256,
	// 	SigningKeyRS384,
	// 	SigningKeyRS512,
	// },
	mapKeys: map[jose.SignatureAlgorithm]op.SigningKey{
		SigningKeyES256.SignatureAlgorithm(): SigningKeyES256,
		SigningKeyES384.SignatureAlgorithm(): SigningKeyES384,
		SigningKeyES512.SignatureAlgorithm(): SigningKeyES512,
		SigningKeyEdDSA.SignatureAlgorithm(): SigningKeyEdDSA,
		SigningKeyHS256.SignatureAlgorithm(): SigningKeyHS256,
		SigningKeyHS384.SignatureAlgorithm(): SigningKeyHS384,
		SigningKeyHS512.SignatureAlgorithm(): SigningKeyHS512,
		SigningKeyPS256.SignatureAlgorithm(): SigningKeyPS256,
		SigningKeyPS384.SignatureAlgorithm(): SigningKeyPS384,
		SigningKeyPS512.SignatureAlgorithm(): SigningKeyPS512,
		SigningKeyRS256.SignatureAlgorithm(): SigningKeyRS256,
		SigningKeyRS384.SignatureAlgorithm(): SigningKeyRS384,
		SigningKeyRS512.SignatureAlgorithm(): SigningKeyRS512,
	},
}

// EdDSA
var _ op.SigningKey = (*signingKeyEdDSA)(nil)

var (
	SigningKeyEdDSA = &signingKeyEdDSA{
		id:        uuid.NewV7().String(),
		algorithm: jose.EdDSA,
		key:       mustEdDSAKey(ed25519.GenerateKey(rand.Reader)),
	}
)

func mustEdDSAKey(_ ed25519.PublicKey, key ed25519.PrivateKey, _ error) *ed25519.PrivateKey {
	return &key
}

type signingKeyEdDSA struct {
	id        string
	algorithm jose.SignatureAlgorithm
	key       *ed25519.PrivateKey
}

func (s *signingKeyEdDSA) SignatureAlgorithm() jose.SignatureAlgorithm {
	return s.algorithm
}

func (s *signingKeyEdDSA) Key() any {
	return s.key
}

func (s *signingKeyEdDSA) ID() string {
	return s.id
}

// rsa256 rsa384 rsa512
var _ op.SigningKey = (*signingKeyRsa)(nil)

var (
	SigningKeyRS256 = &signingKeyRsa{
		id:        uuid.NewV7().String(),
		algorithm: jose.RS256,
		key:       mustRSAKey(rsa.GenerateKey(rand.Reader, 2048)),
	}
	SigningKeyRS384 = &signingKeyRsa{
		id:        uuid.NewV7().String(),
		algorithm: jose.RS384,
		key:       mustRSAKey(rsa.GenerateKey(rand.Reader, 3072)),
	}
	SigningKeyRS512 = &signingKeyRsa{
		id:        uuid.NewV7().String(),
		algorithm: jose.RS512,
		key:       mustRSAKey(rsa.GenerateKey(rand.Reader, 4096)),
	}
)

func mustRSAKey(key *rsa.PrivateKey, _ error) *rsa.PrivateKey {
	return key
}

type signingKeyRsa struct {
	id        string
	algorithm jose.SignatureAlgorithm
	key       *rsa.PrivateKey
}

func (s *signingKeyRsa) SignatureAlgorithm() jose.SignatureAlgorithm {
	return s.algorithm
}

func (s *signingKeyRsa) Key() any {
	return s.key
}

func (s *signingKeyRsa) ID() string {
	return s.id
}

// es256 es384 es512
var _ op.SigningKey = (*signingKeyEs)(nil)

var (
	SigningKeyES256 = &signingKeyEs{
		id:        uuid.NewV7().String(),
		algorithm: jose.ES256,
		key:       mustESKey(ecdsa.GenerateKey(elliptic.P256(), rand.Reader)),
	}
	SigningKeyES384 = &signingKeyEs{
		id:        uuid.NewV7().String(),
		algorithm: jose.ES384,
		key:       mustESKey(ecdsa.GenerateKey(elliptic.P384(), rand.Reader)),
	}
	SigningKeyES512 = &signingKeyEs{
		id:        uuid.NewV7().String(),
		algorithm: jose.ES512,
		key:       mustESKey(ecdsa.GenerateKey(elliptic.P521(), rand.Reader)),
	}
)

func mustESKey(key *ecdsa.PrivateKey, _ error) *ecdsa.PrivateKey {
	return key
}

type signingKeyEs struct {
	id        string
	algorithm jose.SignatureAlgorithm
	key       *ecdsa.PrivateKey
}

func (s *signingKeyEs) SignatureAlgorithm() jose.SignatureAlgorithm {
	return s.algorithm
}

func (s *signingKeyEs) Key() any {
	return s.key
}

func (s *signingKeyEs) ID() string {
	return s.id
}

// hs256 hs384 hs512
var _ op.SigningKey = (*signingKeyHs)(nil)

var (
	SigningKeyHS256 = &signingKeyHs{
		id:        uuid.NewV7().String(),
		algorithm: jose.HS256,
		key:       mustHSKey(32),
	}
	SigningKeyHS384 = &signingKeyHs{
		id:        uuid.NewV7().String(),
		algorithm: jose.HS384,
		key:       mustHSKey(48),
	}
	SigningKeyHS512 = &signingKeyHs{
		id:        uuid.NewV7().String(),
		algorithm: jose.HS512,
		key:       mustHSKey(64),
	}
)

func mustHSKey(len int) []byte {
	key := make([]byte, len)
	_, _ = rand.Read(key)
	return key
}

type signingKeyHs struct {
	id        string
	algorithm jose.SignatureAlgorithm
	key       []byte
}

func (s *signingKeyHs) SignatureAlgorithm() jose.SignatureAlgorithm {
	return s.algorithm
}

func (s *signingKeyHs) Key() any {
	return s.key
}

func (s *signingKeyHs) ID() string {
	return s.id
}

// ps256 ps384 ps512
var _ op.SigningKey = (*signingKeyPs)(nil)

var (
	SigningKeyPS256 = &signingKeyPs{
		id:        uuid.NewV7().String(),
		algorithm: jose.PS256,
		key:       mustPSKey(rsa.GenerateKey(rand.Reader, 2048)),
	}
	SigningKeyPS384 = &signingKeyPs{
		id:        uuid.NewV7().String(),
		algorithm: jose.PS384,
		key:       mustPSKey(rsa.GenerateKey(rand.Reader, 3072)),
	}
	SigningKeyPS512 = &signingKeyPs{
		id:        uuid.NewV7().String(),
		algorithm: jose.PS512,
		key:       mustPSKey(rsa.GenerateKey(rand.Reader, 4096)),
	}
)

func mustPSKey(key *rsa.PrivateKey, _ error) *rsa.PrivateKey {
	return key
}

type signingKeyPs struct {
	id        string
	algorithm jose.SignatureAlgorithm
	key       *rsa.PrivateKey
}

func (s *signingKeyPs) SignatureAlgorithm() jose.SignatureAlgorithm {
	return s.algorithm
}

func (s *signingKeyPs) Key() any {
	return s.key
}

func (s *signingKeyPs) ID() string {
	return s.id
}
