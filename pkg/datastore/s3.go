package datastore

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
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
	downloader *s3manager.Downloader
}

type S3FileDatastoreConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Region          string
	BucketName      string
}

func NewS3FileDatastore(cfg *S3FileDatastoreConfig) FileDataStore {
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
		downloader: s3manager.NewDownloader(sess),
	}
}

func (s *awsS3) Put(ctx context.Context, fl *File) error {
	const op = errors.Op("S3FileDatastore.Put")
	_, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:   bytes.NewReader(fl.Data),
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fl.ID),
	})
	if err != nil {
		return errors.E(op, err, errors.CodeUnexpected, errors.SeverityError)
	}

	return nil
}

func (s *awsS3) Get(ctx context.Context, id string) (*File, error) {
	const op = errors.Op("S3FileDatastore.Get")
	f := File{
		ID: id,
	}
	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := s.downloader.DownloadWithContext(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(id),
	})
	if err != nil {
		return nil, errors.E(op, err, errors.CodeUnexpected, errors.SeverityError)
	}
	f.Data = buf.Bytes()

	return &f, nil
}

func (s *awsS3) Delete(ctx context.Context, id string) error {
	const op = errors.Op("S3FileDatastore.Delete")
	_, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(id),
	})
	if err != nil {
		return errors.E(op, err, errors.CodeUnexpected, errors.SeverityError)
	}

	err = s.client.WaitUntilObjectNotExistsWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(id),
	})
	if err != nil {
		return errors.E(op, fmt.Errorf("failed to wait for object to be deleted : %w", err), errors.CodeUnexpected, errors.SeverityError)
	}

	return nil
}
