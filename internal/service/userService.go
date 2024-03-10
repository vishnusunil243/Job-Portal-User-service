package service

import (
	"context"
	"fmt"

	"github.com/vishnusunil243/Job-Portal-User-service/entities"
	"github.com/vishnusunil243/Job-Portal-User-service/internal/adapters"
	"github.com/vishnusunil243/Job-Portal-User-service/internal/helper"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
)

type UserService struct {
	adapters adapters.AdapterInterface
	pb.UnimplementedUserServiceServer
}

func NewUserService(adapters adapters.AdapterInterface) *UserService {
	return &UserService{
		adapters: adapters,
	}
}
func (user *UserService) UserSignup(ctx context.Context, req *pb.UserSignupRequest) (*pb.UserSignupResponse, error) {
	if req.Email == "" {
		return nil, fmt.Errorf("email can't be empty")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name can't be empty")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password can't be empty")
	}
	if req.Phone == "" {
		return nil, fmt.Errorf("phone can't be empty")
	}
	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
	}
	res, err := user.adapters.UserSignup(reqEntity)
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
func (user *UserService) UserLogin(ctx context.Context, req *pb.LoginRequest) (*pb.UserSignupResponse, error) {
	if req.Email == "" {
		return &pb.UserSignupResponse{}, fmt.Errorf("please enter a valid email")
	}
	userData, err := user.adapters.GetUserByEmail(req.Email)
	if err != nil {
		return &pb.UserSignupResponse{}, err
	}
	if userData.Email == "" {
		return &pb.UserSignupResponse{}, fmt.Errorf("invalid credentials")
	}
	if !helper.CompareHashedPassword(userData.Password, req.Password) {
		return &pb.UserSignupResponse{}, fmt.Errorf("invalid credentials please try again")
	}
	return &pb.UserSignupResponse{
		Id:    userData.ID.String(),
		Name:  userData.Name,
		Email: userData.Email,
		Phone: userData.Phone,
	}, nil
}
func (user *UserService) AdminLogin(ctx context.Context, req *pb.LoginRequest) (*pb.UserSignupResponse, error) {
	if req.Email == "" {
		return &pb.UserSignupResponse{}, fmt.Errorf("please enter a valid email")
	}
	adminData, err := user.adapters.GetAdminByEmail(req.Email)
	if err != nil {
		return &pb.UserSignupResponse{}, err
	}

	if adminData.Email == "" {
		return &pb.UserSignupResponse{}, fmt.Errorf("invalid credentials")
	}
	if !helper.CompareHashedPassword(adminData.Password, req.Password) {
		return &pb.UserSignupResponse{}, fmt.Errorf("invalid credentials")
	}
	return &pb.UserSignupResponse{
		Id:    adminData.ID.String(),
		Name:  adminData.Name,
		Email: adminData.Email,
		Phone: adminData.Phone,
	}, nil
}
