package s3

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"net/url"
)

type File struct {
	Name        string
	ContentType string
	Data        []byte
}

type Connector struct {
	session    *session.Session
	bucket     string
	disableSsl bool
}

func NewConnector(s *session.Session, bucket string, disableSsl bool) *Connector {
	return &Connector{
		session:    s,
		bucket:     bucket,
		disableSsl: disableSsl,
	}
}

func (c *Connector) Upload(file File) (string, error) {
	out, err := s3manager.NewUploader(c.session).Upload(&s3manager.UploadInput{
		Bucket:      aws.String(c.bucket),
		ContentType: aws.String(file.ContentType),
		Key:         aws.String(file.Name),
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

	if c.disableSsl {
		u.Scheme = "http"
	} else {
		u.Scheme = "https"
	}

	return u.String(), nil
}
