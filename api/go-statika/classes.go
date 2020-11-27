package statika

type UploadAuthenticationRequest struct {
	Token  string `json:"token"`  // Auth token. Must match whatever was set in UPLOAD_TOKEN env var
	Bucket string `json:"bucket"` // Target folder for this key. Other buckets will be hidden and unavailable.
}

type AuthenticationResponse struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

type UploadToken string
