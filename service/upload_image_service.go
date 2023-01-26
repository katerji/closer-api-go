package service

import (
	"closer-api-go/closerjwt"
	"closer-api-go/proto"
	"encoding/base64"
	"fmt"
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
	var buffer []byte
	var message string
	var imageExtension string
	fmt.Println(message)
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fileName := fmt.Sprintf("%s%d%d%s", os.TempDir(), user.Id, time.Now().Unix(), imageExtension)
			file, err := os.Create(fileName)
			defer func(name string) {
				err := os.Remove(name)
				if err != nil {

				}
			}(file.Name())
			if err != nil {
				fmt.Println(err)
			}
			_, err = file.Write(buffer)
			if err != nil {
				fmt.Println(err)
			}
			return stream.Send(&proto.UploadFileResponse{
				Success:           true,
				Progress:          1.0,
				Base64EncodedBlur: blurAndResizeImage(file),
			})
		}
		if err != nil {
			log.Printf("Error receiving file: %v", err)
		}
		buffer = append(buffer, data.GetImageBytes()...)
		imageExtension = data.FileExtension
		message = data.Message
		if err := stream.Send(&proto.UploadFileResponse{
			Success:  false,
			Progress: float64(len(buffer)) / float64(data.ImageSize),
		}); err != nil {
			return err
		}
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
	defer tempFile.Close()
	err = jpeg.Encode(tempFile, resizedImage, nil)
	readFile, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(readFile)
}
