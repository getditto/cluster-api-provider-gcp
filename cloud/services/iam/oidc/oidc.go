package oidc

import (
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/go-jose/go-jose/v3"
)

const (
	jwksKey                 = "/openid/v1/jwks"
	opendIDConfigurationKey = "/.well-known/openid-configuration"
)

// oidcDiscovery represents the OpenID Connect discovery document.
type oidcDiscovery struct {
	Issuer                string   `json:"issuer"`
	JWKSURI               string   `json:"jwks_uri"`
	AuthorizationEndpoint string   `json:"authorization_endpoint"`
	ResponseTypes         []string `json:"response_types_supported"`
	SubjectTypes          []string `json:"subject_types_supported"`
	SigningAlgs           []string `json:"id_token_signing_alg_values_supported"`
	ClaimsSupported       []string `json:"claims_supported"`
}

// jwksDocument represents a JWKS document.
type jwksDocument struct {
	Keys []jose.JSONWebKey `json:"keys"`
}

func buildDiscoveryJSON(issuerURL string) ([]byte, error) {
	d := oidcDiscovery{
		Issuer:                issuerURL,
		JWKSURI:               fmt.Sprintf("%v/openid/v1/jwks", issuerURL),
		AuthorizationEndpoint: "urn:kubernetes:programmatic_authorization",
		ResponseTypes:         []string{"id_token"},
		SubjectTypes:          []string{"public"},
		SigningAlgs:           []string{"RS256"},
		ClaimsSupported:       []string{"sub", "iss"},
	}
	return json.MarshalIndent(d, "", "")
}

func (s *Service) buildIssuerURL() string {
	// e.g. storage.googleapis.com/<bucketname>/<clustername>-sa
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.scope.Bucket().Name, s.scope.Name())
}

// createJwksKey generates a JSON Web Key (JWK) from the given private key bytes (in PEM format).
// It returns a pointer to the JSONWebKey or an error if the operation fails.
func createJwksKey(privKeyBytes []byte) (*jose.JSONWebKey, error) {
	keyBlock, _ := pem.Decode(privKeyBytes)
	if keyBlock == nil {
		return nil, errors.New("failed to decode PEM block for private key")
	}

	cert, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	keyID, err := keyIDFromPublicKey(&cert.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key ID from public key: %w", err)
	}

	return &jose.JSONWebKey{
		Key:       &cert.PublicKey,
		KeyID:     keyID,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}, nil
}

// keyIDFromPublicKey derives a key ID non-reversibly from a public key.
// taken from: https://github.com/kubernetes/kubernetes/blob/master/pkg/serviceaccount/jwt.go
func keyIDFromPublicKey(publicKey interface{}) (string, error) {
	publicKeyDERBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to serialize public key to DER format: %v", err)
	}

	hasher := crypto.SHA256.New()
	hasher.Write(publicKeyDERBytes)
	publicKeyDERHash := hasher.Sum(nil)

	keyID := base64.RawURLEncoding.EncodeToString(publicKeyDERHash)

	return keyID, nil
}
