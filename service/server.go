package main

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"sync"
)

type Engine struct {
	// root folder
	root string

	// auth keys
	keyMaster string
	keyUpload string

	port int

	// all fields below will be instantiated internally
	pa           PersistenceAgent
	router       *mux.Router
	static       *http.Server
	handlers     *ApiHandler
	srv          *http.Server
	wg           *sync.WaitGroup
	startedAsync bool
}

func CreateEngine(keyMaster string, keyUpload string, rootFolder string, port int) (e Engine, err error) {
	mh := GetEnvOrDefault("MEMCACHED_HOST", "localhost")
	pa, err := NewPersistenceAgent(mh, 11211)
	if err != nil {
		return
	}

	handlers, err := NewApiHandler(keyMaster, keyUpload, rootFolder)
	if err != nil {
		return
	}

	router, err := buildRouter(handlers, rootFolder, pa)
	if err != nil {
		return
	}

	e = Engine{
		root:      rootFolder,
		keyMaster: keyMaster,
		keyUpload: keyUpload,
		port:      port,
		pa:        pa,
		srv:       &http.Server{Addr: ":" + strconv.Itoa(port), Handler: router},
		wg:        &sync.WaitGroup{},
	}
	return
}

func buildRouter(handlers *ApiHandler, rootFolder string, pa PersistenceAgent) (router *mux.Router, err error) {
	router = mux.NewRouter()

	api := router.PathPrefix("/rest/v1").Subrouter()

	// API endpoints
	api.HandleFunc("/auth/upload", CorsHandler(handlers.LoginUpload)).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/auth/master", CorsHandler(handlers.LoginMaster)).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/ping", CorsHandler(handlers.Ping)).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	api.HandleFunc("/file", CorsHandler(handlers.Upload)).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/meta/{bucket}/{FileId}", CorsHandler(handlers.GetMeta)).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/meta/{bucket}/{FileId}", CorsHandler(handlers.UpdateMeta)).Methods(http.MethodPost, http.MethodOptions)

	// catch-all handler for static files serving
	catcher, err := NewCatcher(rootFolder, pa)
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

		e.srv.ListenAndServe()
	}()
	return
}

func (e Engine) Start() error {
	e.startedAsync = false
	return e.srv.ListenAndServe()
}

func (e Engine) Stop() (err error) {
	if e.startedAsync {
		e.srv.Shutdown(context.TODO())
		// do it
		e.wg.Wait()
	}
	e.startedAsync = false
	return
}
