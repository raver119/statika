package main

import (
	"io"

	"github.com/raver119/statika/classes"
)

/**
 * This Storage implementation acts as a proxy on top of any (fixed) number of other storages
 */
type MirrorStorage struct {
	back []Storage
}

func (s *MirrorStorage) Name() string {
	return "Mirror storage"
}

func (s *MirrorStorage)  Put(bucket, name string, r io.ReadSeeker) (fileName string, err error) {
	for _, b := range s.back {
		b.Put(bucket, name, r)
	}
	return
}

func (s *MirrorStorage) Get(bucket string, name string) (r classes.CloseableReader, err error) {
	return
}

func (s *MirrorStorage)	List(bucket string) (f []classes.FileEntry, err error) {
	return
}

func (s *MirrorStorage) Delete(bucket string, name string) (err error) {
	return
}

func (s *MirrorStorage) PutMeta(bucket string, filename string, meta classes.MetaInfo) (err error) {
	return
}

func (s *MirrorStorage) GetMeta(bucket string, filename string) (meta classes.MetaInfo, err error) {
	return
}

func (s *MirrorStorage) DeleteMeta(bucket string, filename string) (err error) {
	return
}