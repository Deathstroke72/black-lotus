// file: internal/interfaces/http/middleware/jwt_middleware_crypto.go
package middleware

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"hash"
)

// Actual crypto implementations

func init() {
	sha256Hash = func() hash.Hash {
		return sha256.New()
	}
	rsaVerifyPKCS1v15 = func(pub *rsa.PublicKey, hashed, sig []byte) error {
		return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed, sig)
	}
}

var (
	sha256Hash        func() hash.Hash
	rsaVerifyPKCS1v15 func(pub *rsa.PublicKey, hashed, sig []byte) error
)