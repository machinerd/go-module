package s3

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/machinerd/go-module/idgen"
)

func (s *S3Service) TransferObjectIfNotExist(filename string) error {
	destKey := fmt.Sprintf("%s/%s", s.config.DestPrefix, filename)

	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(destKey),
	}

	_, err := s.svc.HeadObject(headInput)
	if err == nil {
		return nil
	}

	sourceKey := fmt.Sprintf("%s/%s/%s", s.config.Bucket, s.config.SourcePrefix, filename)

	err = s.TransferObject(sourceKey, destKey)
	if err != nil {
		return fmt.Errorf("failed to transfer object from %s to %s: %w", sourceKey, destKey, err)
	}

	return nil
}

func (s *S3Service) ParseImgSrc(content *string, prefix string) (*string, error) {
	if content == nil {
		return nil, errors.New("content empty")
	}
	re := regexp.MustCompile(`src="([^"]+)"[^>]*>`)
	newContent := *content
	for _, match := range re.FindAllStringSubmatch(*content, -1) {
		key := match[1]
		key = strings.ReplaceAll(key, s.config.CdnUri+"/", "")
		key = filepath.Base(key)

		source := filepath.Join(s.config.Bucket, s.config.SourcePrefix, key)
		dest := filepath.Join(s.config.DestPrefix, prefix, key)

		err := s.TransferObject(source, dest)
		if err != nil {
			return nil, err
		}
		newSrc := fmt.Sprintf("%s/%s", s.config.CdnUri, dest)
		newContent = strings.ReplaceAll(newContent, match[1], newSrc)
	}
	return &newContent, nil
}

func (s *S3Service) CopyObject(filename string) error {
	fileInfoArr := strings.Split(filename, ".")

	var newFileExt string
	newFileName := idgen.MakeUUID()
	if len(fileInfoArr) > 1 {
		newFileExt = fileInfoArr[1]
		newFileName = fmt.Sprintf("%s/%s.%s", s.config.FilePrefix, newFileName, newFileExt)
	}

	sourceKey := fmt.Sprintf("%s/%s/%s", s.config.Bucket, s.config.SourcePrefix, filename)
	destKey := fmt.Sprintf("%s/%s", s.config.DestPrefix, newFileName)

	if err := s.TransferObject(sourceKey, destKey); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *S3Service) TransferObject(sourceKey string, destKey string) error {
	input := &s3.CopyObjectInput{
		Bucket:     aws.String(s.config.Bucket),
		Key:        aws.String(destKey),
		CopySource: aws.String(sourceKey),
	}
	_, err := s.svc.CopyObject(input)
	if err != nil {
		return err
	}
	return err
}
