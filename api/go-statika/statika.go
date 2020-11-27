package statika

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

/*
	This "class" provides methods for tokens management
*/
type GateKeeper struct {
	restClient *resty.Client
	endpoint   string

	masterKey string
	uploadKey string
}

/*
	This function creates new Statika client

*/
func New(endpoint string, masterKey string, uploadKey string) (c *GateKeeper, err error) {
	rc := resty.New()

	c = &GateKeeper{
		restClient: rc,
		endpoint:   endpoint,
		masterKey:  masterKey,
		uploadKey:  uploadKey,
	}
	return
}

func (gk GateKeeper) IssueUploadToken(bucket string) (token UploadToken, err error) {
	endpoint := fmt.Sprintf("%v/rest/v1/auth/upload", gk.endpoint)
	client := resty.New()

	upReq := UploadAuthenticationRequest{
		Token:  gk.uploadKey,
		Bucket: bucket,
	}

	resp, err := client.R().
		SetBody(upReq).
		Post(endpoint)

	if err != nil {
		return
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		err = fmt.Errorf("statika authentication failed: bad token")
		return
	}

	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("statika request failed with statusCode %v", resp.StatusCode())
		return
	}

	var response AuthenticationResponse
	err = json.Unmarshal(resp.Body(), &response)
	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("statika authentication failed: JSON versions mismatch")
		return
	}

	token = UploadToken(response.Token)
	return
}
