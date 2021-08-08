package api

import (
	"github.com/go-resty/resty/v2"
)

type MultiClient struct {
	endpoint    string
	uploadToken UploadToken
	buckets     []string

	resty *resty.Client
}

func (mc MultiClient) Check(bucket string) (ok bool) {
	for _, v := range mc.buckets {
		if v == bucket {
			return true
		}
	}

	return false
}

func (mc MultiClient) Bucket(bucket string) (c Client) {
	return Client{
		endpoint:    mc.endpoint,
		uploadToken: mc.uploadToken,
		bucket:      bucket,
		resty:       mc.resty,
	}
}
