package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	. "github.com/raver119/statika/classes"
	. "github.com/raver119/statika/utils"
	. "github.com/raver119/statika/wt"
)

type Catcher struct {
	storage   *Storage
	tokenizer Tokenizer
}

func NewCatcher(storage *Storage) (c Catcher, err error) {
	return Catcher{storage: storage, tokenizer: NewTokenizer()}, nil
}

/*
	This method does 3 things:
	1) Serves static files
	2) Updates traffic counters
	3) Handles DELETE requests
*/
func (c Catcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// only GET/DELETE is supported
	if r.Method != http.MethodGet && r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// direct access to metainfo files is forbidden
	if strings.HasSuffix(r.URL.Path, META_EXTENSION) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	path, err := url.QueryUnescape(r.URL.Path)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	bucket, path, err := SplitPath(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch bucket: %v", err), http.StatusInternalServerError)
		return
	}

	if path == "" {
		http.Error(w, fmt.Sprintf("failed to fetch bucket: %v", path), http.StatusBadRequest)
		return
	}

	path, err = SanitizeFileName(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to sanitize file name: %v", err), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		// read file from storage & validate it actually exists
		var reader CloseableReader
		var err error
		init := TimeIt(func() {
			reader, err = (*c.storage).Get(bucket, path)
		})

		if err == errNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("File not found: %v ", path)))
			return
		}

		if IsTimingEnabled() {
			//to get storage read time I'll fetch data into this temporary buffer
			b := bytes.NewBuffer([]byte{})
			read := TimeIt(func() {
				_ = TransferBytes(reader, b)
			})

			// report timing
			w.Header().Add("Server-Timing", fmt.Sprintf(`get;desc="Storage GET";dur=%v, read;desc="Storage READ";dur=%v`, init, read))
			w.WriteHeader(http.StatusOK)

			// and now transfer fetched data to the client
			_ = TransferBytes(b, w)
		} else {
			// return it to the end user
			w.WriteHeader(http.StatusOK)
			_ = TransferBytes(reader, w)
		}
		_ = reader.Close()
	} else {
		// delete method
		// validate access first
		authToken := r.Header.Get("Authorization")
		ok, _ := c.tokenizer.ValidateUploadToken(authToken, bucket)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		_ = (*c.storage).Delete(bucket, path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(ResponseOK())
	}
}
