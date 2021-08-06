package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
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

	// request must be formatted as /bucket/filename
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("Not found: %v ", r.URL.Path)))
		return
	}

	bucket := parts[1]
	fileName := parts[2]

	if r.Method == http.MethodGet {
		// read file from storage & validate it actually exists
		var reader CloseableReader
		var err error
		init := TimeIt(func() {
			reader, err = (*c.storage).Get(bucket, fileName)
		})

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(fmt.Sprintf("File not found: %v ", r.URL.Path)))
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

		_ = (*c.storage).Delete(bucket, fileName)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(responseOK())
	}
}
