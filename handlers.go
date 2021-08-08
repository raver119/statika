package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	. "github.com/raver119/statika/classes"
	. "github.com/raver119/statika/utils"
	. "github.com/raver119/statika/wt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ApiHandler struct {
	masterKey string
	userKey   string

	storage   *Storage
	tokenizer Tokenizer
}

func NewApiHandler(masterKey string, userKey string, storage *Storage) (*ApiHandler, error) {
	return &ApiHandler{masterKey: masterKey, userKey: userKey, storage: storage, tokenizer: NewTokenizer()}, nil
}

func (srv *ApiHandler) validateUploadToken(r *http.Request, bucket string) (ok bool) {
	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		token = r.FormValue("token")
		if len(token) == 0 {
			return false
		}
	}

	ok, err := srv.tokenizer.ValidateUploadToken(token, bucket)
	if ok && err == nil {
		return ok
	} else {
		return false
	}
}

/*
	Login endpoint
*/
func (srv *ApiHandler) LoginUpload(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if !OptionallyReport("unable to read message body", w, err) {
		return
	}

	// deserialize request
	var req UploadAuthenticationRequest
	err = json.Unmarshal(body, &req)
	if !OptionallyReport("unable to deserialize UploadAuthenticationRequest", w, err) {
		return
	}

	// make sure this is authorized request
	if srv.userKey != req.Token {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// create upload token
	token, err := srv.tokenizer.CreateUploadToken(req)
	if !OptionallyReport("unable to create upload token", w, err) {
		return
	}

	response := AuthenticationResponse{
		Token:   token,
		Expires: time.Now().Unix() + 86400*365,
	}

	// and send it back
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response.ToJSON())
}

func (srv *ApiHandler) LoginMaster(w http.ResponseWriter, r *http.Request) {

}

/*
	This method retrieves Meta information
*/
func (srv *ApiHandler) GetMeta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	fileName := vars["fileName"]

	// TODO: check meta file existence first, 404 would be more informative
	meta, err := (*srv.storage).GetMeta(bucket, fileName)
	if !OptionallyReport("failed to get meta", w, err) {
		return
	}

	b, _ := json.Marshal(meta)
	_, _ = w.Write(b)
}

/*
	This method updates Meta information
*/
func (srv *ApiHandler) SetMeta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	fileName := vars["fileName"]

	body, err := ioutil.ReadAll(r.Body)
	if !OptionallyReport("failed to read body", w, err) {
		return
	}

	var meta MetaInfo
	err = json.Unmarshal(body, &meta)
	if !OptionallyReport("failed to deserialize meta", w, err) {
		return
	}

	err = (*srv.storage).PutMeta(bucket, fileName, meta)
	if !OptionallyReport("failed to store meta", w, err) {
		return
	}

	_, _ = w.Write(ResponseOK())
}

/*
	This method removes Meta information
*/
func (srv *ApiHandler) DeleteMeta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	fileName := vars["fileName"]

	err := (*srv.storage).DeleteMeta(bucket, fileName)
	if !OptionallyReport("failed to delete meta", w, err) {
		return
	}

	_, _ = w.Write(ResponseOK())
}

func (srv *ApiHandler) Ping(w http.ResponseWriter, r *http.Request) {
	// validate authorization
	if !srv.validateUploadToken(r, "") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(ResponseOK())
}

func (srv *ApiHandler) processUpload(bucket string, fileName string, reader io.ReadSeeker, w http.ResponseWriter) (ur UploadResponse, err error) {
	_, err = (*srv.storage).Put(bucket, fileName, reader)
	if !OptionallyReport("put failed", w, err) {
		return
	}

	// TODO: use ID wisely
	response := UploadResponse{Id: uuid.New().String(), FileName: fmt.Sprintf("/%v/%v", bucket, fileName)}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response.ToJSON())
	return
}

// C
func (srv *ApiHandler) Upload(w http.ResponseWriter, r *http.Request) {
	var ct = r.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "multipart/form-data;") {
		bucket := r.FormValue("bucket")
		if len(bucket) == 0 {
			http.Error(w, "Missing bucket name", http.StatusBadRequest)
			return
		}

		token := r.FormValue("token")
		if len(token) == 0 {
			http.Error(w, "Missing token name", http.StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if !OptionallyReport("failed to fetch file from the form", w, err) {
			return
		}

		if !srv.validateUploadToken(r, bucket) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		_, _ = srv.processUpload(bucket, fileHeader.Filename, file, w)
	} else {
		// read request
		body, err := ioutil.ReadAll(r.Body)
		if !OptionallyReport("failed to read body", w, err) {
			return
		}

		// deserialize request
		var req UploadRequest
		err = json.Unmarshal(body, &req)
		if !OptionallyReport("failed to deserialized UploadRequest", w, err) {
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
		if !OptionallyReport("unable to decode payload", w, err) {
			return
		}

		_, _ = srv.processUpload(req.Bucket, req.Filename, bytes.NewReader(b), w)
	}
}

// TEST_U
func (srv *ApiHandler) Update(w http.ResponseWriter, r *http.Request) {

}

// D
func (srv *ApiHandler) Delete(w http.ResponseWriter, r *http.Request) {

}

func (srv *ApiHandler) List(w http.ResponseWriter, r *http.Request) {
	bucket := mux.Vars(r)["bucket"]
	if !srv.validateUploadToken(r, bucket) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	entries, err := (*srv.storage).List(bucket)
	if !OptionallyReport("unable to list bucket", w, err) {
		return
	}

	// TODO: decide if the full path really preferred here
	//for i, v := range entries {
	//	entries[i].FileName = fmt.Sprintf("/%v/%v", bucket, v.FileName)
	//}

	var response = ListResponse{
		Bucket: bucket,
		Files:  entries,
	}

	_, _ = w.Write(response.ToJSON())
}
