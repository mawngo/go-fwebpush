package jwt

import (
	"testing"
)

const simplePayload = `simple-string-payload`

func TestSignerAlg(t *testing.T) {
	testCases := []struct {
		s    Signer
		want Algorithm
	}{
		{must(NewSignerES(ES256, ecdsaPrivateKey256)), ES256},
		{must(NewSignerES(ES384, ecdsaPrivateKey384)), ES384},
		{must(NewSignerES(ES512, ecdsaPrivateKey521)), ES512},
	}

	for _, tc := range testCases {
		mustEqual(t, tc.s.Algorithm(), tc.want)
	}
}

func TestVerifierAlg(t *testing.T) {
	testCases := []struct {
		v    Verifier
		want Algorithm
	}{
		{must(NewVerifierES(ES256, ecdsaPublicKey256)), ES256},
		{must(NewVerifierES(ES384, ecdsaPublicKey384)), ES384},
		{must(NewVerifierES(ES512, ecdsaPublicKey521)), ES512},
	}

	for _, tc := range testCases {
		mustEqual(t, tc.v.Algorithm(), tc.want)
	}
}

func TestSignerBadParams(t *testing.T) {
	testCases := []struct {
		err error
	}{
		{getErr(NewSignerES("xxx", ecdsaPrivateKey256))},
	}

	for _, tc := range testCases {
		mustFail(t, tc.err)
	}
}

func TestVerifierBadParams(t *testing.T) {
	testCases := []struct {
		err error
	}{
		{getErr(NewVerifierES("xxx", ecdsaPublicKey256))},
	}

	for _, tc := range testCases {
		mustFail(t, tc.err)
	}
}
