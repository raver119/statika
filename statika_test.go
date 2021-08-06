package main

import (
	"github.com/raver119/statika/api"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	endpoint  = "http://localhost:9191"
	masterKey = "TEST_MASTER_KEY"
	uploadKey = "TEST_UPLOAD_KEY"
)

func TestGateKeeper_IssueUploadToken(t *testing.T) {
	var ls Storage = NewLocalStorage("/tmp")
	e, err := CreateEngine(masterKey, uploadKey, &ls, 9191)
	require.NoError(t, err)
	require.NoError(t, e.StartAsync())
	time.Sleep(1 * time.Second)

	defer func() {
		_ = e.Stop()
	}()

	gk, err := api.NewGateKeeper(endpoint, masterKey, uploadKey)
	require.NoError(t, err)

	token, err := gk.IssueUploadToken("test_bucket")
	require.NoError(t, err)
	//log.Printf("Token: <%v>", token)

	// now confirm token is actually valid and has access to the bucket
	tokenizer := NewTokenizer()
	ok, err := tokenizer.ValidateUploadToken(string(token), "test_bucket")
	require.NoError(t, err)
	require.True(t, ok)

	// now test bad bucket
	ok, err = tokenizer.ValidateUploadToken(string(token), "wrong_bucket")
	require.NoError(t, err)
	require.False(t, ok)

	// now bad token
	ok, err = tokenizer.ValidateUploadToken("ab"+string(token), "test_bucket")
	require.Error(t, err)
	require.False(t, ok)
}

func TestGateKeeper_IssueUploadToken_2(t *testing.T) {
	var ls Storage = NewLocalStorage("/tmp")
	e, err := CreateEngine(masterKey, uploadKey, &ls, 9191)
	require.NoError(t, err)
	require.NoError(t, e.StartAsync())
	time.Sleep(1 * time.Second)
	defer func() {
		_ = e.Stop()
	}()

	gk, err := api.NewGateKeeper(endpoint, masterKey, "bad key")
	require.NoError(t, err)

	_, err = gk.IssueUploadToken("test_bucket")
	require.Error(t, err)
}
