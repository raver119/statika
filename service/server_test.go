package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	TEST_M = "TEST_MASTER_KEY"
	TEST_U = "TEST_UPLOAD_KEY"
)

func TestServer_StartStop(t *testing.T) {
	eng, err := CreateEngine(TEST_M, TEST_U, "/tmp", 80)
	require.NoError(t, err)

	err = eng.StartAsync()
	require.NoError(t, err)

	err = eng.Stop()
	require.NoError(t, err)
}
