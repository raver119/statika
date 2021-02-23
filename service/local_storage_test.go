package main

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"strings"
	"testing"
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
	bucket := "random_bucket_name"
	testString := "test string"
	ls := NewLocalStorage("/tmp")

	input := strings.NewReader(testString)
	_, err := ls.Put(bucket, "test1.txt", input)
	require.NoError(t, err)

	_, err = ls.Put(bucket, "test2.txt", input)
	require.NoError(t, err)

	files, err := ls.List(bucket)
	require.NoError(t, err)
	assert.Equal(t, []FileEntry{{FileName: "test1.txt"}, {FileName: "test2.txt"}}, files)

	assert.NoError(t, ls.Delete(bucket, "test1.txt"))
	assert.NoError(t, ls.Delete(bucket, "test2.txt"))
}

func TestLocalStorage_GetMeta(t *testing.T) {
	bucket := uuid.New().String()
	//testString := "test string"
	ls := NewLocalStorage("/tmp")

	meta := MetaInfo{
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
	require.Error(t, err)
}
