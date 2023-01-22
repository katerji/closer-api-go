package awsclient

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"mime/multipart"
	"os"
)

var s3Client *s3.S3
const Bucket = "closer-media"

func GetS3Client() *s3.S3 {
	if s3Client == nil {
		accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
		secret := os.Getenv("AWS_SECRET_ACCESS_KEY")
		region := os.Getenv("AWS_DEFAULT_REGION")
		awsSession, _ := session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Credentials: credentials.NewStaticCredentials(
					accessKey,
					secret,
					"",
				),
				Region: aws.String(region),
			},
		})
		s3Client = s3.New(awsSession)
	}
	return s3Client
}

func UploadToS3(file multipart.File, fileName string) {
	s3Client := GetS3Client()
	putObjectInput := s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(Bucket),
		Key:    aws.String(fileName),
	}
	_, err := s3Client.PutObject(&putObjectInput)
	if err != nil {
		fmt.Println(err)
	}
}
