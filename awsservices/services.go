package awsservices

import (
	"bytes"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pulpfree/gdps-fs-sum-dwnld/config"

	log "github.com/sirupsen/logrus"
)

// S3Service struct
type S3Service struct {
	cfg     *config.Config
	session *session.Session
}

// NewS3 function
func NewS3(cfg *config.Config) (serv *S3Service, err error) {

	serv = &S3Service{
		cfg: cfg,
	}

	serv.session, err = session.NewSession(&aws.Config{
		Region: aws.String(cfg.AWSRegion),
	})
	if err != nil {
		return nil, err
	}

	return serv, err
}

// PutFile method
func (s *S3Service) PutFile(prefix string, file *bytes.Buffer) (key string, err error) {

	uploader := s3manager.NewUploader(s.session)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:             aws.String(s.cfg.S3Bucket),
		Key:                aws.String(prefix),
		Body:               file,
		ContentType:        aws.String("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"),
		ContentDisposition: aws.String("attachment"),
	})
	if err != nil {
		return "", err
	}

	return prefix, nil
}

// GetSignedURL method
func (s *S3Service) GetSignedURL(prefix string, file *bytes.Buffer) (signedURL string, err error) {

	_, err = s.PutFile(prefix, file)
	if err != nil {
		log.Errorf("Failed to upload file: %s", err.Error())
		return "", err
	}

	svc := s3.New(s.session)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.cfg.S3Bucket),
		Key:    aws.String(prefix),
	})

	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Errorf("Failed to sign request: %s", err.Error())
		return "", err
	}

	return urlStr, err
}
