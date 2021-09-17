package main

import (
	"context"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/raver119/statika/utils"
)

type Engine struct {
	storage *Storage

	// auth keys
	keyMaster string
	keyUpload string

	port int

	srv          *http.Server
	wg           *sync.WaitGroup
	startedAsync bool
}

func CreateEngine(keyMaster string, keyUpload string, storage *Storage, port int) (e Engine, err error) {
	if err != nil {
		return
	}

	handlers, err := NewApiHandler(keyMaster, keyUpload, storage)
	if err != nil {
		return
	}

	router, err := buildRouter(handlers, storage)
	if err != nil {
		return
	}

	e = Engine{
		storage:   storage,
		keyMaster: keyMaster,
		keyUpload: keyUpload,
		port:      port,
		srv:       &http.Server{Addr: ":" + strconv.Itoa(port), Handler: router},
		wg:        &sync.WaitGroup{},
	}
	return
}

func buildRouter(handlers *ApiHandler, storage *Storage) (router *mux.Router, err error) {
	router = mux.NewRouter()

	api := router.PathPrefix("/rest/v1").Subrouter()

	// API endpoints
	api.HandleFunc("/auth/upload", utils.CorsHandler(handlers.LoginUpload)).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/auth/master", utils.CorsHandler(handlers.LoginMaster)).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/ping", utils.CorsHandler(handlers.Ping)).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	api.HandleFunc("/file", utils.CorsHandler(handlers.Upload)).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/files/{bucket}", utils.CorsHandler(handlers.List)).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/meta/{bucket}/{fileName}", utils.CorsHandler(handlers.GetMeta)).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/meta/{bucket}/{fileName}", utils.CorsHandler(handlers.SetMeta)).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/meta/{bucket}/{fileName}", utils.CorsHandler(handlers.DeleteMeta)).Methods(http.MethodDelete, http.MethodOptions)

	// catch-all handler for static files serving
	catcher, err := NewCatcher(storage)
	if err != nil {
		return
	}

	router.PathPrefix("/").Handler(catcher)

	return
}

func (e Engine) StartAsync() (err error) {
	e.startedAsync = true
	e.wg.Add(1)

	go func() {
		defer e.wg.Done()

		err = e.srv.ListenAndServe()
	}()
	return
}

func (e Engine) Start() error {
	e.startedAsync = false
	return e.srv.ListenAndServe()
}

func (e Engine) Stop() (err error) {
	if e.startedAsync {
		err = e.srv.Shutdown(context.TODO())
		// do it
		e.wg.Wait()
	}
	e.startedAsync = false
	return
}
