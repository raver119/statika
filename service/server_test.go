package main

import (
	"testing"
)

const (
	TEST_M = "TEST_MASTER_KEY"
	TEST_U = "TEST_UPLOAD_KEY"
)

func TestServer_StartStop(t *testing.T) {
	eng, err := CreateEngine(TEST_M, TEST_U, "/tmp", 80)
	if err != nil {
		t.Fatal(err)
	}

	err = eng.StartAsync()
	if err != nil {
		t.Fatal(err)
	}

	err = eng.Stop()
	if err != nil {
		t.Fatal(err)
	}
}
