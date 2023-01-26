package utils

import (
	"closer-api-go/closerjwt"
	"closer-api-go/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthenticateAndGetUserFromMetaData(md metadata.MD) (model.User, error) {
	tokenStr, ok := md["authorization"]
	if !ok {
		return model.User{}, status.Errorf(codes.Unauthenticated, "missing authorization token")
	}
	user, err := closerjwt.VerifyToken(tokenStr[0])
	if err != nil {
		return model.User{}, status.Errorf(codes.Unauthenticated, "missing authorization token")
	}
	return user, nil
}
