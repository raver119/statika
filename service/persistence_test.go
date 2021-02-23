package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	host = GetEnvOrDefault("MEMCACHED_HOST", "localhost")
)

/*
	PLEASE NOTE: This set of tests requires memcached instance to be running on MEMCACHED_HOST env var (which defaults to 127.0.0.1)
*/

func TestPersistenceAgent_CreateUploadToken(t *testing.T) {
	pa, err := NewPersistenceAgent(host, 11211)
	require.NoError(t, err)

	req := UploadAuthenticationRequest{Token: "SOME_TOKEN", Bucket: "SOME_BUCKET"}
	token, err := pa.CreateUploadToken(req)
	require.NoError(t, err)

	assert.True(t, pa.CheckUploadToken(token, req.Bucket))
	assert.False(t, pa.CheckUploadToken(token, "random bucked name"))
	assert.False(t, pa.CheckUploadToken("random access token", req.Bucket))
	assert.False(t, pa.CheckUploadToken("random access token", "random bucked name"))
	assert.True(t, pa.TouchUploadToken(token))
}
