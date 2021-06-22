package scrypto

import (
	"crypto"
	"io"
)

type CurveParameters struct {
	Name    string // the canonical name of the curve
	BitSize int    // the size of the underlying field
}

// Ecdh is the interface defining all functions
// necessary for ECDH.
type Ecdh interface {
	// GenerateKey generates a private/public key pair using entropy from rand.
	// If rand is nil, crypto/rand.Reader will be used.
	GenerateKey(rand io.Reader) (private crypto.PrivateKey, public crypto.PublicKey, err error)

	// Parameters returns the curve parameters - like the field size.
	Parameters() CurveParameters

	// PublicKey returns the public key corresponding to the given private one.
	PublicKey(private crypto.PrivateKey) (public crypto.PublicKey, err error)

	//// Check returns a non-nil error if the peers public key cannot used for the
	//// key exchange - for instance the public key isn't a point on the elliptic curve.
	//// It's recommended to check peer's public key before computing the secret.
	//Check(peersPublic crypto.PublicKey) (err error)

	// ComputeSecret returns the secret value computed from the given private key
	// and the peers public key.
	ComputeSecret(private crypto.PrivateKey, peersPublic crypto.PublicKey) (secret []byte, err error)

	PublicKeyToBytes(publicKey crypto.PublicKey) (key []byte, err error)

	PrivateKeyToBytes(privateKey crypto.PrivateKey) (key []byte, err error)
}
