package statika

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"net/http"
)

type Client struct {
	endpoint    string
	uploadToken UploadToken
	bucket      string

	resty *resty.Client
}

func (c Client) UploadFile(fileName string, reader io.Reader) (ur UploadResponse, err error) {
	response, err := c.resty.R().
		SetFileReader("file", fileName, reader).
		SetFormData(map[string]string{
			"token":  string(c.uploadToken),
			"bucket": c.bucket,
		}).
		Post(fmt.Sprintf("%v/rest/v1/file", c.endpoint))

	if err != nil {
		return UploadResponse{}, err
	}

	if response.StatusCode() != http.StatusOK {
		return UploadResponse{}, fmt.Errorf("failed to upload file, ErrorCode: %v; Message: %v", response.StatusCode(), string(response.Body()))
	}

	err = json.Unmarshal(response.Body(), &ur)
	return
}

func (c Client) DeleteFile(fileName string) (err error) {
	response, err := c.resty.R().
		SetAuthToken(string(c.uploadToken)).
		Delete(fmt.Sprintf("%v%v", c.endpoint, fileName))
	if err != nil {
		return err
	}

	if response.StatusCode() != http.StatusOK {
		err = fmt.Errorf("http request returned unexpected error code: %v", response.StatusCode())
	}

	return
}
