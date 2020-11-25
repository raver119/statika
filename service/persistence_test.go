package main

import (
	"testing"
)

var (
	host = GetEnvOrDefault("MEMCACHED_HOST", "localhost")
)




/*
	PLEASE NOTE: This set of tests requires memcached instance to be running on MEMCACHED_HOST env var (which defaults to 127.0.0.1)
*/

func TestPersistenceAgent_CreateUploadToken(t *testing.T) {
	pa, err := NewPersistenceAgent(host, 11211)
	if err != nil {
		t.Fatal(err)
	}

	req := UploadAuthenticationRequest{Token:  "SOME_TOKEN", Bucket: "SOME_BUCKET", }
	token, err := pa.CreateUploadToken(req)
	if err != nil {
		t.Fatal(err)
	}

	if !pa.CheckUploadToken(token, req.Bucket) {
		t.Fatal("unexpected access failed")
	}

	if pa.CheckUploadToken(token, "random bucked name") {
		t.Fatal("unexpected access granted")
	}

	if pa.CheckUploadToken("random access token", req.Bucket) {
		t.Fatal("unexpected access granted")
	}

	if pa.CheckUploadToken("random access token", "random bucked name") {
		t.Fatal("unexpected access granted")
	}

	err = pa.TouchUploadToken(token)

}
