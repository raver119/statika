package main

import (
	"fmt"
	"net/http"
	"strings"
)

type Catcher struct {
	storage *Storage
	pa      PersistenceAgent
}

func NewCatcher(storage *Storage, pa PersistenceAgent) (c Catcher, err error) {
	return Catcher{storage: storage, pa: pa}, nil
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
		reader, err := (*c.storage).Get(bucket, fileName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(fmt.Sprintf("File not found: %v ", r.URL.Path)))
			return
		}

		// return it to the end user
		w.WriteHeader(http.StatusOK)
		_ = TransferBytes(reader, w)
		_ = reader.Close()
	} else {
		// delete method
		// validate access first
		authToken := r.Header.Get("Authorization")
		ok := c.pa.CheckUploadToken(authToken, bucket)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		_ = (*c.storage).Delete(bucket, fileName)
		w.WriteHeader(http.StatusOK)
		w.Write(responseOK())
	}
}
