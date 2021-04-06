package main

import (
	"github.com/google/uuid"
	"github.com/raver119/statika/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestClient_UploadFile(t *testing.T) {
	var ls Storage = NewLocalStorage("/tmp")
	e, err := CreateEngine(masterKey, uploadKey, &ls, 9191)
	require.NoError(t, err)
	require.NoError(t, e.StartAsync())
	time.Sleep(1 * time.Second)
	defer func() {
		_ = e.Stop()
	}()

	testBucket := uuid.New().String()
	gk, err := api.NewGateKeeper(endpoint, masterKey, uploadKey)
	require.NoError(t, err)

	client, err := gk.IssueClient(testBucket)
	require.NoError(t, err)

	content := "random text content"

	ur, err := client.UploadFile("file.txt", strings.NewReader(content))
	require.NoError(t, err)

	assert.Equal(t, "/"+testBucket+"/file.txt", ur.FileName)

	// check file existence and content
	response, err := http.Get(endpoint + ur.FileName)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, content, string(body))

	require.NoError(t, client.DeleteFile(ur.FileName))

	// file must be absent now
	response, err = http.Get(endpoint + ur.FileName)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestClient_ListFiles(t *testing.T) {
	var ls Storage = NewLocalStorage("/tmp")
	e, err := CreateEngine(masterKey, uploadKey, &ls, 9191)
	require.NoError(t, err)
	require.NoError(t, e.StartAsync())
	time.Sleep(1 * time.Second)
	defer func() {
		_ = e.Stop()
	}()

	testBucket := uuid.New().String()
	gk, err := api.NewGateKeeper(endpoint, masterKey, uploadKey)
	require.NoError(t, err)

	token, err := gk.IssueUploadToken(testBucket)
	require.NoError(t, err)

	client, err := gk.NewClient(testBucket, token)
	require.NoError(t, err)

	content := "random text content"

	_, err = client.UploadFile("file5.txt", strings.NewReader(content))
	require.NoError(t, err)

	_, err = client.UploadFile("file6.txt", strings.NewReader(content))
	require.NoError(t, err)

	files, err := client.ListFiles()
	require.NoError(t, err)

	require.Equal(t, []api.FileEntry{{"file5.txt"}, {"file6.txt"}}, files)
}

func TestClient_Meta(t *testing.T) {
	var ls Storage = NewLocalStorage("/tmp")
	e, err := CreateEngine(masterKey, uploadKey, &ls, 9191)
	require.NoError(t, err)
	require.NoError(t, e.StartAsync())
	time.Sleep(1 * time.Second)
	defer func() {
		_ = e.Stop()
	}()

	testBucket := uuid.New().String()
	gk, err := api.NewGateKeeper(endpoint, masterKey, uploadKey)
	require.NoError(t, err)

	client, err := gk.IssueClient(testBucket)
	require.NoError(t, err)

	meta := api.MetaInfo{
		"alpha": "1",
		"beta":  "2",
	}

	require.NoError(t, client.SetMeta("file5.txt", meta))

	restored, err := client.GetMeta("file5.txt")
	require.NoError(t, err)
	assert.Equal(t, meta, restored)

	require.NoError(t, client.DeleteMeta("file5.txt"))

	// must be empty map
	restored, err = client.GetMeta("file5.txt")
	require.NoError(t, err)
	assert.Equal(t, api.MetaInfo{}, restored)
}
