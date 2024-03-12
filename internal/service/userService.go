package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vishnusunil243/Job-Portal-User-service/entities"
	"github.com/vishnusunil243/Job-Portal-User-service/internal/adapters"
	"github.com/vishnusunil243/Job-Portal-User-service/internal/helper"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
	"google.golang.org/protobuf/types/known/emptypb"
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
func (user *UserService) CreateProfile(ctx context.Context, req *pb.GetUserById) (*emptypb.Empty, error) {
	if err := user.adapters.CreateProfile(req.Id); err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
func (user *UserService) AddCategory(ctx context.Context, req *pb.AddCategoryRequest) (*emptypb.Empty, error) {
	reqEntity := entities.Category{
		Name: req.Category,
	}
	fmt.Println(req.Category)
	err := user.adapters.AdminAddCategory(reqEntity)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
func (user *UserService) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*emptypb.Empty, error) {
	reqEntity := entities.Category{
		ID:   int(req.Id),
		Name: req.Category,
	}
	err := user.adapters.AdminUpdateCategory(reqEntity)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return nil, nil
}
func (user *UserService) GetAllCategory(req *emptypb.Empty, srv pb.UserService_GetAllCategoryServer) error {
	categories, err := user.adapters.GetAllCategory()
	if err != nil {
		return err
	}
	for _, category := range categories {
		res := &pb.UpdateCategoryRequest{
			Id:       int32(category.ID),
			Category: category.Name,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}
func (user *UserService) AddSkillAdmin(ctx context.Context, req *pb.AddSkillRequest) (*emptypb.Empty, error) {
	reqEntity := entities.Skill{
		CategoryId: int(req.CategoryId),
		Name:       req.Skill,
	}
	err := user.adapters.AdminAddSkill(reqEntity)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) AdminUpdateSkill(ctx context.Context, req *pb.SkillResponse) (*emptypb.Empty, error) {
	reqEntity := entities.Skill{
		ID:         int(req.Id),
		CategoryId: int(req.CategoryId),
		Name:       req.Skill,
	}
	if err := user.adapters.AdminUpdateSkill(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) GetAllSkills(e *emptypb.Empty, srv pb.UserService_GetAllSkillsServer) error {
	skills, err := user.adapters.AdminGetAllSkills()
	if err != nil {
		return err
	}
	for _, skill := range skills {
		res := &pb.SkillResponse{
			Id:         int32(skill.SkillId),
			Skill:      skill.SkillName,
			CategoryId: int32(skill.CategoryId),
			Category:   skill.CategoryName,
		}
		err := srv.Send(res)
		if err != nil {
			return err
		}
	}
	return nil
}
func (user *UserService) GetAllSkillsUser(req *pb.GetUserById, srv pb.UserService_GetAllSkillsUserServer) error {
	profileId, err := user.adapters.GetProfileIdByUserId(req.Id)
	if err != nil {
		return err
	}
	skills, err := user.adapters.UserGetAllSkills(profileId)
	if err != nil {
		return err
	}
	for _, skill := range skills {
		res := &pb.SkillResponse{
			Id:         int32(skill.SkillId),
			CategoryId: int32(skill.CategoryId),
			Skill:      skill.SkillName,
			Category:   skill.CategoryName,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}
func (user *UserService) AddSkillUser(ctx context.Context, req *pb.DeleteSkillRequest) (*emptypb.Empty, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.UserSkill{
		ProfileId: profileId,
		SkillId:   int(req.SkillId),
	}
	if err := user.adapters.UserAddSkill(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) DeleteSkillUser(ctx context.Context, req *pb.DeleteSkillRequest) (*emptypb.Empty, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.UserSkill{
		ProfileId: profileId,
		SkillId:   int(req.SkillId),
	}
	if err := user.adapters.UserDeleteSkill(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) AddLinkUser(ctx context.Context, req *pb.AddLinkRequest) (*emptypb.Empty, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.Link{
		ProfileId: profileId,
		Title:     req.Title,
		Url:       req.Url,
	}
	if err := user.adapters.AddLink(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) DeleteLinkUser(ctx context.Context, req *pb.DeleteLinkRequest) (*emptypb.Empty, error) {
	if err := user.adapters.DeleteLink(req.Id); err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) GetAllLinksUser(req *pb.GetUserById, srv pb.UserService_GetAllLinksUserServer) error {
	profile, err := user.adapters.GetProfileIdByUserId(req.Id)
	if err != nil {
		return err
	}

	links, err := user.adapters.GetAllLinksUser(profile)
	if err != nil {
		return err
	}
	for _, link := range links {
		res := &pb.LinkResponse{
			Id:    link.ID.String(),
			Title: link.Title,
			Url:   link.Url,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}
