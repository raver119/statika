package api

type MetaInfo map[string]string

type UploadAuthenticationRequest struct {
	Token   string   `json:"token"`   // Auth token. Must match whatever was set in UPLOAD_TOKEN env var
	Buckets []string `json:"buckets"` // Target folder for this key. Other bucket will be hidden and unavailable.
}

type AuthenticationResponse struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

type UploadToken string

type UploadResponse struct {
	Id       string `json:"id"`
	FileName string `json:"filename"`
}

type FileEntry struct {
	FileName string `json:"filename"`
}

type ListResponse struct {
	Bucket string      `json:"bucket"`
	Files  []FileEntry `json:"files"`
}
