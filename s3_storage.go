package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/raver119/statika/classes"
	"github.com/raver119/statika/utils"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

type S3Storage struct {
	bucket    string
	awsConfig *aws.Config
	s3client  *s3.S3
	mode      string
}

var allowedAccessModes = []string{"private", "public-read", "public-read-write", "authenticated-read"}

func validMode(mode string) bool {
	for _, v := range allowedAccessModes {
		if v == mode {
			return true
		}
	}

	return false
}

func NewS3Storage(bucket string, endpoint string, region string) (s S3Storage, err error) {
	spacesKey := utils.GetEnvOrDefault("S3_KEY", "")
	spacesSecret := utils.GetEnvOrDefault("S3_SECRET", "")

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(spacesKey, spacesSecret, ""),
		Endpoint:    aws.String(endpoint),
		Region:      aws.String(region),
	}

	// private by default
	mode := utils.GetEnvOrDefault("ACL", "private")
	if !validMode(mode) {
		return S3Storage{}, fmt.Errorf("unknown ACL provided: [%v]", mode)
	}

	// this client will be reused
	c, err := buildClient(s3Config)
	if err != nil {
		return S3Storage{}, err
	}

	return S3Storage{bucket: bucket, awsConfig: s3Config, s3client: c, mode: mode}, nil
}

// endpoint looks like "https://nyc3.digitaloceanspaces.com", region is hardcoded
func NewSpacesStorage(bucket string, endpoint string) (s S3Storage, err error) {
	return NewS3Storage(bucket, endpoint, "us-east-1")
}

func (s S3Storage) Name() string {
	return "S3 storage"
}

func buildClient(awsConfig *aws.Config) (c *s3.S3, err error) {
	newSession, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}

	c = s3.New(newSession)
	return
}

func (s S3Storage) client() (c *s3.S3, err error) {
	return s.s3client, nil
}

func (s S3Storage) Get(bucket string, name string) (r classes.CloseableReader, err error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(bucket + "/" + name),
	}

	result, err := c.GetObject(input)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	return result.Body, err
}

func (s S3Storage) Put(bucket string, name string, r io.ReadSeeker) (fileName string, err error) {
	c, err := s.client()
	if err != nil {
		return fileName, err
	}

	// use this bucket to upload file
	object := s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(bucket + "/" + name),
		Body:   r,
		ACL:    aws.String(s.mode),
		Metadata: map[string]*string{
			"x-amz-meta-my-key": aws.String("your-value"), //required?
		},
	}

	_, err = c.PutObject(&object)
	return fmt.Sprintf("%v/%v", bucket, name), err
}

func (s S3Storage) List(bucket string) (f []classes.FileEntry, err error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}

	input := &s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(bucket + "/"),
	}

	objects, err := c.ListObjects(input)
	if err != nil {
		return nil, err
	}

	for _, obj := range objects.Contents {
		// strip bucket name from the file name
		f = append(f, classes.FileEntry{FileName: strings.Replace(aws.StringValue(obj.Key), bucket+"/", "", 1)})
	}

	return
}

func (s S3Storage) Delete(bucket string, name string) (err error) {
	c, err := s.client()
	if err != nil {
		return err
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(bucket + "/" + name),
	}

	_, err = c.DeleteObject(input)
	return err
}

func (s S3Storage) PutMeta(bucket string, filename string, meta classes.MetaInfo) (err error) {
	b, _ := json.Marshal(meta)
	_, err = s.Put(bucket, filename+classes.META_EXTENSION, bytes.NewReader(b))
	return
}

func (s S3Storage) GetMeta(bucket string, filename string) (meta classes.MetaInfo, err error) {
	b, err := s.Get(bucket, filename+classes.META_EXTENSION)
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			// just a 404, return empty MetaInfo
			return classes.MetaInfo{}, nil
		}
		return nil, err
	}

	body, _ := ioutil.ReadAll(b)
	err = json.Unmarshal(body, &meta)
	return
}

func (s S3Storage) DeleteMeta(bucket string, filename string) (err error) {
	err = s.Delete(bucket, filename+classes.META_EXTENSION)
	return
}
