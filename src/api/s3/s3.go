package s3

import (
	"app/src/utils"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadFile(filename string, content string, conf utils.Config) {
	bucket := aws.String(conf.Backup.S3.BucketName)

	key := aws.String(filename)

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(conf.Backup.S3.AccessKey, conf.Backup.S3.SecretKey, ""),
		Endpoint:         aws.String(conf.Backup.S3.Endpoint),
		Region:           aws.String(conf.Backup.S3.Region),
		S3ForcePathStyle: aws.Bool(conf.Backup.S3.S3ForcePathStyle),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		panic(err)
	}

	s3Client := s3.New(newSession)

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader(content),
		Bucket: bucket,
		Key:    key,
	})

	if err != nil {
		fmt.Printf("Failed to upload object %s%s, %s\n", *bucket, *key, err.Error())
		return
	}
	fmt.Printf("Successfully uploaded key %s\n", *key)
}

// not used
func keyExists(bucket string, key string, conf utils.Config) (bool, error) {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(conf.Backup.S3.AccessKey, conf.Backup.S3.SecretKey, ""),
		Endpoint:         aws.String(conf.Backup.S3.Endpoint),
		Region:           aws.String(conf.Backup.S3.Region),
		S3ForcePathStyle: aws.Bool(conf.Backup.S3.S3ForcePathStyle),
	}
	newSession := session.New(s3Config)

	s3Client := s3.New(newSession)

	_, err := s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				return false, nil
			default:
				return false, err
			}
		}
		return false, err
	}
	return true, nil
}
