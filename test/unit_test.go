package userServiceTest

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vishnusunil243/Job-Portal-User-service/entities"
	mock_adapters "github.com/vishnusunil243/Job-Portal-User-service/internal/adapters/mockAdapters"
	"github.com/vishnusunil243/Job-Portal-User-service/internal/helper"
	"github.com/vishnusunil243/Job-Portal-User-service/internal/service"
	mock_usecases "github.com/vishnusunil243/Job-Portal-User-service/internal/usecases/mockUsecases"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
)

func TestUserLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adapter := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecases(ctrl)
	userService := service.NewUserService(adapter, usecase)
	hashedPass, _ := helper.HashPassword("valid")
	testUUID := uuid.New()
	tests := []struct {
		name               string
		request            *pb.LoginRequest
		mockGetUserByEmail func(string) (entities.User, error)
		wantError          bool
		expectedResult     *pb.UserSignupResponse
	}{
		{
			name: "Success",
			request: &pb.LoginRequest{
				Email:    "valid@gmail.com",
				Password: "valid",
			},
			mockGetUserByEmail: func(s string) (entities.User, error) {
				return entities.User{
					ID:         testUUID,
					Name:       "valid",
					Email:      "valid@gmail.com",
					Phone:      "8888888888",
					Password:   hashedPass,
					Subscribed: true,
				}, nil
			},
			wantError: false,
			expectedResult: &pb.UserSignupResponse{
				Id:    testUUID.String(),
				Name:  "valid",
				Email: "valid@gmail.com",
				Phone: "8888888888",
			},
		},
		{
			name: "Fail",
			request: &pb.LoginRequest{
				Email:    "invalid",
				Password: "invalid",
			},
			mockGetUserByEmail: func(s string) (entities.User, error) {
				return entities.User{}, nil
			},
			wantError:      true,
			expectedResult: &pb.UserSignupResponse{},
		},
		{
			name: "Blocked",
			request: &pb.LoginRequest{
				Email:    "valid@gmail.com",
				Password: "valid",
			},
			mockGetUserByEmail: func(s string) (entities.User, error) {
				return entities.User{
					ID:        testUUID,
					Name:      "valid",
					Email:     "valid@gmail.com",
					Phone:     "8888888888",
					Password:  hashedPass,
					IsBlocked: true,
				}, nil
			},
			wantError:      true,
			expectedResult: &pb.UserSignupResponse{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapter.EXPECT().GetUserByEmail(gomock.Any()).DoAndReturn(test.mockGetUserByEmail).AnyTimes().Times(1)
			result, err := userService.UserLogin(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
				if err == nil {
					t.Errorf("expected an error but didn't find one")
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, test.expectedResult, result)

			}
		})
	}
}
func TestAddSkillAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adapter := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecases(ctrl)
	userService := service.NewUserService(adapter, usecase)
	tests := []struct {
		name               string
		request            *pb.AddSkillRequest
		mockGetSkillByName func(string) (entities.Skill, error)
		wantError          bool
	}{
		{
			name: "Success",
			request: &pb.AddSkillRequest{
				Skill:      "valid",
				CategoryId: 1,
			},
			mockGetSkillByName: func(s string) (entities.Skill, error) {
				return entities.Skill{}, nil
			},
			wantError: false,
		},
		{
			name: "Fail",
			request: &pb.AddSkillRequest{
				Skill:      "valid",
				CategoryId: 1,
			},
			mockGetSkillByName: func(s string) (entities.Skill, error) {
				return entities.Skill{
					CategoryId: 1,
					Name:       "valid",
				}, nil
			},
			wantError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapter.EXPECT().GetSkillByName(gomock.Any()).DoAndReturn(test.mockGetSkillByName).AnyTimes().Times(1)
			if !test.wantError {
				adapter.EXPECT().AdminAddSkill(gomock.Any()).Return(nil).Times(1)
			}
			_, err := userService.AddSkillAdmin(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestUserSignup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adapter := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecases(ctrl)
	userService := service.NewUserService(adapter, usecase)
	tests := []struct {
		name               string
		request            *pb.UserSignupRequest
		mockGetUserByEmail func(string) (entities.User, error)
		mockGetUserByPhone func(string) (entities.User, error)
		mockUserSignup     func(entities.User) (entities.User, error)
		wantError          bool
		expectedResult     *pb.UserSignupResponse
	}{
		{
			name: "Success",
			request: &pb.UserSignupRequest{
				Email:    "valid@gmail.com",
				Name:     "valid",
				Password: "valid",
				Phone:    "8888888888",
			},
			mockGetUserByEmail: func(s string) (entities.User, error) {
				return entities.User{}, nil
			},
			mockGetUserByPhone: func(s string) (entities.User, error) {
				return entities.User{}, nil
			},
			mockUserSignup: func(u entities.User) (entities.User, error) {
				return entities.User{
					Name:  "valid",
					Email: "valid@gmail.com",
					Phone: "8888888888",
				}, nil
			},
			wantError: false,
			expectedResult: &pb.UserSignupResponse{
				Name:  "valid",
				Email: "valid@gmail.com",
				Phone: "8888888888",
			},
		},
		{
			name: "EmailNotUnique",
			request: &pb.UserSignupRequest{
				Email:    "valid@gmail.com",
				Name:     "valid",
				Password: "1234",
				Phone:    "8888888888",
			},
			mockGetUserByEmail: func(s string) (entities.User, error) {
				return entities.User{
					Name:  "valid",
					Email: "valid@gmail.com",
					Phone: "8888888888",
				}, nil
			},
			mockGetUserByPhone: func(s string) (entities.User, error) {
				return entities.User{}, nil
			},
			mockUserSignup: func(u entities.User) (entities.User, error) {
				return entities.User{}, nil
			},
			wantError: true,
			expectedResult: &pb.UserSignupResponse{
				Email: "valid@gmail.com",
				Name:  "valid",
				Phone: "8888888888",
			},
		},
		{
			name: "PhoneNotUnique",
			request: &pb.UserSignupRequest{
				Email:    "valid@gmail.com",
				Name:     "valid",
				Password: "valid",
				Phone:    "8888888888",
			},
			mockGetUserByEmail: func(s string) (entities.User, error) {
				return entities.User{}, nil
			},
			mockGetUserByPhone: func(s string) (entities.User, error) {
				return entities.User{
					Name:  "valid",
					Email: "valid@gmail.com",
					Phone: "88888888888",
				}, nil
			},
			mockUserSignup: func(u entities.User) (entities.User, error) {
				return entities.User{}, nil
			},
			wantError: true,
			expectedResult: &pb.UserSignupResponse{
				Name:                     "valid",
				Email:                    "valid@gmail.com",
				Phone:                    "8888888888",
				ExperienceInCurrentField: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapter.EXPECT().GetUserByEmail(gomock.Any()).DoAndReturn(test.mockGetUserByEmail).Times(1)
			if test.name != "EmailNotUnique" {
				adapter.EXPECT().GetUserByPhone(gomock.Any()).DoAndReturn(test.mockGetUserByPhone).AnyTimes().Times(1)
			}
			if !test.wantError {
				adapter.EXPECT().UserSignup(gomock.Any()).DoAndReturn(test.mockUserSignup).AnyTimes().Times(1)
			}
			res, err := userService.UserSignup(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				res.Id = ""
				assert.NotNil(t, res)
				assert.Equal(t, test.expectedResult, res)
			}
		})
	}
}
