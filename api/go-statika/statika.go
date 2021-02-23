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
	This function creates new Statika GateKeeper client, it's responsible

*/
func New(endpoint string, masterKey string, uploadKey string) (GateKeeper, error) {
	rc := resty.New()

	return GateKeeper{
		restClient: rc,
		endpoint:   endpoint,
		masterKey:  masterKey,
		uploadKey:  uploadKey,
	}, nil
}

// IssueClient This method gets new upload token from Statika server, and returns a client instance with this token
func (gk GateKeeper) IssueClient(bucket string) (Client, error) {
	token, err := gk.IssueUploadToken(bucket)
	if err != nil {
		return Client{}, err
	}

	return Client{
		uploadToken: token,
		bucket:      bucket,
		endpoint:    gk.endpoint,
		resty:       gk.restClient,
	}, nil
}

// NewClient This method creates new client with existing upload token
func (gk GateKeeper) NewClient(bucket string, uploadToken UploadToken) (Client, error) {

	return Client{
		uploadToken: uploadToken,
		bucket:      bucket,
		endpoint:    gk.endpoint,
		resty:       gk.restClient,
	}, nil
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
