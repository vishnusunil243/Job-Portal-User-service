package service

import (
	"context"
	"fmt"

	"github.com/vishnusunil243/Job-Portal-User-service/entities"
	"github.com/vishnusunil243/Job-Portal-User-service/repository/Interface"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
)

type UserService struct {
	repo Interface.UserRepo
	pb.UnimplementedUserServiceServer
}

func NewUserService(repo Interface.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}
func (user *UserService) UserSignup(ctx context.Context, req *pb.UserSignupRequest) (*pb.UserSignupResponse, error) {
	if req.Email == "" {
		return nil, fmt.Errorf("email can't be empty")
	}
	reqEntity := entities.User{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}
	res, err := user.repo.UserSignup(reqEntity)
	if err != nil {
		return nil, err
	}
	return &pb.UserSignupResponse{
		Id:    res.ID.String(),
		Name:  res.Name,
		Email: res.Email,
		Phone: res.Phone,
	}, nil
}
