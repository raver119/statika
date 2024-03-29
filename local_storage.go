package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/raver119/statika/classes"
	"github.com/raver119/statika/utils"
)

var reName = regexp.MustCompile(`\..*?$`)
var reXt = regexp.MustCompile(`^.*?\.`)

/*
	Local storage can be used as primary storage or mirror backup storage
*/
type LocalStorage struct {
	rootFolder string
}

func NewLocalStorage(root string) *LocalStorage {
	return &LocalStorage{rootFolder: root}
}

func (s *LocalStorage) Name() string {
	return "Local storage"
}

func (s *LocalStorage) prepareFolder(bucket string) (err error) {
	b := base64.StdEncoding.EncodeToString([]byte(bucket))

	folder := fmt.Sprintf("%v/%v", s.rootFolder, b)
	if !utils.FileExists(folder, false) {
		err = os.MkdirAll(folder, 0755)
	}

	return
}

func (s *LocalStorage) Get(bucket string, name string) (r classes.CloseableReader, err error) {
	path := s.rootFolder + "/" + utils.MasterFileName(bucket, name)
	if !utils.FileExists(path, true) {
		err = errNotFound
		return
	}

	r, err = os.Open(path)
	return
}

func (s *LocalStorage) Put(bucket string, name string, r io.ReadSeeker) (fileName string, err error) {
	fileName = utils.MasterFileName(bucket, name)

	err = s.prepareFolder(bucket)
	if err != nil {
		return
	}

	f, err := os.Create(s.rootFolder + "/" + fileName)
	if err != nil {
		return
	}

	err = utils.TransferBytes(r, f)
	if err != nil {
		return
	}

	err = f.Close()
	return
}

func (s *LocalStorage) Delete(bucket string, name string) (err error) {
	err = os.Remove(s.rootFolder + "/" + utils.MasterFileName(bucket, name))

	// meta should be deleted regardless of result
	_ = os.Remove(s.rootFolder + "/" + utils.MasterMetaName(bucket, name))
	return
}

func (s *LocalStorage) List(bucket string) (f []classes.FileEntry, err error) {
	// bucket must be base64-encoded
	b := base64.StdEncoding.EncodeToString([]byte(bucket))

	files, err := ioutil.ReadDir(s.rootFolder + "/" + b)
	if err != nil {
		// just return empty array, as if there's no files
		return f, nil
	}

	for _, v := range files {
		// FIXME: reconsider this eventually
		if v.IsDir() || strings.HasSuffix(v.Name(), classes.META_EXTENSION) {
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

		f = append(f, classes.FileEntry{FileName: string(dec) + "." + ext})
	}

	// TODO: make optional sort-by-date, alpha-sort etc.
	sort.Slice(f, func(i, j int) bool {
		return f[i].FileName < f[j].FileName
	})

	return
}

func (s *LocalStorage) PutMeta(bucket string, filename string, meta classes.MetaInfo) (err error) {
	bf := utils.MasterMetaName(bucket, filename)

	err = s.prepareFolder(bucket)
	if err != nil {
		return
	}

	// TODO: check if file exists first eventually, however it might be bad for performance
	path := fmt.Sprintf("%v/%v", s.rootFolder, bf)
	file, err := os.Create(path)
	if err != nil {
		return
	}

	r, err := json.Marshal(meta)
	if err != nil {
		return
	}

	err = utils.TransferBytes(bytes.NewReader(r), file)
	if err != nil {
		return
	}

	err = file.Close()
	return
}

func (s *LocalStorage) GetMeta(bucket string, filename string) (meta classes.MetaInfo, err error) {
	bf := utils.MasterMetaName(bucket, filename)

	path := fmt.Sprintf("%v/%v", s.rootFolder, bf)
	if !utils.FileExists(path, true) {
		return map[string]string{}, err
	}

	reader, err := os.Open(path)
	if err != nil {
		return
	}

	metaBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}

	_ = reader.Close()
	err = json.Unmarshal(metaBytes, &meta)
	return
}

func (s *LocalStorage) DeleteMeta(bucket string, filename string) (err error) {
	bf := utils.MasterMetaName(bucket, filename)

	path := fmt.Sprintf("%v/%v", s.rootFolder, bf)
	err = os.Remove(path)
	return
}
