package main

import (
	"github.com/raver119/statika/classes"
	"github.com/raver119/statika/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"strings"
	"testing"
)

func TestS3Storage_Get(t *testing.T) {
	spacesBucket := utils.GetEnvOrDefault("S3_BUCKET", "")
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

func TestS3Storage_List(t *testing.T) {
	spacesBucket := utils.GetEnvOrDefault("S3_BUCKET", "")
	storage, err := NewSpacesStorage(spacesBucket, "https://nyc3.digitaloceanspaces.com")
	require.NoError(t, err)

	content := "test file content"
	_, err = storage.Put("test", "filename1.txt", strings.NewReader(content))
	require.NoError(t, err)

	_, err = storage.Put("test", "filename2.txt", strings.NewReader(content))
	require.NoError(t, err)

	_, err = storage.Put("test3", "filename3.txt", strings.NewReader(content))
	require.NoError(t, err)

	f, err := storage.List("test")
	assert.Equal(t, 2, len(f))

	assert.NoError(t, storage.Delete("test", "filename1.txt"))
	assert.NoError(t, storage.Delete("test", "filename2.txt"))
	assert.NoError(t, storage.Delete("test3", "filename3.txt"))
}

func TestS3Storage_GetMeta(t *testing.T) {
	spacesBucket := utils.GetEnvOrDefault("S3_BUCKET", "")
	storage, err := NewSpacesStorage(spacesBucket, "https://nyc3.digitaloceanspaces.com")
	require.NoError(t, err)

	meta := classes.MetaInfo{
		"alpha": "1",
		"beta":  "2",
	}

	err = storage.PutMeta("test", "filename.txt", meta)
	require.NoError(t, err)

	restored, err := storage.GetMeta("test", "filename.txt")
	assert.Equal(t, meta, restored)
	assert.Equal(t, "2", meta["beta"])

	require.NoError(t, storage.DeleteMeta("test", "filename.txt"))
}

func Test_validMode(t *testing.T) {
	tests := []struct {
		name string
		mode string
		want bool
	}{
		{"test_0", "private", true},
		{"test_1", "public-read", true},
		{"test_2", "random-access-mode", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validMode(tt.mode); got != tt.want {
				t.Errorf("validMode() = %v, want %v", got, tt.want)
			}
		})
	}
}
