package jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"reflect"
	"testing"
)

func TestMarshalHeader(t *testing.T) {
	testCases := []struct {
		h    *Header
		want string
	}{
		{
			&Header{Algorithm: RS256},
			`{"alg":"RS256"}`,
		},
		{
			&Header{Algorithm: RS256, Type: "JWT"},
			`{"alg":"RS256","typ":"JWT"}`,
		},
		{
			&Header{Algorithm: RS256, ContentType: "token"},
			`{"alg":"RS256","cty":"token"}`,
		},
		{
			&Header{Algorithm: RS256, Type: "JWT", ContentType: "token"},
			`{"alg":"RS256","typ":"JWT","cty":"token"}`,
		},
		{
			&Header{Algorithm: RS256, Type: "JwT", ContentType: "token"},
			`{"alg":"RS256","typ":"JwT","cty":"token"}`,
		},
		{
			&Header{Algorithm: RS256, Type: "JwT", ContentType: "token", KeyID: "test"},
			`{"alg":"RS256","typ":"JwT","cty":"token","kid":"test"}`,
		},
	}

	for _, tc := range testCases {
		raw, err := tc.h.MarshalJSON()
		mustOk(t, err)
		mustEqual(t, string(raw), tc.want)
	}
}

func TestNewKey(t *testing.T) {
	key, err := GenerateRandomBits(512)
	mustOk(t, err)

	// 8 bits to 1 byte
	const byteCount = int(512.0 / 8)
	mustEqual(t, len(key), byteCount)
}

func getErr[T any](_ T, err error) error {
	return err
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func mustParseECKey(s string) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(s))
	if block == nil {
		panic("invalid PEM")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return key
}

func mustOk(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatal(err)
	}
}

func mustFail(tb testing.TB, err error) {
	tb.Helper()
	if err == nil {
		tb.Fatal()
	}
}

func mustEqual[T any](tb testing.TB, have, want T) {
	tb.Helper()
	if !reflect.DeepEqual(have, want) {
		tb.Fatalf("\nhave: %+v\nwant: %+v\n", have, want)
	}
}
