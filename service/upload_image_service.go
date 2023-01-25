package service

import (
	"closer-api-go/closerjwt"
	"closer-api-go/proto"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
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
	fmt.Println(message)
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			return stream.Send(&proto.UploadFileResponse{
				Success:  true,
				Message:  "Image uploaded successfully",
				Progress: 1.0,
			})
		}
		if err != nil {
			log.Printf("Error receiving file: %v", err)
		}
		buffer = append(buffer, data.GetImageBytes()...)
		message = data.Message
		if err := stream.Send(&proto.UploadFileResponse{
			Success:  false,
			Message:  fmt.Sprintf("Received %d bytes", len(buffer)),
			Progress: float64(len(buffer)) / float64(data.ImageSize),
		}); err != nil {
			return err
		}
	}
}
