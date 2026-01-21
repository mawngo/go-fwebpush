package jwt_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	jwt2 "github.com/mawngo/go-fwebpush/internal/jwt"
	"io"
	"testing"
	"time"
)

func BenchmarkAlgES(b *testing.B) {
	esAlgos := map[jwt2.Algorithm]elliptic.Curve{
		jwt2.ES256: elliptic.P256(),
		jwt2.ES384: elliptic.P384(),
		jwt2.ES512: elliptic.P521(),
	}
	for algo, curve := range esAlgos {
		key, errKey := ecdsa.GenerateKey(curve, rand.Reader)
		if errKey != nil {
			b.Fatal(errKey)
		}
		signer, errSigner := jwt2.NewSignerES(algo, key)
		if errSigner != nil {
			b.Fatal(errSigner)
		}
		verifier, errVerifier := jwt2.NewVerifierES(algo, &key.PublicKey)
		if errVerifier != nil {
			b.Fatal(errVerifier)
		}

		builder := jwt2.NewBuilder(signer)
		b.Run("Sign-"+string(algo), func(b *testing.B) {
			runSignerBench(b, builder)
		})
		b.Run("Verify-"+string(algo), func(b *testing.B) {
			runVerifyBench(b, builder, verifier)
		})
	}
}

func runSignerBench(b *testing.B, builder *jwt2.Builder) {
	b.Helper()
	b.ReportAllocs()

	claims := jwt2.RegisteredClaims{
		ID:       "id",
		Issuer:   "sdf",
		IssuedAt: time.Now().Unix(),
	}

	var dummy int
	for range b.N {
		token, err := builder.Build(claims)
		if err != nil {
			b.Fatal(err)
		}
		dummy += int(token.PayloadPart()[0])
	}
	sink(dummy)
}

func runVerifyBench(b *testing.B, builder *jwt2.Builder, verifier jwt2.Verifier) {
	b.Helper()
	const tokensCount = 32
	tokens := make([]*jwt2.Token, 0, tokensCount)
	for range tokensCount {
		token, err := builder.Build(jwt2.RegisteredClaims{
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
	for range b.N / tokensCount {
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
