package main

import "io"

type Storage interface {
	Put(bucket string, name string, r io.ReadSeeker) (fileName string, err error)
	Get(bucket string, name string) (r CloseableReader, err error)
	List(bucket string) (f []FileEntry, err error)
	Delete(bucket string, name string) (err error)

	PutMeta(bucket string, filename string, meta MetaInfo) (err error)
	GetMeta(bucket string, filename string) (meta MetaInfo, err error)
	DeleteMeta(bucket string, filename string) (err error)
}
