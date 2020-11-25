package main

import (
	"encoding/json"
	"io"
)

type Storage interface {
	// must be a stream instead
	Put(bucket string, name string, r io.Reader) (string, error)
	Get(bucket string, name string) (CloseableReader, error)
	Delete(bucket string, name string) error
}

type CloseableReader interface {
	 Read(b []byte) (int, error)
	 Close() error
}


/*
	API data structures
 */
type MasterAuthenticationRequest struct {
	Token		string		`json:"token"`
}

type UploadAuthenticationRequest struct {
	Token		string		`json:"token"`	// Auth token. Must match whatever was set in UPLOAD_TOKEN env var
	Bucket		string		`json:"bucket"`	// Target folder for this key. Other buckets will be hidden and unavailable.
}

type AuthenticationResponse struct {
	Token		string		`json:"token"`
	Expires		int64		`json:"expires"`
}

func (ar AuthenticationResponse) ToJSON() []byte {
	js, _ := json.Marshal(ar)
	return js
}


type UploadRequest struct {
	Filename	string				`json:"filename"`
	Bucket 		string				`json:"bucket"`
	Meta		map[string]string	`json:"meta"`
	Payload		string				`json:"payload"`
}

type UploadResponse struct {
	Id			string				`json:"id"`
	FileName	string				`json:"file_name"`
}

func (ur UploadResponse) ToJSON() []byte {
	js, _ := json.Marshal(ur)
	return js
}

type UpdateMetaRequest struct {
	Id			string 				`json:"id"`
	Bucket 		string				`json:"bucket"`
	Meta		map[string]string	`json:"meta"`
}
