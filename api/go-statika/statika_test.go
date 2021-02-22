package statika

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	endpoint  = "http://localhost:8080"
	masterKey = "TEST_MASTER_KEY"
	uploadKey = "TEST_UPLOAD_KEY"
)

func TestGateKeeper_IssueUploadToken(t *testing.T) {
	gk, err := New(endpoint, masterKey, uploadKey)
	require.NoError(t, err)

	_, err = gk.IssueUploadToken("test_bucket")
	require.NoError(t, err)
}

func TestGateKeeper_IssueUploadToken_2(t *testing.T) {
	gk, err := New(endpoint, masterKey, "bad key")
	require.NoError(t, err)

	_, err = gk.IssueUploadToken("test_bucket")
	require.Error(t, err)
}
