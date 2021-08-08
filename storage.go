package main

import (
	"github.com/raver119/statika/classes"
	"io"
)

type Storage interface {
	Name() string

	Put(bucket string, name string, r io.ReadSeeker) (fileName string, err error)
	Get(bucket string, name string) (r classes.CloseableReader, err error)
	List(bucket string) (f []classes.FileEntry, err error)
	Delete(bucket string, name string) (err error)

	PutMeta(bucket string, filename string, meta classes.MetaInfo) (err error)
	GetMeta(bucket string, filename string) (meta classes.MetaInfo, err error)
	DeleteMeta(bucket string, filename string) (err error)
}
