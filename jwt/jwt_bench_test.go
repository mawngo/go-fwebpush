package jwt_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/mawngo/go-fwebpush/jwt"
	"io"
	"testing"
	"time"
)

func BenchmarkAlgES(b *testing.B) {
	esAlgos := map[jwt.Algorithm]elliptic.Curve{
		jwt.ES256: elliptic.P256(),
		jwt.ES384: elliptic.P384(),
		jwt.ES512: elliptic.P521(),
	}
	for algo, curve := range esAlgos {
		key, errKey := ecdsa.GenerateKey(curve, rand.Reader)
		if errKey != nil {
			b.Fatal(errKey)
		}
		signer, errSigner := jwt.NewSignerES(algo, key)
		if errSigner != nil {
			b.Fatal(errSigner)
		}
		verifier, errVerifier := jwt.NewVerifierES(algo, &key.PublicKey)
		if errVerifier != nil {
			b.Fatal(errVerifier)
		}

		builder := jwt.NewBuilder(signer)
		b.Run("Sign-"+string(algo), func(b *testing.B) {
			runSignerBench(b, builder)
		})
		b.Run("Verify-"+string(algo), func(b *testing.B) {
			runVerifyBench(b, builder, verifier)
		})
	}
}

func runSignerBench(b *testing.B, builder *jwt.Builder) {
	b.Helper()
	b.ReportAllocs()

	claims := jwt.RegisteredClaims{
		ID:       "id",
		Issuer:   "sdf",
		IssuedAt: time.Now().Unix(),
	}

	var dummy int
	for i := 0; i < b.N; i++ {
		token, err := builder.Build(claims)
		if err != nil {
			b.Fatal(err)
		}
		dummy += int(token.PayloadPart()[0])
	}
	sink(dummy)
}

func runVerifyBench(b *testing.B, builder *jwt.Builder, verifier jwt.Verifier) {
	b.Helper()
	const tokensCount = 32
	tokens := make([]*jwt.Token, 0, tokensCount)
	for i := 0; i < tokensCount; i++ {
		token, err := builder.Build(jwt.RegisteredClaims{
			ID:       "id",
			Issuer:   "sdf",
			IssuedAt: time.Now().Unix(),
		})
		if err != nil {
			b.Fatal(err)
		}
		tokens = append(tokens, token)
	}

	b.ReportAllocs()
	var dummy int
	for i := 0; i < b.N/tokensCount; i++ {
		for _, token := range tokens {
			err := verifier.Verify(token)
			if err != nil {
				b.Fatal(err)
			}
			dummy++
		}
	}
	sink(dummy)
}

func sink(v any) {
	fmt.Fprint(io.Discard, v)
}
