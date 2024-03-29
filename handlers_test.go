package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/raver119/statika/classes"
	"github.com/raver119/statika/wt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TEST_B = "my_test_bucket"
const TEST_P = 9191

func TestApiHandler_Login(t *testing.T) {
	token, err := wt.DevTokenizer("alpha").CreateUploadToken(classes.UploadAuthenticationRequest{
		Buckets: []string{"css"},
	})

	require.NoError(t, err)
	log.Printf("Token: %v", token)
}

/*
	This test validates login + file upload procedure
*/
func TestApiHandler_LoginUpload(t *testing.T) {
	var ls Storage = NewLocalStorage("/tmp")
	engine, err := CreateEngine(TEST_M, TEST_U, &ls, TEST_P)
	require.NoError(t, err)

	err = engine.StartAsync()
	require.NoError(t, err)
	time.Sleep(time.Second)

	client := resty.New()

	negativeAuthReq := classes.UploadAuthenticationRequest{Token: "bad upload key", Bucket: TEST_B}
	positiveAuthReq := classes.UploadAuthenticationRequest{Token: TEST_U, Bucket: TEST_B}

	// negative test
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetBody(negativeAuthReq).
		Post("http://localhost:9191/rest/v1/auth/upload")

	require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())

	// positive test
	resp, err = client.R().
		SetHeader("Accept", "application/json").
		SetBody(positiveAuthReq).
		Post("http://localhost:9191/rest/v1/auth/upload")

	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode())

	var authResp classes.AuthenticationResponse
	err = json.Unmarshal(resp.Body(), &authResp)
	require.NoError(t, err)

	fileContent := "file content"

	// now try to upload file
	positiveUploadReq := classes.UploadRequest{
		Filename: "file name.txt",
		Bucket:   TEST_B,
		Meta:     nil,
		Payload:  base64.StdEncoding.EncodeToString([]byte(fileContent)),
	}

	// unauthorized file upload
	resp, err = client.R().
		SetHeader("Accept", "application/json").
		SetBody(positiveUploadReq).
		Post("http://localhost:9191/rest/v1/file")

	require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())

	// authorized file upload
	resp, err = client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(authResp.Token).
		SetBody(positiveUploadReq).
		Post("http://localhost:9191/rest/v1/file")
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode())

	var uploadResp classes.UploadResponse
	err = json.Unmarshal(resp.Body(), &uploadResp)
	require.NoError(t, err)

	assert.Equal(t, "/my_test_bucket/file name.txt", uploadResp.FileName)

	// now, it should be possible to request file
	r, err := http.Get(fmt.Sprintf("http://localhost:%v/%v", TEST_P, uploadResp.FileName))
	require.NoError(t, err)

	bytes, _ := ioutil.ReadAll(r.Body)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, fileContent, string(bytes))
	require.NoError(t, engine.Stop())
}

func TestApiHandler_UpdateDelete(t *testing.T) {
	var ls Storage = NewLocalStorage("/tmp")
	engine, err := CreateEngine(TEST_M, TEST_U, &ls, TEST_P)
	require.NoError(t, err)

	err = engine.StartAsync()
	require.NoError(t, err)

	time.Sleep(time.Second)

	client := resty.New()

	positiveAuthReq := classes.UploadAuthenticationRequest{Token: TEST_U, Bucket: TEST_B}

	ar, err := client.R().
		SetBody(positiveAuthReq).
		Post("http://localhost:9191/rest/v1/auth/upload")
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, ar.StatusCode())

	var authResp classes.AuthenticationResponse
	require.NoError(t, json.Unmarshal(ar.Body(), &authResp))

	originalContent := "my original content"
	updatedContent := "my updated content"

	positiveUploadReq := classes.UploadRequest{
		Filename: "somefile.txt",
		Bucket:   TEST_B,
		Meta:     nil,
		Payload:  base64.StdEncoding.EncodeToString([]byte(originalContent)),
	}

	ur, err := client.R().
		SetAuthToken(authResp.Token).
		SetBody(positiveUploadReq).
		Post("http://localhost:9191/rest/v1/file")
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, ur.StatusCode())

	var uploadResp classes.UploadResponse
	require.NoError(t, json.Unmarshal(ur.Body(), &uploadResp))

	// test update
	positiveUploadReq = classes.UploadRequest{
		Filename: "somefile.txt",
		Bucket:   TEST_B,
		Meta:     nil,
		Payload:  base64.StdEncoding.EncodeToString([]byte(updatedContent)),
	}
	ur, err = client.R().
		SetAuthToken(authResp.Token).
		SetBody(positiveUploadReq).
		Post("http://localhost:9191/rest/v1/file")
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, ur.StatusCode())

	// now make sure content is updated
	r, err := http.Get(fmt.Sprintf("http://localhost:%v/%v", TEST_P, uploadResp.FileName))
	require.NoError(t, err)

	bytes, _ := ioutil.ReadAll(r.Body)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, updatedContent, string(bytes))

	// test negative delete
	dr, err := client.R().
		Delete(fmt.Sprintf("http://localhost:9191%v", uploadResp.FileName))
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, dr.StatusCode())

	// test positive delete
	dr, err = client.R().
		SetAuthToken(authResp.Token).
		Delete(fmt.Sprintf("http://localhost:9191%v", uploadResp.FileName))
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, dr.StatusCode())

	// file must be absent
	r, err = http.Get(fmt.Sprintf("http://localhost:%v/%v", TEST_P, uploadResp.FileName))
	require.NoError(t, err)

	bytes, _ = ioutil.ReadAll(r.Body)
	assert.Equal(t, http.StatusNotFound, r.StatusCode)

	assert.NoError(t, engine.Stop())
}

func TestApiHandler_FormUpload(t *testing.T) {
	var ls Storage = NewLocalStorage("/tmp")
	engine, err := CreateEngine(TEST_M, TEST_U, &ls, TEST_P)
	require.NoError(t, err)

	err = engine.StartAsync()
	require.NoError(t, err)

	time.Sleep(time.Second)

	client := resty.New()
	positiveAuthReq := classes.UploadAuthenticationRequest{Token: TEST_U, Bucket: TEST_B}

	ar, err := client.R().
		SetBody(positiveAuthReq).
		Post("http://localhost:9191/rest/v1/auth/upload")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, ar.StatusCode())

	var authResp classes.AuthenticationResponse
	require.NoError(t, json.Unmarshal(ar.Body(), &authResp))

	var content string = "another_file content"

	ur, err := client.R().
		SetFileReader("file", "another_file.txt", bytes.NewReader([]byte(content))).
		SetFormData(map[string]string{
			"token":  authResp.Token,
			"bucket": TEST_B,
		}).
		Post("http://localhost:9191/rest/v1/file")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, ur.StatusCode(), string(ur.Body()))

	var uploadResp classes.UploadResponse
	require.NoError(t, json.Unmarshal(ur.Body(), &uploadResp))

	// test positive delete
	dr, err := client.R().
		SetAuthToken(authResp.Token).
		Delete(fmt.Sprintf("http://localhost:9191%v", uploadResp.FileName))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, dr.StatusCode())

	assert.NoError(t, engine.Stop())
}

func TestApiHandler_List(t *testing.T) {
	randomBucket := uuid.New().String()

	var ls Storage = NewLocalStorage("/tmp")
	engine, err := CreateEngine(TEST_M, TEST_U, &ls, TEST_P)
	require.NoError(t, err)

	err = engine.StartAsync()
	require.NoError(t, err)

	time.Sleep(time.Second)

	client := resty.New()
	positiveAuthReq := classes.UploadAuthenticationRequest{Token: TEST_U, Bucket: randomBucket}

	ar, err := client.R().
		SetBody(positiveAuthReq).
		Post("http://localhost:9191/rest/v1/auth/upload")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, ar.StatusCode())

	var authResp classes.AuthenticationResponse
	require.NoError(t, json.Unmarshal(ar.Body(), &authResp))

	ur, err := client.R().
		SetFileReader("file", "file1.txt", strings.NewReader("pew-pew-zomg")).
		SetFormData(map[string]string{
			"token":  authResp.Token,
			"bucket": randomBucket,
		}).
		Post("http://localhost:9191/rest/v1/file")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, ur.StatusCode(), string(ur.Body()))

	ur, err = client.R().
		SetFileReader("file", "file2.txt", strings.NewReader("pew-pew-zomg")).
		SetFormData(map[string]string{
			"token":  authResp.Token,
			"bucket": randomBucket,
		}).
		Post("http://localhost:9191/rest/v1/file")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, ur.StatusCode(), string(ur.Body()))

	lr, err := client.R().
		SetAuthToken(authResp.Token).
		Get(fmt.Sprintf("http://localhost:9191/rest/v1/files/%v", randomBucket))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, lr.StatusCode(), string(lr.Body()))

	var listResp classes.ListResponse
	require.NoError(t, json.Unmarshal(lr.Body(), &listResp))
	require.Equal(t, randomBucket, listResp.Bucket)
	require.Equal(t, []classes.FileEntry{{FileName: "file1.txt"}, {FileName: "file2.txt"}}, listResp.Files)
}

func TestApiHandler_Meta(t *testing.T) {
	randomBucket := uuid.New().String()

	var ls Storage = NewLocalStorage("/tmp")
	engine, err := CreateEngine(TEST_M, TEST_U, &ls, TEST_P)
	require.NoError(t, err)

	err = engine.StartAsync()
	require.NoError(t, err)

	time.Sleep(time.Second)

	client := resty.New()
	positiveAuthReq := classes.UploadAuthenticationRequest{Token: TEST_U, Bucket: randomBucket}

	ar, err := client.R().
		SetBody(positiveAuthReq).
		Post("http://localhost:9191/rest/v1/auth/upload")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, ar.StatusCode())

	var authResp classes.AuthenticationResponse
	require.NoError(t, json.Unmarshal(ar.Body(), &authResp))

	meta := classes.MetaInfo{
		"alpha": "1",
		"beta":  "2",
	}

	fileName := "some_file.txt"

	pm, err := client.R().
		SetAuthToken(authResp.Token).
		SetBody(meta).
		Post(fmt.Sprintf("http://localhost:9191/rest/v1/meta/%v/%v", randomBucket, fileName))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, pm.StatusCode())

	gm, err := client.R().
		SetAuthToken(authResp.Token).
		SetBody(meta).
		Get(fmt.Sprintf("http://localhost:9191/rest/v1/meta/%v/%v", randomBucket, fileName))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, gm.StatusCode())

	var restored classes.MetaInfo
	require.NoError(t, json.Unmarshal(gm.Body(), &restored))
	assert.Equal(t, meta, restored)

	dm, err := client.R().
		SetAuthToken(authResp.Token).
		SetBody(meta).
		Delete(fmt.Sprintf("http://localhost:9191/rest/v1/meta/%v/%v", randomBucket, fileName))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, dm.StatusCode())
}

func TestApiHandler_Decode(t *testing.T) {
	decoded, err := url.QueryUnescape("cmF2ZXIxMTk%3D")
	require.NoError(t, err)
	require.Equal(t, "cmF2ZXIxMTk=", decoded)
}
