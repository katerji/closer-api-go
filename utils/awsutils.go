package utils

import (
	"closer-api-go/awsclient"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

func UploadToS3(file *os.File, fileName string) {
	s3Client := awsclient.GetS3Client()
	putObjectInput := s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(awsclient.Bucket),
		Key:    aws.String(fileName),
	}
	_, err := s3Client.PutObject(&putObjectInput)
	if err != nil {
		fmt.Println(err)
	}
}

func GetFileFromS3(s3Path string) (s3.GetObjectOutput, error) {
	output, err := awsclient.GetS3Client().GetObject(&s3.GetObjectInput{
		Bucket: aws.String(awsclient.Bucket),
		Key:    aws.String(s3Path),
	})
	if err != nil {
		fmt.Println(err)
		return s3.GetObjectOutput{}, err
	}
	return *output, nil
}
