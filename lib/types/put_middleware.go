package types

import "net/http"

type PutMiddleware func(next http.HandlerFunc) http.HandlerFunc
