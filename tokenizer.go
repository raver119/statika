package main

import (
	"encoding/json"
	"fmt"
	"github.com/cristalhq/jwt/v3"
	"github.com/google/uuid"
	"strings"
	"time"
)

type uploadClaims struct {
	jwt.RegisteredClaims
	Buckets []string
}

type Tokenizer struct {
	key []byte
}

func NewTokenizer() Tokenizer {
	//log.Printf("KEY: %v", k)
	return Tokenizer{key: []byte(GetEnvOrPanic("MASTER_KEY"))}
}

func (t Tokenizer) CreateUploadToken(req UploadAuthenticationRequest) (token string, err error) {
	signer, err := jwt.NewSignerHS(jwt.HS256, t.key)
	if err != nil {
		return "", err
	}

	// TODO: once req.Bucket removed this will be removed as well
	buckets := append(req.Buckets, req.Bucket)

	// now validate buckets
	for _, b := range buckets {
		if b == "" {
			return "", fmt.Errorf("empty bucket was requested")
		}
	}

	claims := &uploadClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Audience: []string{"statika"},
			ID:       uuid.NewString(),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
		Buckets: buckets,
	}

	builder := jwt.NewBuilder(signer)

	tkn, err := builder.Build(claims)
	if err != nil {
		return "", err
	}

	return tkn.String(), nil
}

func (t Tokenizer) ValidateUploadToken(token string, bucket string) (ok bool, err error) {
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.Replace(token, "Bearer ", "", 1)
	}

	verifier, err := jwt.NewVerifierHS(jwt.HS256, t.key)
	if err != nil {
		return false, err
	}

	tkn, err := jwt.ParseString(token)
	if err != nil {
		return false, err
	}

	err = verifier.Verify(tkn.Payload(), tkn.Signature())
	if err != nil {
		return false, err
	}

	// now validate the bucket
	var claims uploadClaims
	err = json.Unmarshal(tkn.RawClaims(), &claims)
	if err != nil {
		return false, err
	}

	// ping is the only request that can get here
	if bucket == "" {
		return true, nil
	}

	for _, b := range claims.Buckets {
		if b == bucket {
			return true, nil
		}
	}

	return false, nil
}
