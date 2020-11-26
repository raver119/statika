package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ApiHandler struct {
	masterKey string
	userKey   string

	pa      PersistenceAgent
	storage Storage
}

func NewApiHandler(masterKey string, userKey string, rootFolder string) (*ApiHandler, error) {
	mh := GetEnvOrDefault("MEMCACHED_HOST", "localhost")
	pa, err := NewPersistenceAgent(mh, 11211)
	return &ApiHandler{pa: pa, masterKey: masterKey, userKey: userKey, storage: NewLocalStorage(rootFolder)}, err
}

func (srv *ApiHandler) validateUploadToken(r *http.Request, bucket string) (ok bool) {
	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		return false
	}

	if len(bucket) > 0 {
		return srv.pa.CheckUploadToken(token, bucket)
	} else {
		return srv.pa.TouchUploadToken(token)
	}
}

/*
	Login endpoint
*/
func (srv *ApiHandler) LoginUpload(w http.ResponseWriter, r *http.Request) {
	// CORS setup etc
	SetupResponseHeaders(&w, r)
	body, err := ioutil.ReadAll(r.Body)
	if !OptionallyReport(w, err) {
		return
	}

	// deserialize request
	var req UploadAuthenticationRequest
	err = json.Unmarshal(body, &req)
	if !OptionallyReport(w, err) {
		return
	}

	req.Bucket = url.QueryEscape(req.Bucket)

	// make sure this is authorized request
	if srv.userKey != req.Token {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("expected token: [%v]; actual token: [%v]", srv.userKey, req.Token)))
		return
	}

	// create upload token
	token, err := srv.pa.CreateUploadToken(req)
	if !OptionallyReport(w, err) {
		return
	}

	response := AuthenticationResponse{
		Token:   token,
		Expires: time.Now().Unix() + int64(srv.pa.Expiration()-60),
	}

	// and send it back
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response.ToJSON())
}

func (srv *ApiHandler) LoginMaster(w http.ResponseWriter, r *http.Request) {
	SetupResponseHeaders(&w, r)
}

/*
	This method retrieves Meta information
*/
func (srv *ApiHandler) GetMeta(w http.ResponseWriter, r *http.Request) {
	SetupResponseHeaders(&w, r)
}

/*
	This method updates Meta information
*/
func (srv *ApiHandler) UpdateMeta(w http.ResponseWriter, r *http.Request) {
	SetupResponseHeaders(&w, r)
}

func (srv *ApiHandler) Ping(w http.ResponseWriter, r *http.Request) {
	// CORS setup
	SetupResponseHeaders(&w, r)

	// validate authorization
	if !srv.validateUploadToken(r, "") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseOK())
}

// C
func (srv *ApiHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// CORS setup
	SetupResponseHeaders(&w, r)

	// read request
	body, err := ioutil.ReadAll(r.Body)
	if !OptionallyReport(w, err) {
		return
	}

	// deserialize request
	var req UploadRequest
	err = json.Unmarshal(body, &req)
	if !OptionallyReport(w, err) {
		return
	}

	req.Filename = strings.TrimSpace(req.Filename)
	if len(req.Filename) == 0 {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	req.Filename = url.QueryEscape(req.Filename)
	req.Bucket = url.QueryEscape(req.Bucket)

	// validate authorization
	if !srv.validateUploadToken(r, req.Bucket) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// if everything is ok - store the file
	b, err := base64.StdEncoding.DecodeString(req.Payload)
	if !OptionallyReport(w, err) {
		return
	}

	_, err = srv.storage.Put(req.Bucket, req.Filename, bytes.NewReader(b))
	if !OptionallyReport(w, err) {
		return
	}

	// TODO: use ID wisely
	response := UploadResponse{Id: uuid.New().String(), FileName: fmt.Sprintf("/%v/%v", req.Bucket, req.Filename)}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response.ToJSON())
}

// TEST_U
func (srv *ApiHandler) Update(w http.ResponseWriter, r *http.Request) {
	SetupResponseHeaders(&w, r)
}

// D
func (srv *ApiHandler) Delete(w http.ResponseWriter, r *http.Request) {
	SetupResponseHeaders(&w, r)
}
