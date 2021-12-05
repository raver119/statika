package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/raver119/statika/utils"
)

type Client struct {
	endpoint    string
	uploadToken UploadToken
	bucket      string

	resty *resty.Client
}

func (c Client) UploadFile(fileName string, reader io.Reader) (ur UploadResponse, err error) {
	fname, err := utils.ExtractFileName(fileName)
	if err != nil {
		return UploadResponse{}, err
	}

	path, err := utils.ExtractPath(fileName)
	if err != nil {
		return UploadResponse{}, err
	}

	response, err := c.resty.R().
		SetFileReader("file", fname, reader).
		SetFormData(map[string]string{
			"token":  string(c.uploadToken),
			"bucket": c.bucket,
			"path": path,
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

func (c Client) ListFiles() (f []FileEntry, err error) {
	response, err := c.resty.R().
		SetAuthToken(string(c.uploadToken)).
		Get(fmt.Sprintf("%v/rest/v1/files/%v", c.endpoint, c.bucket))
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		err = fmt.Errorf("http request returned unexpected error code: %v", response.StatusCode())
		return
	}

	var listResp ListResponse
	err = json.Unmarshal(response.Body(), &listResp)

	if err != nil {
		return nil, err
	}

	for _, v := range listResp.Files {
		f = append(f, FileEntry{FileName: v.FileName})
	}

	return
}

func (c Client) SetMeta(fileName string, meta MetaInfo) (err error) {
	path := fmt.Sprintf(fmt.Sprintf("%v/rest/v1/meta/%v/%v", c.endpoint, c.bucket, fileName))
	response, err := c.resty.R().
		SetAuthToken(string(c.uploadToken)).
		SetBody(meta).
		Post(path)

	if err != nil {
		return err
	}

	if response.StatusCode() != http.StatusOK {
		err = fmt.Errorf("http request returned unexpected url: %v; error code: %v; body: %v;", path, response.StatusCode(), string(response.Body()))
	}
	return
}

func (c Client) Get(fileName string) (r []byte, err error) {
	response, err := c.resty.R().Get(fmt.Sprintf("%v/%v/%v", c.endpoint, c.bucket, fileName))

	if err != nil {
		return nil, err
	}

	if response.StatusCode() == http.StatusNotFound {
		return nil, fmt.Errorf("file [%v] was not found", fileName)
	}

	return response.Body(), nil
}

func (c Client) GetMeta(fileName string) (meta MetaInfo, err error) {
	response, err := c.resty.R().
		SetAuthToken(string(c.uploadToken)).
		Get(fmt.Sprintf("%v/rest/v1/meta/%v/%v", c.endpoint, c.bucket, fileName))

	if err != nil {
		return nil, err
	}

	if response.StatusCode() == http.StatusNotFound {
		return map[string]string{}, nil
	}

	if response.StatusCode() != http.StatusOK {
		err = fmt.Errorf("http request returned unexpected error code: %v; %v", response.StatusCode(), response.String())
		return nil, err
	}

	err = json.Unmarshal(response.Body(), &meta)
	return
}

func (c Client) DeleteMeta(fileName string) (err error) {
	response, err := c.resty.R().
		SetAuthToken(string(c.uploadToken)).
		Delete(fmt.Sprintf("%v/rest/v1/meta/%v/%v", c.endpoint, c.bucket, fileName))
	if err != nil {
		return err
	}

	if response.StatusCode() != http.StatusOK {
		err = fmt.Errorf("http request returned unexpected error code: %v", response.StatusCode())
	}
	return
}

func (c Client) Ping() (err error) {
	response, err := c.resty.R().
		SetAuthToken(string(c.uploadToken)).
		Get(fmt.Sprintf("%v/rest/v1/ping", c.endpoint))

	if err != nil {
		return err
	}

	if response.StatusCode() == http.StatusUnauthorized {
		err = fmt.Errorf("unauthorized. probably expired token: %v", response.String())
		return
	}

	if response.StatusCode() != http.StatusOK {
		err = fmt.Errorf("http request returned unexpected error code: %v", response.StatusCode())
		return
	}

	return
}
