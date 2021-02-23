package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
)

var reName = regexp.MustCompile("\\..*?$")
var reXt = regexp.MustCompile("^.*?\\.")

/*
	Local storage can be used as primary storage or mirror backup storage
*/
type LocalStorage struct {
	rootFolder string
}

func NewLocalStorage(root string) LocalStorage {
	return LocalStorage{rootFolder: root}
}

func (s LocalStorage) masterFileName(bucket string, fileName string) string {
	b, f := EncodePath(bucket, fileName)
	folder := fmt.Sprintf("%v/%v", s.rootFolder, b)
	if !FileExists(folder, false) {
		_ = os.MkdirAll(folder, 0755)
	}

	return fmt.Sprintf("%v/%v", b, f)
}

func (s LocalStorage) Get(bucket string, name string) (r CloseableReader, err error) {
	path := s.rootFolder + "/" + s.masterFileName(bucket, name)
	if !FileExists(path, true) {
		err = fmt.Errorf("requested file doesn't exist: [%v/%v]", bucket, name)
		return
	}

	r, err = os.Open(path)
	return
}

func (s LocalStorage) Put(bucket string, name string, r io.Reader) (fileName string, err error) {
	fileName = s.masterFileName(bucket, name)
	f, err := os.Create(s.rootFolder + "/" + fileName)
	if err != nil {
		return
	}

	err = TransferBytes(r, f)
	if err != nil {
		return
	}

	err = f.Close()
	return
}

func (s LocalStorage) Delete(bucket string, name string) (err error) {
	err = os.Remove(s.rootFolder + "/" + s.masterFileName(bucket, name))
	return
}

func (s LocalStorage) List(bucket string) (f []FileEntry, err error) {
	// bucket must be base64-encoded
	b := base64.StdEncoding.EncodeToString([]byte(bucket))

	files, err := ioutil.ReadDir(s.rootFolder + "/" + b)
	if err != nil {
		return nil, err
	}

	for _, v := range files {

		// FIXME: reconsider this eventually
		if v.IsDir() {
			continue
		}

		// get name and extension separately
		name := reName.ReplaceAllString(v.Name(), "")
		ext := reXt.ReplaceAllString(v.Name(), "")

		// filename was b64-encoded, so decode it
		dec, err := base64.StdEncoding.DecodeString(name)
		if err != nil {
			return nil, err
		}

		f = append(f, FileEntry{FileName: string(dec) + "." + ext})
	}

	// output will be alpha-sorted, always
	sort.Slice(f, func(i, j int) bool {
		return f[i].FileName < f[j].FileName
	})

	return
}
