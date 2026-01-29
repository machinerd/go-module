package s3

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type Config struct {
	Bucket       string
	SourcePrefix string 
	DestPrefix   string 
	CdnUri       string
	Key          string
}

type S3Service struct {
	svc    s3iface.S3API
	config Config
}

func NewS3Service(svc s3iface.S3API, cfg Config) *S3Service {
	return &S3Service{
		svc:    svc,
		config: cfg,
	}
}
