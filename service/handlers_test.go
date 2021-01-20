package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	if err != nil {
		t.Fatal(err)
	}

	err = engine.StartAsync()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)

	client := resty.New()

	negativeAuthReq := UploadAuthenticationRequest{"bad upload key", TEST_B}
	positiveAuthReq := UploadAuthenticationRequest{TEST_U, TEST_B}

	// negative test
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetBody(negativeAuthReq).
		Post("http://localhost:8080/rest/v1/auth/upload")

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode() != http.StatusUnauthorized {
		t.Fatalf("unexpected status code: %v", resp.StatusCode())
	}


	// positive test
	resp, err = client.R().
					SetHeader("Accept", "application/json").
					SetBody(positiveAuthReq).
					Post("http://localhost:8080/rest/v1/auth/upload")

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode() != http.StatusOK {
		t.Fatalf("unexpected status code: %v", resp.StatusCode())
	}

	var authResp AuthenticationResponse
	err = json.Unmarshal(resp.Body(), &authResp)
	if err != nil {
		t.Fatal(err)
	}

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

	if resp.StatusCode() != http.StatusUnauthorized {
		t.Fatalf("unexpected status code: %v", resp.StatusCode())
	}

	// authorized file upload
	resp, err = client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(authResp.Token).
		SetBody(positiveUploadReq).
		Post("http://localhost:8080/rest/v1/file")

	if resp.StatusCode() != http.StatusOK {
		t.Fatalf("unexpected status code: %v; body: %v", resp.StatusCode(), string(resp.Body()))
	}

	var uploadResp UploadResponse
	err = json.Unmarshal(resp.Body(), &uploadResp)
	if err != nil {
		t.Fatal(err)
	}

	if uploadResp.FileName != "/my+test+bucket/file+name.txt" {
		t.Fatalf("unexpected file name after upload: [%v]", uploadResp.FileName)
	}

	// now, it should be possible to request file
	r, err := http.Get(fmt.Sprintf("http://localhost:%v/%v", TEST_P, uploadResp.FileName))
	if err != nil {
		t.Fatal(err)
	}

	bytes, _ := ioutil.ReadAll(r.Body)
	if r.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %v; Body: [%v]", r.StatusCode, string(bytes))
	}

	if fileContent != string(bytes) {
		t.Fatalf("response file content doesn't match expectations: [%v]\n", string(bytes))
	}

	err = engine.Stop()
	if err != nil {
		t.Fatal(err)
	}
}

func TestApiHandler_UpdateDelete(t *testing.T) {
	engine, err := CreateEngine(TEST_M, TEST_U, "/tmp", TEST_P)
	if err != nil {
		t.Fatal(err)
	}

	err = engine.StartAsync()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)

	client := resty.New()

	positiveAuthReq := UploadAuthenticationRequest{TEST_U, TEST_B}

	ar, err := client.R().
						SetBody(positiveAuthReq).
						Post("http://localhost:8080/rest/v1/auth/upload")
	if err != nil {
		t.Fatal(err)
	}

	if ar.StatusCode() != http.StatusOK {
		t.Fatalf("unexpected status code: %v", ar.StatusCode())
	}

	var authResp AuthenticationResponse
	_ = json.Unmarshal(ar.Body(), &authResp)


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
	if err != nil {
		t.Fatal(err)
	}

	if ur.StatusCode() != http.StatusOK {
		t.Fatalf("unexpected status code: %v", ur.StatusCode())
	}

	var uploadResp UploadResponse
	err = json.Unmarshal(ur.Body(), &uploadResp)
	if err != nil {
		t.Fatal(err)
	}

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
	if err != nil {
		t.Fatal(err)
	}

	if ur.StatusCode() != http.StatusOK {
		t.Fatalf("unexpected status code: %v", ur.StatusCode())
	}
	// now make sure content is updated
	r, err := http.Get(fmt.Sprintf("http://localhost:%v/%v", TEST_P, uploadResp.FileName))
	if err != nil {
		t.Fatal(err)
	}

	bytes, _ := ioutil.ReadAll(r.Body)
	if r.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %v; Body: [%v]", r.StatusCode, string(bytes))
	}

	if updatedContent != string(bytes) {
		t.Fatalf("response file content doesn't match expectations: [%v]\n", string(bytes))
	}

	// test negative delete
	dr, err := client.R().
		Delete(fmt.Sprintf("http://localhost:8080%v", uploadResp.FileName))
	if err != nil {
		t.Fatal(err)
	}
	if dr.StatusCode() != http.StatusUnauthorized {
		t.Fatalf("unexpected status code: %v; Body: [%v]", dr.StatusCode(), string(dr.Body()))
	}

	// test positive delete
	dr, err = client.R().
		SetAuthToken(authResp.Token).
		Delete(fmt.Sprintf("http://localhost:8080%v", uploadResp.FileName))
	if err != nil {
		t.Fatal(err)
	}
	if dr.StatusCode() != http.StatusOK {
		t.Fatalf("unexpected status code: %v; Body: [%v]", dr.StatusCode(), string(dr.Body()))
	}

	// file must be absent
	r, err = http.Get(fmt.Sprintf("http://localhost:%v/%v", TEST_P, uploadResp.FileName))
	if err != nil {
		t.Fatal(err)
	}

	bytes, _ = ioutil.ReadAll(r.Body)
	if r.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected status code: %v; Body: [%v]", r.StatusCode, string(bytes))
	}

	err = engine.Stop()
	if err != nil {
		t.Fatal(err)
	}
}
