package storage

import (
	"io"
	"log"
	"net/url"

	"github.com/minio/minio-go"
	"github.com/sarulabs/di"
	"github.com/spf13/viper"
)

type s3ServiceImpl struct {
	logger    *log.Logger
	client    *minio.Client
	endpoint  string
	bucket    string
	accessKey string
	secretKey string
	useHTTPS  bool
}

func createS3Service(c di.Container, logger *log.Logger) serviceInternal {
	bucket := viper.GetString("S3_BUCKET")
	accessKey := viper.GetString("S3_ACCESS_KEY")
	secretKey := viper.GetString("S3_SECRET_KEY")
	domain := viper.GetString("S3_DOMAIN")

	if bucket == "" || accessKey == "" || secretKey == "" {
		return nil
	}

	var endpoint string
	var useHTTPS bool
	if domain != "" {
		u, err := url.Parse(domain)
		if err != nil {
			endpoint = domain
			useHTTPS = false
		} else {
			endpoint = u.Host
			useHTTPS = u.Scheme == "https"
		}
	} else {
		endpoint = "s3.amazonaws.com"
		useHTTPS = true
	}

	impl := &s3ServiceImpl{
		logger:    logger,
		client:    nil,
		endpoint:  endpoint,
		bucket:    bucket,
		accessKey: accessKey,
		secretKey: secretKey,
		useHTTPS:  useHTTPS,
	}

	return impl
}

// Initialize service
func (s *s3ServiceImpl) Initialize() error {
	client, err := minio.New(s.endpoint, s.accessKey, s.secretKey, s.useHTTPS)
	if err != nil {
		log.Fatalln(err)
	}
	s.client = client

	done := make(chan struct{})
	defer close(done)

	channel := client.ListObjectsV2(s.bucket, "", true, done)
	for object := range channel {
		if object.Err != nil {
			s.logger.Printf("unable to access s3 service %s: %v", s.endpoint, object.Err)
			return err
		}
		break
	}

	s.logger.Printf("using s3 service %s as storage (bucket \"%s\")", s.endpoint, s.bucket)
	return nil
}

// Upload new file
func (s *s3ServiceImpl) Upload(filename FileRef, reader io.Reader) (FileRef, error) {
	n, err := s.client.PutObject(s.bucket, string(filename), reader, -1, minio.PutObjectOptions{})
	if err != nil {
		s.logger.Printf("unable to write s3 file \"%s:%s\": %v", s.bucket, filename, err)
		return emptyFileRef, err
	}

	s.logger.Printf("new s3 file has been written: \"%s:%s\" (%d bytes)", s.bucket, filename, n)
	return FileRef(filename), nil
}

// Download existing file
func (s *s3ServiceImpl) Download(filename FileRef) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(s.bucket, string(filename), minio.GetObjectOptions{})
	if err != nil {
		s.logger.Printf("unable to read s3 file \"%s:%s\": %v", s.bucket, filename, err)
		return nil, err
	}

	return obj, nil
}

// List existing files
func (s *s3ServiceImpl) List() ([]FileRef, error) {
	items := make([]FileRef, 0)

	done := make(chan struct{})
	defer close(done)

	channel := s.client.ListObjectsV2(s.bucket, "", true, done)
	for object := range channel {
		if object.Err != nil {
			s.logger.Printf("unable to access s3 service: %v", object.Err)
			return nil, object.Err
		}

		items = append(items, FileRef(object.Key))
	}

	return items, nil
}

// Delete existing file
func (s *s3ServiceImpl) Delete(filename FileRef) error {
	err := s.client.RemoveObject(s.bucket, string(filename))
	if err != nil {
		e, ok := err.(minio.ErrorResponse)

		if ok && e.Code == "NoSuchKey" {
			s.logger.Printf("won't remove s3 file \"%s:%s\" since it doesn't exist", s.bucket, filename)
			return nil
		}

		s.logger.Printf("unable to remove s3 file \"%s:%s\": %v", s.bucket, filename, err)
		return err
	}

	s.logger.Printf("s3 file \"%s:%s\" has been removed", s.bucket, filename)
	return nil
}
