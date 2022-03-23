package s3

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"net/url"
)

type File struct {
	Name        string
	Format      string
	ContentType string
	Data        []byte
}

type Connector struct {
	session      *session.Session
	bucket       string
	EnsureBucket bool
}

func NewConnector(s *session.Session, bucket string, ensureBucket bool) *Connector {
	return &Connector{
		session:      s,
		bucket:       bucket,
		EnsureBucket: ensureBucket,
	}
}

func (c *Connector) Upload(file File) (string, error) {
	key := fmt.Sprintf("%s%s", file.Name, file.Format)

	out, err := s3manager.NewUploader(c.session).Upload(&s3manager.UploadInput{
		Bucket:      aws.String(c.bucket),
		ContentType: aws.String(file.ContentType),
		Key:         aws.String(key),
		ACL:         aws.String("public-read"),
		Body:        bytes.NewBuffer(file.Data),
	})
	if err != nil {
		return "", err
	}

	u, err := url.Parse(out.Location)
	if err != nil {
		return "", err
	}

	if u.Scheme != "https" {
		u.Scheme = "https"
	}

	return u.String(), nil
}
