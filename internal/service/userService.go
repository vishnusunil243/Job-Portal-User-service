package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vishnusunil243/Job-Portal-User-service/entities"
	"github.com/vishnusunil243/Job-Portal-User-service/internal/adapters"
	"github.com/vishnusunil243/Job-Portal-User-service/internal/helper"
	"github.com/vishnusunil243/Job-Portal-User-service/internal/usecases"
	"github.com/vishnusunil243/Job-Portal-User-service/kafka"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	CompanyClient      pb.CompanyServiceClient
	NotificationClient pb.EmailServiceClient
)

type UserService struct {
	adapters adapters.AdapterInterface
	usecases usecases.Usecases
	pb.UnimplementedUserServiceServer
}

func NewUserService(adapters adapters.AdapterInterface, usecases usecases.Usecases) *UserService {
	return &UserService{
		adapters: adapters,
		usecases: usecases,
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
	check1, err := user.adapters.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if check1.Name != "" {
		return nil, fmt.Errorf("an account already exists with the given email")
	}
	check2, err := user.adapters.GetUserByPhone(req.Phone)
	if err != nil {
		return nil, err

	}
	if check2.Name != "" {
		return nil, fmt.Errorf("an account already exists with the given phone number")
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
	if userData.IsBlocked {
		return &pb.UserSignupResponse{}, fmt.Errorf("you have been blocked by the admin")
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
	check1, err := user.adapters.GetCategoryByName(req.Category)
	if err != nil {
		return nil, err
	}
	if check1.Name != "" {
		return nil, fmt.Errorf("category already exists")
	}
	err = user.adapters.AdminAddCategory(reqEntity)
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
	check1, err := user.adapters.GetCategoryByName(req.Category)
	if err != nil {
		return nil, err
	}
	if check1.Name != "" {
		return nil, fmt.Errorf("category already exists")
	}
	err = user.adapters.AdminUpdateCategory(reqEntity)
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
	check1, err := user.adapters.GetSkillByName(req.Skill)
	if err != nil {
		return nil, err
	}
	if check1.Name != "" {
		return nil, fmt.Errorf("skill already exist")
	}
	err = user.adapters.AdminAddSkill(reqEntity)
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
	check1, err := user.adapters.GetSkillByName(req.Skill)
	if err != nil {
		return nil, err
	}
	if check1.Name != "" {
		return nil, fmt.Errorf("skill already exist")
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
	check, err := user.adapters.GetSkillById(int(req.SkillId))
	if err != nil {
		return nil, err
	}
	if check.SkillId == 0 {
		return nil, fmt.Errorf("please enter a valid skill id")
	}
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	check1, err := user.adapters.GetUserSkillById(profile, int(req.SkillId))
	if err != nil {
		return nil, err
	}
	if check1.SkillId != 0 {
		return nil, fmt.Errorf("you already have added this skill please add a new one")
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
	check, err := user.adapters.GetLinkByTitle(profile, req.Title)
	if err != nil {
		return nil, err
	}
	if check.Title != "" {
		return nil, fmt.Errorf("link with the given title is already present please delete the existing one or add a new title")
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
func (user *UserService) GetSkillById(ctx context.Context, req *pb.GetSkillByIdRequest) (*pb.SkillResponse, error) {
	res, err := user.adapters.GetSkillById(int(req.Id))
	if err != nil {
		return &pb.SkillResponse{}, err
	}
	return &pb.SkillResponse{
		Id:         int32(res.SkillId),
		Skill:      res.SkillName,
		CategoryId: int32(res.CategoryId),
		Category:   res.CategoryName,
	}, nil
}
func (user *UserService) GetCategoryById(ctx context.Context, req *pb.GetCategoryByIdRequest) (*pb.UpdateCategoryRequest, error) {
	category, err := user.adapters.GetCategoryById(int(req.Id))
	if err != nil {
		return nil, err
	}
	res := &pb.UpdateCategoryRequest{
		Id:       int32(category.ID),
		Category: category.Name,
	}
	return res, nil
}
func (user *UserService) GetUser(ctx context.Context, req *pb.GetUserById) (*pb.UserSignupResponse, error) {
	userData, err := user.adapters.GetUserById(req.Id)
	if err != nil {
		return nil, err
	}
	experience, err := user.adapters.GetExperience(req.Id)
	if err != nil {
		return nil, err
	}
	res := &pb.UserSignupResponse{
		Id:                       userData.ID.String(),
		Name:                     userData.Name,
		Email:                    userData.Email,
		Phone:                    userData.Phone,
		ExperienceInCurrentField: experience,
	}
	return res, nil
}
func (user *UserService) JobApply(ctx context.Context, req *pb.JobApplyRequest) (*emptypb.Empty, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.Jobs{
		UserId: userId,
		JobId:  req.JobId,
	}
	jobs, err := user.adapters.GetAppliedJob(req.UserId, req.JobId)
	if err != nil {
		return nil, err
	}
	if jobs.JobStatusId != 0 {
		return nil, fmt.Errorf("you have already applied for this job")
	}
	jobskills, err := CompanyClient.GetAllJobSkill(context.Background(), &pb.GetJobById{
		Id: req.JobId,
	})
	if err != nil {
		return nil, err
	}
	jobSkillData := []*pb.JobSkillResponse{}
	for {
		jobSkill, err := jobskills.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		jobSkillData = append(jobSkillData, jobSkill)
	}
	jobData, err := CompanyClient.GetJob(context.Background(), &pb.GetJobById{
		Id: req.JobId,
	})
	if err != nil {
		return nil, err
	}
	if jobData.Status == "expired" {
		return nil, fmt.Errorf("this job opening has stopped accepting applications")
	}
	jobExperience := helper.GetNumberInString(jobData.Minexperience)
	userExp, err := user.adapters.GetExperience(req.UserId)
	if err != nil {
		return nil, err
	}
	userExperience := helper.GetNumberInString(userExp)
	if err := user.usecases.JobApply(reqEntity, jobSkillData, jobExperience, userExperience); err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) UserEditPhone(ctx context.Context, req *pb.EditPhoneRequest) (*emptypb.Empty, error) {
	check, err := user.adapters.GetUserByPhone(req.Phone)
	if err != nil {
		return nil, err
	}
	reqUserId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	if check.ID != uuid.Nil && check.ID != reqUserId {
		return nil, fmt.Errorf("account already exists with the given phone please provide a new phone")
	}
	reqEntity := entities.User{
		ID:    reqUserId,
		Phone: req.Phone,
	}
	if err := user.adapters.UserEditPhone(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) UserEditName(ctx context.Context, req *pb.EditNameRequest) (*emptypb.Empty, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.User{
		ID:   userId,
		Name: req.Name,
	}
	if err := user.adapters.UserEditName(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) UserAddAddress(ctx context.Context, req *pb.AddAddressRequest) (*emptypb.Empty, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		return nil, err
	}
	address, err := user.adapters.GetAddressByProfileId(profile)
	if err != nil {
		return nil, err
	}
	if address.Country != "" {
		return nil, fmt.Errorf("you have already added an address please edit the existing one")
	}
	reqEntity := entities.Address{
		Country:   req.Country,
		State:     req.State,
		District:  req.District,
		City:      req.City,
		ProfileId: profileId,
	}
	if err := user.adapters.UserAddAddress(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) UserEditAddress(ctx context.Context, req *pb.AddressResponse) (*emptypb.Empty, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.Address{
		Country:   req.Country,
		State:     req.State,
		District:  req.District,
		City:      req.City,
		ProfileId: profileId,
	}
	if err := user.adapters.UserEditAddress(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) UserGetAddress(ctx context.Context, req *pb.GetUserById) (*pb.AddressResponse, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.Id)
	if err != nil {
		return nil, err
	}
	address, err := user.adapters.GetAddressByProfileId(profile)
	if err != nil {
		return nil, err
	}
	addressId := ""
	if address.Id != uuid.Nil {
		addressId = address.Id.String()
	}
	res := &pb.AddressResponse{
		Id:       addressId,
		Country:  address.Country,
		State:    address.State,
		District: address.District,
		City:     address.City,
	}
	return res, nil
}
func (user *UserService) UserUploadProfileImage(ctx context.Context, req *pb.UserImageRequest) (*pb.UserImageResponse, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	url, err := user.usecases.UploadImage(req, profile)
	if err != nil {
		return nil, err
	}
	res := &pb.UserImageResponse{
		Url: url,
	}
	return res, nil
}
func (user *UserService) UserGetProfilePic(ctx context.Context, req *pb.GetUserById) (*pb.UserImageResponse, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.Id)
	if err != nil {
		return nil, err
	}
	image, err := user.adapters.GetProfilePic(profile)
	if err != nil {
		return nil, err
	}
	return &pb.UserImageResponse{
		Url: image,
	}, nil
}
func (user *UserService) UserAppliedJobs(req *pb.GetUserById, srv pb.UserService_UserAppliedJobsServer) error {
	jobs, err := user.adapters.GetAppliedJobs(req.Id)
	if err != nil {
		return err
	}
	for _, job := range jobs {
		jobData, err := CompanyClient.GetJob(context.Background(), &pb.GetJobById{
			Id: job.JobId.String(),
		})
		if err != nil {
			return err
		}
		res := &pb.AppliedJobResponse{
			JobId:       job.JobId.String(),
			Designation: jobData.Designation,
			Company:     jobData.Company,
			Status:      job.Status,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}
func (user *UserService) GetAppliedUsersByJobId(req *pb.JobIdRequest, srv pb.UserService_GetAppliedUsersByJobIdServer) error {
	users, err := user.usecases.GetAppliedUsersByJobId(req.JobId)
	if err != nil {
		return err
	}
	for _, usr := range users {
		userData, err := user.adapters.GetUserById(usr.UserId.String())
		if err != nil {
			return err
		}
		res := &pb.GetUserResponse{
			Id:        usr.UserId.String(),
			Name:      userData.Name,
			Email:     userData.Email,
			Phone:     userData.Phone,
			Weightage: float32(usr.Weightage),
		}
		fmt.Println(res)
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}
func (user *UserService) AddExperience(ctx context.Context, req *pb.AddExperienceRequest) (*emptypb.Empty, error) {
	err := user.adapters.AddExperience(req.UserId, req.Experience)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) AddToShortlist(ctx context.Context, req *pb.AddToShortListRequest) (*emptypb.Empty, error) {
	weightage, err := user.adapters.GetWeightage(req.UserId, req.JobId)
	if err != nil {
		return nil, err
	}
	check, err := user.adapters.GetShortlistByUserIdAndJobId(req.UserId, req.JobId)
	if err != nil {
		return nil, err
	}
	if check.UserId != uuid.Nil {
		return nil, fmt.Errorf("this user is already added")
	}
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	jobID, err := uuid.Parse(req.JobId)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.Shortlist{
		UserId:    userID,
		JobId:     jobID,
		Weightage: weightage,
	}
	if err := user.adapters.UpdateAppliedJobStatus(2, req.JobId, req.UserId); err != nil {
		return nil, err
	}
	jobData, err := CompanyClient.GetJob(context.Background(), &pb.GetJobById{
		Id: req.JobId,
	})
	if err != nil {
		return nil, err
	}
	userData, err := user.adapters.GetUserById(req.UserId)
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		_, err = NotificationClient.AddNotification(context.Background(), &pb.AddNotificationRequest{
			UserId:  req.UserId,
			Message: fmt.Sprintf(`{"message": "congrats you have been shortlisted for the position %s at company %s"}`, jobData.Designation, jobData.Company),
		})
		if err != nil {
			log.Print("error while sending notification ", err)
		}
		fmt.Println("notification sent.............")
	}()
	err = kafka.ProduceShortlistUserMessage(userData.Email, jobData.Company, jobData.Designation, req.JobId)
	if err != nil {
		fmt.Println(err)
	}
	if err := user.adapters.AddToShortlist(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) GetShortlist(req *pb.JobIdRequest, srv pb.UserService_GetShortlistServer) error {
	shortlists, err := user.adapters.GetShortlist(req.JobId)
	if err != nil {
		return err
	}
	for _, shortlist := range shortlists {
		userData, err := user.adapters.GetUserById(shortlist.UserId.String())
		if err != nil {
			return err
		}
		res := &pb.GetUserResponse{
			Id:        userData.ID.String(),
			Name:      userData.Name,
			Email:     userData.Email,
			Phone:     userData.Phone,
			Weightage: float32(shortlist.Weightage),
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}
func (user *UserService) AddEducation(ctx context.Context, req *pb.EducationRequest) (*emptypb.Empty, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	start, err := helper.ConvertStringToDate(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("please provide a valid start date")
	}
	end, err := helper.ConvertStringToDate(req.EndDate)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.Education{
		Degree:      req.Degree,
		Institution: req.Institution,
		UserId:      userId,
		StartDate:   start,
		EndDate:     end,
	}
	if err := user.adapters.AddEducation(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) EditEducation(ctx context.Context, req *pb.EducationResponse) (*emptypb.Empty, error) {
	Id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	start, err := helper.ConvertStringToDate(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("please provide a valid start date")
	}
	end, err := helper.ConvertStringToDate(req.EndDate)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.Education{
		ID:          Id,
		Degree:      req.Degree,
		Institution: req.Institution,
		StartDate:   start,
		EndDate:     end,
	}
	if err := user.adapters.EditEducation(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) GetEducation(req *pb.GetUserById, srv pb.UserService_GetEducationServer) error {
	educations, err := user.adapters.GetEducation(req.Id)
	if err != nil {
		return err
	}
	for _, education := range educations {
		res := &pb.EducationResponse{
			Id:          education.ID.String(),
			Degree:      education.Degree,
			Institution: education.Institution,
			StartDate:   education.StartDate.String(),
			EndDate:     education.EndDate.String(),
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}
func (user *UserService) BlockUser(ctx context.Context, req *pb.GetUserById) (*emptypb.Empty, error) {
	if err := user.adapters.AddToBlockList(req.Id); err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) UnblockUser(ctx context.Context, req *pb.GetUserById) (*emptypb.Empty, error) {
	if err := user.adapters.RemoveFromBlockList(req.Id); err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) RemoveEducation(ctx context.Context, req *pb.EducationById) (*emptypb.Empty, error) {
	if req.EducationId == "" {
		return nil, fmt.Errorf("please provide a valid education id")
	}
	if err := user.adapters.DeleteEducation(req.EducationId); err != nil {
		return nil, err
	}
	return nil, nil
}

func (user *UserService) InterviewScheduleForUser(ctx context.Context, req *pb.InterviewScheduleRequest) (*emptypb.Empty, error) {
	date, err := helper.ConvertStringToDate(req.Date)
	if err != nil {
		return nil, err
	}
	if time.Now().Before(date) {
		return nil, fmt.Errorf("please provide a valid date")
	}
	if err := user.adapters.UpdateAppliedJobStatus(3, req.JobId, req.UserId); err != nil {
		return nil, err
	}
	jobData, err := CompanyClient.GetJob(context.Background(), &pb.GetJobById{
		Id: req.JobId,
	})
	if err != nil {
		return nil, err
	}
	userData, err := user.adapters.GetUserById(req.UserId)
	if err != nil {
		log.Print("error getting user id")
	}
	go func() {
		NotificationClient.AddNotification(context.Background(), &pb.AddNotificationRequest{
			UserId:  req.UserId,
			Message: fmt.Sprintf(`{"message": "interview has been scheduled for %s by %s for the position of %s"}`, req.Date, jobData.Company, jobData.Designation),
		})
	}()
	if err := kafka.InterviewScheduledMessage(userData.Email, jobData.Company, req.Date, jobData.Designation); err != nil {
		log.Print("error while sending notifications", err)
	}
	if err := user.adapters.InterviewSchedule(req.UserId, req.JobId, date); err != nil {
		return nil, err
	}
	return nil, nil
}
func (user *UserService) GetInterviewsForUser(req *pb.GetUserById, srv pb.UserService_GetInterviewsForUserServer) error {
	users, err := user.adapters.GetInterviewsForUser(req.Id)
	if err != nil {
		return err
	}
	for _, user := range users {
		jobData, err := CompanyClient.GetJob(context.Background(), &pb.GetJobById{
			Id: user.JobId.String(),
		})
		if err != nil {
			return err
		}
		res := &pb.InterviewResponse{
			JobId:       user.JobId.String(),
			Designation: jobData.Designation,
			Company:     jobData.Company,
			Date:        user.InterviewDate.String(),
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}
