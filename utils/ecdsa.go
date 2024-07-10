package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%064x%064x", s.R, s.S)
}

func String2BigIntTuples(s string) (big.Int, big.Int) {
	bx, _ := hex.DecodeString(s[:64])
	by, _ := hex.DecodeString(s[64:])

	var bix big.Int
	var biy big.Int

	_ = bix.SetBytes(bx)
	_ = biy.SetBytes(by)

	return bix, biy
}

func String2Signature(s string) *Signature {
	x, y := String2BigIntTuples(s)
	return &Signature{&x, &y}
}

func String2PublicKey(s string) *ecdsa.PublicKey {
	x, y := String2BigIntTuples(s)

	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}
}

func String2PrivateKey(s string, publicKey *ecdsa.PublicKey) *ecdsa.PrivateKey {
	b, _ := hex.DecodeString(s[:])

	var bi big.Int
	_ = bi.SetBytes(b)
	return &ecdsa.PrivateKey{PublicKey: *publicKey, D: &bi}
}
