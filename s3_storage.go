package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
)

type S3Storage struct {
	awsConfig *aws.Config
}

func NewS3Storage(endpoint string, region string) (s S3Storage, err error) {
	spacesKey := GetEnvOrDefault("S3_KEY", "")
	spacesSecret := GetEnvOrDefault("S3_SECRET", "")

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(spacesKey, spacesSecret, ""),
		Endpoint:    aws.String(endpoint),
		Region:      aws.String(region),
	}

	return S3Storage{awsConfig: s3Config}, nil
}

// endpoint looks like "https://nyc3.digitaloceanspaces.com", region is hardcoded
func NewSpacesStorage(endpoint string) (s S3Storage, err error) {
	return NewS3Storage(endpoint, "us-east-1")
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
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
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

	// try to create bucket first
	params := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}

	_, err = c.CreateBucket(params)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				// do nothing, it already exists
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				// same, do nothing
			default:
				// all other errors indicate real error
				return fileName, err
			}
		}
	}

	// use this bucket to upload file
	object := s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
		Body:   r,
		ACL:    aws.String("private"), // all files will be accessed through proxy anyway
		Metadata: map[string]*string{
			"x-amz-meta-my-key": aws.String("your-value"), //required
		},
	}

	_, err = c.PutObject(&object)
	return
}

func (s S3Storage) List(bucket string) (f []FileEntry, err error) {
	c, err := s.client()
	if err != nil {
		return nil, err
	}

	input := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
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
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
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
