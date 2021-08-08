package api

import (
	"encoding/json"
	"fmt"
	"github.com/cristalhq/jwt/v3"
	"github.com/go-resty/resty/v2"
	. "github.com/raver119/statika/wt"
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
func NewGateKeeper(endpoint string, masterKey string, uploadKey string) (GateKeeper, error) {
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

func (gk GateKeeper) NewMultiClient(uploadToken UploadToken, buckets ...string) (MultiClient, error) {
	// validate buckets
	tkn, err := jwt.ParseString(string(uploadToken))
	if err != nil {
		return MultiClient{}, err
	}

	var claims UploadClaims
	err = json.Unmarshal(tkn.RawClaims(), &claims)
	if err != nil {
		return MultiClient{}, err
	}

	exists := func(needle string, haystack []string) bool {
		for _, h := range haystack {
			if h == needle {
				return true
			}
		}
		return false
	}

	for _, qb := range buckets {
		if !exists(qb, claims.Buckets) {
			return MultiClient{}, fmt.Errorf("non-approved bucket was requested")
		}
	}

	return MultiClient{
		endpoint:    gk.endpoint,
		uploadToken: uploadToken,
		buckets:     buckets,
		resty:       gk.restClient,
	}, nil
}

func (gk GateKeeper) IssueUploadToken(bucket ...string) (token UploadToken, err error) {
	endpoint := fmt.Sprintf("%v/rest/v1/auth/upload", gk.endpoint)
	client := resty.New()

	upReq := UploadAuthenticationRequest{
		Token:   gk.uploadKey,
		Buckets: bucket,
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
		err = fmt.Errorf("statika request failed with statusCode %v; message: [%v]", resp.StatusCode(), resp.String())
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
