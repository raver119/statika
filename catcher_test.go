package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/raver119/statika/wt"
	"github.com/stretchr/testify/require"
)

func TestCatcher_ServeHTTP(t *testing.T) {
	var storage Storage = NewLocalStorage("/tmp")
	_, err := storage.Put("pictures", "filename.txt", strings.NewReader("pictures"))
	require.NoError(t, err)

	_, err = storage.Put("pictures", "sub/filename.txt", strings.NewReader("sub pictures"))
	require.NoError(t, err)

	_, err = storage.Put("pictures", "sub/file name.txt", strings.NewReader("sub pictures"))
	require.NoError(t, err)

	c := Catcher{
		storage:   &storage,
		tokenizer: DevTokenizer("123456"),
	}

	tests := []struct {
		name     string
		r        *http.Request
		wantCode int
		wantBody string
	}{
		{"test_0", httptest.NewRequest(http.MethodGet, "http://localhost/images/filename.txt", nil), http.StatusNotFound, ""},
		{"test_1", httptest.NewRequest(http.MethodGet, "http://localhost/pictures/filename.txt", nil), http.StatusOK, "pictures"},
		{"test_2", httptest.NewRequest(http.MethodGet, "http://localhost/pictures/../filename.txt", nil), http.StatusOK, "pictures"},
		{"test_3", httptest.NewRequest(http.MethodGet, "http://localhost/pictures/sub/filename.txt", nil), http.StatusOK, "sub pictures"},
		{"test_4", httptest.NewRequest(http.MethodGet, "http://localhost/filename.txt", nil), http.StatusBadRequest, ""},
		{"test_5", httptest.NewRequest(http.MethodGet, "http://localhost/../filename.txt", nil), http.StatusInternalServerError, ""},
		{"test_6", httptest.NewRequest(http.MethodGet, "http://localhost/pictures/sub/file+name.txt", nil), http.StatusOK, "sub pictures"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.ServeHTTP(w, tt.r)

			require.Equal(t, tt.wantCode, w.Code)

			if tt.wantCode == http.StatusOK {
				require.Equal(t, tt.wantBody, w.Body.String())
			}
		})
	}
}
