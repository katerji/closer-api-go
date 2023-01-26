package service

import (
	"closer-api-go/awsclient"
	"closer-api-go/closerjwt"
	"closer-api-go/proto"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nfnt/resize"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type Server struct {
	proto.UnimplementedFileUploaderServer
}

func (s *Server) InitGrpc() {
	lis, err := net.Listen("tcp", ":1030")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterFileUploaderServer(grpcServer, &Server{})
	err = grpcServer.Serve(lis)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *Server) UploadImage(stream proto.FileUploader_UploadImageServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.Unauthenticated, "missing metadata")
	}
	tokenStr, ok := md["authorization"]
	if !ok {
		return status.Errorf(codes.Unauthenticated, "missing authorization token")
	}
	user, err := closerjwt.VerifyToken(tokenStr[0])
	fmt.Println(user)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}
	file, _ := os.CreateTemp("", "")
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {

		}
	}(file.Name())
	var fileExtension string
	var chatId int
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			if err != nil {
				fmt.Println(err)
			}
			fileName := fmt.Sprintf("%d/%d/%d%s", chatId, user.Id, time.Now().Unix(), fileExtension)
			UploadToS3(file, fileName)
			return stream.Send(&proto.UploadFileResponse{
				Success:           true,
				Progress:          1.0,
				Base64EncodedBlur: blurAndResizeImage(file),
			})
		}
		if err != nil {
			log.Printf("Error receiving file: %v", err)
		}
		file.Write(data.GetImageBytes())
		stat, _ := file.Stat()
		if err := stream.Send(&proto.UploadFileResponse{
			Success:  false,
			Progress: float64(stat.Size()) / float64(data.ImageSize),
		}); err != nil {
			return err
		}
		fileExtension = data.FileExtension
		chatId = int(data.ChatId)
	}
}

func blurAndResizeImage(file *os.File) string {
	_, err := file.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	imageToResize, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	resizedImage := resize.Resize(9, 0, imageToResize, resize.Lanczos3)
	tempFile, _ := os.CreateTemp("", "")
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			fmt.Println(err)
		}
	}(tempFile.Name())
	err = jpeg.Encode(tempFile, resizedImage, nil)
	readFile, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(readFile)
}

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
