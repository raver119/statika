package main

import (
	"io"

	"github.com/raver119/statika/classes"
)

/**
 * This Storage implementation acts as a proxy on top of any (fixed) number of other storages.
 * PLEASE NOTE: This is naive mirror implementation suitable for playground only. It has no error handling. It should NOT be used in production environment.
 */
type MirrorStorage struct {
	back []Storage
}

func NewMirrorStorage(backends... Storage) *MirrorStorage {
	return &MirrorStorage{back: backends}
}

func (s *MirrorStorage) Name() string {
	return "Mirror storage"
}

func (s *MirrorStorage)  Put(bucket, name string, r io.ReadSeeker) (fileName string, err error) {
	for _, b := range s.back {
		fileName, err = b.Put(bucket, name, r)
		
		_, errf := r.Seek(0, 0)
		if errf != nil {
			return "", errf
		}
	}
	return
}

func (s *MirrorStorage) Get(bucket, name string) (r classes.CloseableReader, err error) {
	for _, b := range s.back {
		// just get the first successful one
		r, err = b.Get(bucket, name)
		if err == nil {
			break
		}
	} 
	return
}

func (s *MirrorStorage)	List(bucket string) (f []classes.FileEntry, err error) {
	for _, b := range s.back {
		// just get the first successful one
		f, err = b.List(bucket)
		if err == nil {
			break
		}
	} 
	return
}

func (s *MirrorStorage) Delete(bucket, name string) (err error) {
	for _, b := range s.back {
		// FIXME: handle error, i.e. re-issue query later
		err = b.Delete(bucket, name)
	}
	return
}

func (s *MirrorStorage) PutMeta(bucket, filename string, meta classes.MetaInfo) (err error) {
	for _, b := range s.back {
		// FIXME: handle error, i.e. re-issue query later
		err = b.PutMeta(bucket, filename, meta)
	}
	return
}

func (s *MirrorStorage) GetMeta(bucket, filename string) (meta classes.MetaInfo, err error) {
	for _, b := range s.back {
		// just get the first successful one
		meta, err = b.GetMeta(bucket, filename)
		if err == nil {
			break
		}
	} 
	return
}

func (s *MirrorStorage) DeleteMeta(bucket, filename string) (err error) {
	for _, b := range s.back {
		err = b.DeleteMeta(bucket, filename)
	}
	return
}