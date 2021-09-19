package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/raver119/statika/api"
	"github.com/raver119/statika/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	// server-timing header field must be present
	assert.Contains(t, response.Header.Get("Server-Timing"), "Storage READ")

	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, content, string(body))

	require.NoError(t, client.DeleteFile(ur.FileName))

	// file must be absent now
	response, err = http.Get(endpoint + ur.FileName)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	// let's ping the server, to make sure our key is still ok
	require.NoError(t, client.Ping())
}

func TestClient_ListFiles(t *testing.T) {
	bucket := utils.GetEnvOrPanic("S3_BUCKET")
	ep := utils.GetEnvOrDefault("S3_ENDPOINT", "https://nyc3.digitaloceanspaces.com")
	_ = utils.GetEnvOrPanic("S3_KEY")
	_ = utils.GetEnvOrPanic("S3_SECRET")
	s3storage, err := NewSpacesStorage(bucket, ep)
	require.NoError(t, err)

	var storages = []Storage{NewLocalStorage("/tmp"), s3storage}

	for _, ls := range storages {
		t.Run(ls.Name(), func(t *testing.T) {
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

			_, err = client.UploadFile("nested/file6.txt", strings.NewReader(content))
			require.NoError(t, err)

			_, err = client.UploadFile("/nested/sub/file7.txt", strings.NewReader(content))
			require.NoError(t, err)

			files, err := client.ListFiles()
			require.NoError(t, err)

			require.Equal(t, []api.FileEntry{{FileName: "file5.txt"}, {FileName: "nested/file6.txt"}, {FileName: "nested/sub/file7.txt"}}, files)

			// ping must succeed
			require.NoError(t, client.Ping())

			// Get must succeed as well
			for _, e := range files {
				r, err := client.Get(e.FileName)
				require.NoError(t, err)

				require.Equal(t, content, string(r))
			}
		})
	}
}

func TestClient_Meta(t *testing.T) {
	bucket := utils.GetEnvOrPanic("S3_BUCKET")
	ep := utils.GetEnvOrDefault("S3_ENDPOINT", "https://nyc3.digitaloceanspaces.com")
	_ = utils.GetEnvOrPanic("S3_KEY")
	_ = utils.GetEnvOrPanic("S3_SECRET")
	s3storage, err := NewSpacesStorage(bucket, ep)
	require.NoError(t, err)

	var storages = []Storage{NewLocalStorage("/tmp"), s3storage}

	for _, ls := range storages {
		t.Run(ls.Name(), func(t *testing.T) {
			e, err := CreateEngine(masterKey, uploadKey, &ls, 9191)
			require.NoError(t, err)
			require.NoError(t, e.StartAsync())
			time.Sleep(1 * time.Second)
			defer func() {
				_ = e.Stop()
			}()

			testBucket := uuid.NewString()
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
			assert.Equal(t, "2", restored["beta"])
			assert.Equal(t, 2, len(restored))

			require.NoError(t, client.DeleteMeta("file5.txt"))

			// must be empty map
			restored, err = client.GetMeta("file5.txt")
			require.NoError(t, err)
			assert.Equal(t, api.MetaInfo{}, restored)

			// ping must succeed
			require.NoError(t, client.Ping())

			// now master a fake client, so it fails on ping
			fake, _ := gk.NewClient("bad bucket", "bad token")
			require.Error(t, fake.Ping())
		})
	}
}
