package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
)

type S3Storage struct {
	bucket    string
	awsConfig *aws.Config
}

func NewS3Storage(bucket string, endpoint string, region string) (s S3Storage, err error) {
	spacesKey := GetEnvOrDefault("S3_KEY", "")
	spacesSecret := GetEnvOrDefault("S3_SECRET", "")

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(spacesKey, spacesSecret, ""),
		Endpoint:    aws.String(endpoint),
		Region:      aws.String(region),
	}

	return S3Storage{bucket: bucket, awsConfig: s3Config}, nil
}

// endpoint looks like "https://nyc3.digitaloceanspaces.com", region is hardcoded
func NewSpacesStorage(bucket string, endpoint string) (s S3Storage, err error) {
	return NewS3Storage(bucket, endpoint, "us-east-1")
}

func (s S3Storage) client() (c *s3.S3, err error) {
	newSession, err := session.NewSession(s.awsConfig)
	if err != nil {
		return nil, err
	}

	c = s3.New(newSession)
	return
}

func (s S3Storage) Get(bucket string, name string) (r CloseableReader, err error) {
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
		ACL:    aws.String("private"), // all files will be accessed through proxy anyway
		Metadata: map[string]*string{
			"x-amz-meta-my-key": aws.String("your-value"), //required
		},
	}

	_, err = c.PutObject(&object)
	return fmt.Sprintf("%v/%v", bucket, name), err
}

func (s S3Storage) List(bucket string) (f []FileEntry, err error) {
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
		f = append(f, FileEntry{FileName: aws.StringValue(obj.Key)})
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

func (s S3Storage) PutMeta(bucket string, filename string, meta MetaInfo) (err error) {
	return
}

func (s S3Storage) GetMeta(bucket string, filename string) (meta MetaInfo, err error) {
	return
}

func (s S3Storage) DeleteMeta(bucket string, filename string) (err error) {
	return
}
