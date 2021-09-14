package wt

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cristalhq/jwt/v3"
	"github.com/google/uuid"
	. "github.com/raver119/statika/classes"
	. "github.com/raver119/statika/utils"
)

type UploadClaims struct {
	jwt.RegisteredClaims
	Buckets []string
}

type Tokenizer struct {
	key []byte
}

func DevTokenizer(key string) Tokenizer {
	//log.Printf("KEY: %v", k)
	return Tokenizer{key: []byte(key)}
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
	if req.Bucket != "" {
		req.Buckets = append(req.Buckets, req.Bucket)
	}

	// now validate buckets
	for _, b := range req.Buckets {
		if b == "" {
			return "", fmt.Errorf("empty bucket was requested")
		}
	}

	claims := &UploadClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Audience: []string{"statika"},
			ID:       uuid.NewString(),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
		Buckets: req.Buckets,
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
	var claims UploadClaims
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
