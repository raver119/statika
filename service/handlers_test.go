package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/resty.v1"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const TEST_B = "my test bucket"
const TEST_P = 8080

/*
	This test validates login + file upload procedure
*/
func TestApiHandler_LoginUpload(t *testing.T) {
	engine, err := CreateEngine(TEST_M, TEST_U, "/tmp", TEST_P)
	require.NoError(t, err)

	err = engine.StartAsync()
	require.NoError(t, err)
	time.Sleep(time.Second)

	client := resty.New()

	negativeAuthReq := UploadAuthenticationRequest{"bad upload key", TEST_B}
	positiveAuthReq := UploadAuthenticationRequest{TEST_U, TEST_B}

	// negative test
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetBody(negativeAuthReq).
		Post("http://localhost:8080/rest/v1/auth/upload")

	require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())

	// positive test
	resp, err = client.R().
		SetHeader("Accept", "application/json").
		SetBody(positiveAuthReq).
		Post("http://localhost:8080/rest/v1/auth/upload")

	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode())

	var authResp AuthenticationResponse
	err = json.Unmarshal(resp.Body(), &authResp)
	require.NoError(t, err)

	fileContent := "file content"

	// now try to upload file
	positiveUploadReq := UploadRequest{
		Filename: "file name.txt",
		Bucket:   TEST_B,
		Meta:     nil,
		Payload:  base64.StdEncoding.EncodeToString([]byte(fileContent)),
	}

	// unauthorized file upload
	resp, err = client.R().
		SetHeader("Accept", "application/json").
		SetBody(positiveUploadReq).
		Post("http://localhost:8080/rest/v1/file")

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())

	// authorized file upload
	resp, err = client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(authResp.Token).
		SetBody(positiveUploadReq).
		Post("http://localhost:8080/rest/v1/file")
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode())

	var uploadResp UploadResponse
	err = json.Unmarshal(resp.Body(), &uploadResp)
	require.NoError(t, err)

	assert.Equal(t, "/my+test+bucket/file+name.txt", uploadResp.FileName)

	// now, it should be possible to request file
	r, err := http.Get(fmt.Sprintf("http://localhost:%v/%v", TEST_P, uploadResp.FileName))
	require.NoError(t, err)

	bytes, _ := ioutil.ReadAll(r.Body)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, fileContent, string(bytes))
	require.NoError(t, engine.Stop())
}

func TestApiHandler_UpdateDelete(t *testing.T) {
	engine, err := CreateEngine(TEST_M, TEST_U, "/tmp", TEST_P)
	require.NoError(t, err)

	err = engine.StartAsync()
	require.NoError(t, err)

	time.Sleep(time.Second)

	client := resty.New()

	positiveAuthReq := UploadAuthenticationRequest{TEST_U, TEST_B}

	ar, err := client.R().
		SetBody(positiveAuthReq).
		Post("http://localhost:8080/rest/v1/auth/upload")
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, ar.StatusCode())

	var authResp AuthenticationResponse
	require.NoError(t, json.Unmarshal(ar.Body(), &authResp))

	originalContent := "my original content"
	updatedContent := "my updated content"

	positiveUploadReq := UploadRequest{
		Filename: "somefile.txt",
		Bucket:   TEST_B,
		Meta:     nil,
		Payload:  base64.StdEncoding.EncodeToString([]byte(originalContent)),
	}

	ur, err := client.R().
		SetAuthToken(authResp.Token).
		SetBody(positiveUploadReq).
		Post("http://localhost:8080/rest/v1/file")
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, ur.StatusCode())

	var uploadResp UploadResponse
	require.NoError(t, json.Unmarshal(ur.Body(), &uploadResp))

	// test update
	positiveUploadReq = UploadRequest{
		Filename: "somefile.txt",
		Bucket:   TEST_B,
		Meta:     nil,
		Payload:  base64.StdEncoding.EncodeToString([]byte(updatedContent)),
	}
	ur, err = client.R().
		SetAuthToken(authResp.Token).
		SetBody(positiveUploadReq).
		Post("http://localhost:8080/rest/v1/file")
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
		Delete(fmt.Sprintf("http://localhost:8080%v", uploadResp.FileName))
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, dr.StatusCode())

	// test positive delete
	dr, err = client.R().
		SetAuthToken(authResp.Token).
		Delete(fmt.Sprintf("http://localhost:8080%v", uploadResp.FileName))
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, dr.StatusCode())

	// file must be absent
	r, err = http.Get(fmt.Sprintf("http://localhost:%v/%v", TEST_P, uploadResp.FileName))
	require.NoError(t, err)

	bytes, _ = ioutil.ReadAll(r.Body)
	assert.Equal(t, http.StatusNotFound, r.StatusCode)

	assert.NoError(t, engine.Stop())
}
