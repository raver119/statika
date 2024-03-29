package classes

import (
	"encoding/json"
	"net/http"
)

const META_EXTENSION = ".statika_metainfo"

type MetaInfo map[string]string

type CloseableReader interface {
	Read(b []byte) (int, error)
	Close() error
}

/*
	API data structures
*/
type MasterAuthenticationRequest struct {
	Token string `json:"token"`
}

type UploadAuthenticationRequest struct {
	Token   string   `json:"token"`   // Auth token. Must match whatever was set in UPLOAD_TOKEN env var
	Bucket  string   `json:"bucket"`  //Deprecated field. Will be removed soon.
	Buckets []string `json:"buckets"` // Target folder for this key. Other buckets will be hidden and unavailable.
}

type AuthenticationResponse struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

func (ar AuthenticationResponse) ToJSON() []byte {
	js, _ := json.Marshal(ar)
	return js
}

type UploadRequest struct {
	Filename string            `json:"filename"`
	Bucket   string            `json:"bucket"`
	Meta     map[string]string `json:"meta"`
	Payload  string            `json:"payload"`
}

type UploadResponse struct {
	Id       string `json:"id"`
	FileName string `json:"filename"`
}

type ListResponse struct {
	Bucket string      `json:"bucket"`
	Files  []FileEntry `json:"files"`
}

func (lr ListResponse) ToJSON() []byte {
	js, _ := json.Marshal(lr)
	return js
}

func (ur UploadResponse) ToJSON() []byte {
	js, _ := json.Marshal(ur)
	return js
}

type FileEntry struct {
	FileName string `json:"filename"`
}

type UpdateMetaRequest struct {
	Id     string   `json:"id"`
	Bucket string   `json:"bucket"`
	Meta   MetaInfo `json:"meta"`
}

type ApiResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (ar ApiResponse) ToJSON() []byte {
	js, _ := json.Marshal(ar)
	return js
}

func ResponseOK() []byte {
	return ApiResponse{StatusCode: http.StatusOK}.ToJSON()
}
