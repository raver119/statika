package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

func TestS3Storage_Get(t *testing.T) {
	spacesBucket := GetEnvOrDefault("S3_BUCKET", "")
	log.Printf("Bucket: %v", spacesBucket)
	storage, err := NewSpacesStorage(spacesBucket, "https://nyc3.digitaloceanspaces.com")
	require.NoError(t, err)

	content := "test file content"
	f, err := storage.Put("test", "filename.txt", strings.NewReader(content))
	require.NoError(t, err)
	assert.Equal(t, "test/filename.txt", f)

	r, err := storage.Get("test", "filename.txt")
	require.NoError(t, err)

	body, err := ioutil.ReadAll(r)
	require.Equal(t, content, string(body))

	err = storage.Delete("test", "filename.txt")
	require.NoError(t, err)
}
