package fwebpush

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/mawngo/go-fwebpush/fastunsafeurl"
	jwt2 "github.com/mawngo/go-fwebpush/internal/jwt"
	"math/big"
	"time"
)

func (p *VAPIDPusher) getCachedKeys(endpoint string, now time.Time) (reusableKey, error) {
	aud, _, err := fastunsafeurl.ParseSchemeHost(endpoint)
	if err != nil {
		return reusableKey{}, fmt.Errorf("error parsing audience: %w", err)
	}
	// Cache disabled.
	if p.vapidTokenTTL <= 0 {
		auth, err := p.doGenLocalKey()
		if err != nil {
			return reusableKey{}, err
		}
		auth.vapid, auth.exp, err = p.doGetVAPIDAuthorizationHeader(aud, now)
		return auth, err
	}

	// We should regenerate the token <additional time> before actually expires.
	// So the min-acceptable expiration should be <additional time> after now.
	nowExp := now.Add(p.vapidTTLBuffer)
	// Most of the time code will run into this path.
	// Cache hit, not expired, use cached vapid.
	p.mu.RLock()
	if auth := p.cache[aud]; nowExp.Before(auth.exp) {
		p.mu.RUnlock()
		return auth, nil
	}
	p.mu.RUnlock()

	p.mu.Lock()
	defer p.mu.Unlock()
	// Someone else has written to the cache.
	if auth := p.cache[aud]; nowExp.Before(auth.exp) {
		return auth, nil
	}
	auth, err := p.doGenLocalKey()
	if err != nil {
		return reusableKey{}, err
	}
	auth.vapid, auth.exp, err = p.doGetVAPIDAuthorizationHeader(aud, now)
	p.cache[aud] = auth
	return auth, nil
}

func (p *VAPIDPusher) doGetVAPIDAuthorizationHeader(aud string, now time.Time) (string, time.Time, error) {
	// Always expire at least <additional time> (so the message won't expire when it reached the server).
	exp := now.Add(p.vapidTokenTTL + p.vapidTTLBuffer)
	privKey := generateVAPIDHeaderKeys(p.vapidPrivateKey)
	signer, err := jwt2.NewSignerES(jwt2.ES256, privKey)
	if err != nil {
		return "", exp, err
	}
	claims := &jwt2.RegisteredClaims{
		Audience:  aud,
		Subject:   p.subject,
		ExpiresAt: exp.Unix(),
	}
	token, err := jwt2.NewBuilder(signer).Build(claims)
	if err != nil {
		return "", exp, err
	}
	return "vapid t=" + token.String() + p.vapidPublicKeyHeaderPart, exp, nil
}

func (p *VAPIDPusher) doGenLocalKey() (reusableKey, error) {
	curve := ecdh.P256()
	// Application server key pairs (single use).
	localPrivateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return reusableKey{}, errors.Join(ErrEncryption, err)
	}
	localPublicKeyBytes := localPrivateKey.PublicKey().Bytes()
	return reusableKey{
		curve:               curve,
		localPrivateKey:     localPrivateKey,
		localPublicKeyBytes: localPublicKeyBytes,
	}, nil
}

// GenerateVAPIDKeys will create a private and public VAPID key pair.
func GenerateVAPIDKeys() (privateKey, publicKey string, err error) {
	// Get the private key from the P256 curve
	curve := ecdh.P256()

	private, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return
	}

	public := private.PublicKey()
	// Convert to base64
	publicKey = base64.RawURLEncoding.EncodeToString(public.Bytes())
	privateKey = base64.RawURLEncoding.EncodeToString(private.Bytes())
	return
}

// Generates the ECDSA public and private keys for the JWT encryption.
func generateVAPIDHeaderKeys(privateKey []byte) *ecdsa.PrivateKey {
	// Public key
	curve := elliptic.P256()
	px, py := curve.ScalarMult(
		curve.Params().Gx,
		curve.Params().Gy,
		privateKey,
	)

	pubKey := ecdsa.PublicKey{
		Curve: curve,
		X:     px,
		Y:     py,
	}

	// Private key
	d := &big.Int{}
	d.SetBytes(privateKey)

	return &ecdsa.PrivateKey{
		PublicKey: pubKey,
		D:         d,
	}
}

// reusableKey is used to cache the VAPID reusable keys and token.
// Does not modify this struct outside this file.
type reusableKey struct {
	vapid string

	curve               ecdh.Curve
	localPrivateKey     *ecdh.PrivateKey
	localPublicKeyBytes []byte
	exp                 time.Time
}
