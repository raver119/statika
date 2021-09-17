package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/raver119/statika/classes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalStorage_Get(t *testing.T) {
	testString := "test string"
	ls := NewLocalStorage("/tmp")
	input := strings.NewReader(testString)
	_, err := ls.Put("", "test.txt", input)
	require.NoError(t, err)

	defer ls.Delete("", "test.txt")

	r, err := ls.Get("", "test.txt")
	assert.NoError(t, err)

	bytes, err := ioutil.ReadAll(r)
	assert.Equal(t, testString, string(bytes))
	assert.NoError(t, r.Close())
}

func TestLocalStorage_Put(t *testing.T) {
	testString := "test string"
	ls := NewLocalStorage("/tmp")
	input := strings.NewReader(testString)
	f, err := ls.Put("", "test.txt", input)
	require.NoError(t, err)

	defer ls.Delete("", "test.txt")

	require.FileExists(t, "/tmp/"+f)

	bytes, err := ioutil.ReadFile("/tmp/" + f)
	require.NoError(t, err)

	assert.Equal(t, testString, string(bytes))
}

func TestLocalStorage_Delete(t *testing.T) {
	ls := NewLocalStorage("/tmp")
	f, err := ls.Put("", "test2.txt", strings.NewReader("some content"))
	require.NoError(t, err)

	err = ls.Delete("", "test2.txt")
	require.NoError(t, err)
	require.NoFileExists(t, "/tmp/"+f)
}

func TestLocalStorage_List(t *testing.T) {
	bucket := uuid.NewString()
	testString := "test string"
	ls := NewLocalStorage("/tmp")

	input := strings.NewReader(testString)
	_, err := ls.Put(bucket, "test1.txt", input)
	require.NoError(t, err)

	_, err = ls.Put(bucket, "test2.txt", input)
	require.NoError(t, err)

	files, err := ls.List(bucket)
	require.NoError(t, err)
	assert.Equal(t, []classes.FileEntry{{FileName: "test1.txt"}, {FileName: "test2.txt"}}, files)

	assert.NoError(t, ls.Delete(bucket, "test1.txt"))
	assert.NoError(t, ls.Delete(bucket, "test2.txt"))
}

func TestLocalStorage_GetMeta(t *testing.T) {
	bucket := uuid.NewString()
	ls := NewLocalStorage("/tmp")

	meta := classes.MetaInfo{
		"alpha": "1",
		"beta":  "2",
	}

	fileName := "random_file.txt"
	require.NoError(t, ls.PutMeta(bucket, fileName, meta))

	restored, err := ls.GetMeta(bucket, fileName)
	require.NoError(t, err)
	assert.Equal(t, meta, restored)

	require.NoError(t, ls.DeleteMeta(bucket, fileName))

	restored, err = ls.GetMeta(bucket, fileName)
	require.NoError(t, err)
	require.Equal(t, classes.MetaInfo{}, restored)
}

func TestLocalStorage_NestedFolders(t *testing.T) {
	bucket := uuid.NewString()
	ls := NewLocalStorage("/tmp")

	fileNames := []string{"nested/file.txt", "nested/deeper/file.txt"}

	meta := classes.MetaInfo{
		"alpha": "1",
		"beta":  "2",
	}

	for _, fileName := range fileNames {
		fname, err := ls.Put(bucket, fileName, strings.NewReader("test"))
		log.Printf("Processing uploaded file: %v", fname)
		require.NoError(t, err)
		require.FileExists(t, fmt.Sprintf("/tmp/%v", fname))
		require.NoError(t, ls.PutMeta(bucket, fileName, meta))

		reader, err := ls.Get(bucket, fileName)
		require.NoError(t, err)

		restored, err := ls.GetMeta(bucket, fileName)
		require.NoError(t, err)
		assert.Equal(t, meta, restored)

		body, err := ioutil.ReadAll(reader)
		require.NoError(t, err)
		require.Equal(t, "test", string(body))

		require.NoError(t, ls.Delete(bucket, fileName))
		require.NoFileExists(t, fmt.Sprintf("/tmp/%v", fname))

		// once file removed Meta should be gone as well
		restored, err = ls.GetMeta(bucket, fileName)
		require.NoError(t, err)
		require.Equal(t, classes.MetaInfo{}, restored)
	}
}
