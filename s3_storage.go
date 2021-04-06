package main

import (
	"io"
)

// TODO: to be implemented
type S3Storage struct {
	//
}

func (s S3Storage) Get(name string) (r CloseableReader, err error) {

	return
}

func (s S3Storage) Put(name string, r io.Reader) (fileName string, err error) {

	return
}

func (s S3Storage) Delete(name string) (err error) {
	return
}
