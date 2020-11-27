package statika

import "testing"

const (
	endpoint  = "http://localhost:8080"
	masterKey = "TEST_MASTER_KEY"
	uploadKey = "TEST_UPLOAD_KEY"
)

func TestGateKeeper_IssueUploadToken(t *testing.T) {
	gk, err := New(endpoint, masterKey, uploadKey)
	if err != nil {
		t.Fatal(err)
	}

	_, err = gk.IssueUploadToken("test_bucket")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGateKeeper_IssueUploadToken_2(t *testing.T) {
	gk, err := New(endpoint, masterKey, "bad key")
	if err != nil {
		t.Fatal(err)
	}

	_, err = gk.IssueUploadToken("test_bucket")
	if err == nil {
		t.Fatalf("Token shouldn't be issued")
	}
}
