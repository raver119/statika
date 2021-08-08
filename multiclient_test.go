package main

import (
	"github.com/google/uuid"
	. "github.com/raver119/statika/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMultiClient_Login(t *testing.T) {
	var ls Storage = NewLocalStorage("/tmp")

	e, err := CreateEngine(masterKey, uploadKey, &ls, 9191)
	require.NoError(t, err)
	require.NoError(t, e.StartAsync())
	time.Sleep(1 * time.Second)
	defer func() {
		_ = e.Stop()
	}()

	bucket1 := uuid.NewString()
	bucket2 := uuid.NewString()
	gk, err := NewGateKeeper(endpoint, masterKey, uploadKey)
	require.NoError(t, err)

	token, err := gk.IssueUploadToken(bucket1, bucket2)
	require.NoError(t, err)

	client, err := gk.NewMultiClient(token, bucket1, bucket2)
	require.NoError(t, err)

	assert.True(t, client.Check(bucket1))
	assert.True(t, client.Check(bucket2))
	assert.False(t, client.Check(uuid.NewString()))

	// test for non-approved access
	_, err = gk.NewMultiClient(token, bucket1, uuid.NewString())
	require.Error(t, err)
}
