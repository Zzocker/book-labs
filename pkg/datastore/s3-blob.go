package datastore

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/Zzocker/book-labs/pkg/errors"
)

type awsS3 struct {
	bucketName string
	client     *s3.S3
	uploader   *s3manager.Uploader
}

type S3BlobStoreConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Region          string
	BucketName      string
}

func NewS3BlobDatastore(cfg *S3BlobStoreConfig) BlobStore {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(cfg.Endpoint),
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, cfg.SessionToken),
		Region:           aws.String(cfg.Region),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		panic(err)
	}

	return &awsS3{
		bucketName: cfg.BucketName,
		client:     s3.New(sess),
		uploader:   s3manager.NewUploader(sess),
	}
}

func (s *awsS3) Put(ctx context.Context, id string, data []byte, metadata map[string]string) error {
	const op = errors.Op("S3BlobDatastore.Put")
	_, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:     bytes.NewReader(data),
		Bucket:   aws.String(s.bucketName),
		Key:      aws.String(id),
		Metadata: aws.StringMap(metadata),
	})
	if err != nil {
		return errors.E(op, err, errors.CodeInternal)
	}

	return nil
}

func (s *awsS3) Get(ctx context.Context, id string) (*BlobFile, error) {
	const op = errors.Op("S3BlobDatastore.Get")

	obj, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(id),
	})
	if err != nil {
		if err, ok := err.(awserr.RequestFailure); ok {
			if err.StatusCode() == http.StatusNotFound {
				return nil, errors.E(op, fmt.Errorf("blob data not found"), errors.CodeNotFound)
			}
		}

		return nil, errors.E(op, err, errors.CodeInternal)
	}
	defer obj.Body.Close()

	data, err := ioutil.ReadAll(obj.Body)
	if err != nil {
		return nil, errors.E(op, fmt.Errorf("failed to read blob data: %w", err), errors.CodeInternal)
	}
	md := map[string]string{}
	for key, value := range obj.Metadata {
		md[key] = *value
	}

	return &BlobFile{
		Data:     data,
		Metadata: md,
	}, nil
}

func (s *awsS3) Del(ctx context.Context, id string) error {
	const op = errors.Op("S3BlobDatastore.Del")
	_, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(id),
	})
	if err != nil {
		return errors.E(op, err, errors.CodeInternal)
	}

	err = s.client.WaitUntilObjectNotExistsWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(id),
	})
	if err != nil {
		return errors.E(op, fmt.Errorf("failed to wait for object to be deleted : %w", err), errors.CodeInternal)
	}

	return nil
}
