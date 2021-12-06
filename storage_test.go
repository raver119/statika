package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/raver119/statika/classes"
	"github.com/stretchr/testify/require"
)

func TestStorage_Universal(t *testing.T) {
	// FIXME: atm this is linux-only. it should become windows-compatible
	require.NoError(t, os.MkdirAll("/tmp/first", 0755))
	require.NoError(t, os.MkdirAll("/tmp/second", 0755))

	bucket := "foo"
	fileName := "file.txt"
	content := uuid.NewString()
	var storages = []Storage{NewLocalStorage("/tmp/first"), NewMirrorStorage(NewLocalStorage("/tmp/second"), NewLocalStorage("/tmp/third"))}

	if spacesBucket, ok := os.LookupEnv("S3_BUCKET"); ok {
		// this part of test is optional
		storage, err := NewSpacesStorage(spacesBucket, "https://nyc3.digitaloceanspaces.com")
		if err == nil {
			storages = append(storages, storage)
		}
	}

	for _, s := range storages {
		storage := s
		t.Run(s.Name(), func(t *testing.T) {
			// put single file there
			_, err := storage.Put(bucket, fileName, strings.NewReader(content))
			require.NoError(t, err)

			// there must be 1 file
			list, err := storage.List(bucket)
			require.NoError(t, err)
			require.Len(t, list, 1)
			require.Equal(t, []classes.FileEntry{{FileName: "file.txt"}}, list)

			// file must be availabe & readable
			reader, err := storage.Get(bucket, fileName)
			require.NoError(t, err)

			// check content
			b, err := ioutil.ReadAll(reader)
			require.NoError(t, err)
			require.Equal(t, content, string(b))

			require.NoError(t, storage.Delete(bucket, fileName))

			// once file's removed, Get must return error
			_, err = storage.Get(bucket, fileName)
			require.Error(t, err)
		})
	}

	// everything must be removed 
	require.NoError(t, os.RemoveAll("/tmp/first"))
	require.NoError(t, os.RemoveAll("/tmp/second"))
}