package service

import (
	"closer-api-go/proto"
	"closer-api-go/utils"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type GrpcServer struct {
	proto.UnimplementedFileServiceServer
}

func (s *GrpcServer) InitGrpc() {
	lis, err := net.Listen("tcp", ":1030")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterFileServiceServer(grpcServer, &GrpcServer{})
	err = grpcServer.Serve(lis)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *GrpcServer) UploadImage(stream proto.FileService_UploadImageServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.Unauthenticated, "missing metadata")
	}
	user, err := utils.AuthenticateAndGetUserFromMetaData(md)
	if err != nil {
		return err
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
			utils.UploadToS3(file, fileName)
			return stream.Send(&proto.UploadImageResponse{
				Success:           true,
				Progress:          1.0,
				Base64EncodedBlur: utils.BlurAndResizeImage(file),
			})
		}
		if err != nil {
			log.Printf("Error receiving file: %v", err)
		}
		file.Write(data.GetImageBytes())
		stat, _ := file.Stat()
		if err := stream.Send(&proto.UploadImageResponse{
			Success:  false,
			Progress: utils.ToFixed(float64(stat.Size())/float64(data.ImageSize), 2),
		}); err != nil {
			return err
		}
		fileExtension = data.FileExtension
		chatId = int(data.ChatId)
	}
}

func (s *GrpcServer) DownloadImage(request *proto.DownloadImageRequest, stream proto.FileService_DownloadImageServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.Unauthenticated, "missing metadata")
	}
	user, err := utils.AuthenticateAndGetUserFromMetaData(md)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "Invalid token")
	}

	messageId := request.GetMessageId()
	message, err := GetMessageById(int(messageId))
	if err != nil {
		stream.Send(&proto.DownloadImageResponse{Message: "message not found", Success: false})
		return nil
	}
	if !IsUserInChat(message.ChatId, user.Id) {
		return status.Errorf(codes.Unauthenticated, "Unauthorized")
	}
	s3Output, err := utils.GetFileFromS3(message.S3Path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var n int
	var totalBytesSent int
	bytes := make([]byte, 1024*1024)
Loop:
	for {
		n, err = io.ReadFull(s3Output.Body, bytes)
		totalBytesSent += n
		switch err {
		case nil:
		case io.EOF:
			err := stream.Send(&proto.DownloadImageResponse{
				ImageBytes: nil,
				Progress:   1.0,
				Success:    true,
				Message:    "",
			})
			if err != nil {
				return err
			}
			break Loop
		case io.ErrUnexpectedEOF:
			err := stream.Send(&proto.DownloadImageResponse{
				ImageBytes: nil,
				Progress:   1.0,
				Success:    true,
				Message:    "",
			})
			if err != nil {
				return err
			}
			break Loop
		default:
			fmt.Println(err)
			return status.Errorf(codes.Internal, "io.ReadAll: %v", err)
		}
		serverErr := stream.Send(&proto.DownloadImageResponse{
			ImageBytes: bytes[:n],
			Progress:   utils.ToFixed(float64(totalBytesSent)/float64(*s3Output.ContentLength), 2),
			Success:    false,
			Message:    "",
		})
		if serverErr != nil {
			return status.Errorf(codes.Internal, "server.Send: %v", serverErr)
		}
	}
	return nil
}
