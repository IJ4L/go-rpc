package gapi

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	db "simplebank.com/db/sqlgen"
	"simplebank.com/pb"
)

func converterUser(user db.User) *pb.User {
	return &pb.User{
		Username:         user.Username,
		Fullname:         user.FullName,
		Email:            user.Email,
		PasswordChangeAt: timestamppb.New(user.PasswordChangedAt),
		CreateAt:         timestamppb.New(user.CreatedAt),
	}
}
