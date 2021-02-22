package statika

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestClient_UploadFile(t *testing.T) {
	testBucket := uuid.New().String()
	gk, err := New(endpoint, masterKey, uploadKey)
	require.NoError(t, err)

	client, err := gk.BuildClient(testBucket)
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
	testBucket := uuid.New().String()
	gk, err := New(endpoint, masterKey, uploadKey)
	require.NoError(t, err)

	client, err := gk.BuildClient(testBucket)
	require.NoError(t, err)

	content := "random text content"

	_, err = client.UploadFile("file5.txt", strings.NewReader(content))
	require.NoError(t, err)

	_, err = client.UploadFile("file6.txt", strings.NewReader(content))
	require.NoError(t, err)

	files, err := client.ListFiles()
	require.NoError(t, err)

	// i.e. http://localhost:9191/bucket_name/file5.txt
	p := fmt.Sprintf("http://localhost:9191/%v/", testBucket)
	require.Equal(t, []FileEntry{{p + "file5.txt"}, {p + "file6.txt"}}, files)
}
