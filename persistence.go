package main

import (
	"encoding/base64"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/uuid"
	"strings"
)

const (
	PrefixUpload = "UPLOAD_"
	PrefixMaster = "MASTER_"
)

type PersistenceAgent struct {
	memCached *memcache.Client
}

func NewPersistenceAgent(memcachedHost string, port int) (a PersistenceAgent, err error) {
	memCache := memcache.New(fmt.Sprintf("%v:%v", memcachedHost, port))
	a = PersistenceAgent{memCached: memCache}
	err = memCache.Ping()
	return
}

func (a PersistenceAgent) CreateUploadToken(req UploadAuthenticationRequest) (authToken string, err error) {
	authToken = uuid.New().String() + "-" + uuid.New().String()
	encAuthToken := base64.StdEncoding.EncodeToString([]byte(a.normalizeToken(authToken)))

	// TODO: make expiration time configurable?
	err = a.memCached.Set(&memcache.Item{
		Key:        PrefixUpload + encAuthToken,
		Value:      []byte(req.Bucket),
		Expiration: a.Expiration(),
	})

	return
}

func (a PersistenceAgent) CreateMasterToken(req UploadAuthenticationRequest) (authToken string, err error) {
	authToken = uuid.New().String() + "-" + uuid.New().String()
	encAuthToken := base64.StdEncoding.EncodeToString([]byte(a.normalizeToken(authToken)))

	// TODO: make expiration time configurable?
	err = a.memCached.Set(&memcache.Item{
		Key:        PrefixMaster + encAuthToken,
		Value:      []byte(req.Bucket),
		Expiration: a.Expiration(),
	})

	return
}

func (a PersistenceAgent) CheckUploadToken(authToken string, bucket string) bool {
	encAuthToken := base64.StdEncoding.EncodeToString([]byte(a.normalizeToken(authToken)))
	b, err := a.memCached.Get(PrefixUpload + encAuthToken)
	if err != nil {
		return false
	}
	allowedBucket := string(b.Value)
	if b == nil || bucket != allowedBucket {
		return false
	}

	return true
}

func (a PersistenceAgent) normalizeToken(authToken string) string {
	if strings.HasPrefix(authToken, "Bearer ") {
		authToken = strings.Replace(authToken, "Bearer ", "", 1)
	}

	return authToken
}

func (a PersistenceAgent) TouchUploadToken(authToken string) bool {
	encAuthToken := base64.StdEncoding.EncodeToString([]byte(a.normalizeToken(authToken)))
	err := a.memCached.Touch(PrefixUpload+encAuthToken, a.Expiration())
	return err == nil
}

func (a PersistenceAgent) TouchMasterToken(authToken string) (err error) {
	encAuthToken := base64.StdEncoding.EncodeToString([]byte(a.normalizeToken(authToken)))
	err = a.memCached.Touch(PrefixMaster+encAuthToken, 3600)
	return
}

func (a PersistenceAgent) Expiration() int32 {
	return 3600
}
