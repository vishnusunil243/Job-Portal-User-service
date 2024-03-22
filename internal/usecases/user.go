package usecases

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/vishnusunil243/Job-Portal-User-service/entities"
	"github.com/vishnusunil243/Job-Portal-User-service/internal/adapters"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
)

type UserUseCase struct {
	userAdapter adapters.AdapterInterface
}

func NewUserUseCase(useradapter adapters.AdapterInterface) *UserUseCase {
	return &UserUseCase{
		userAdapter: useradapter,
	}
}
func (user *UserUseCase) UploadImage(req *pb.UserImageRequest, profileId string) (string, error) {
	minioClient, err := minio.New(os.Getenv("MINIO_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESSKEY"), os.Getenv("MINIO_SECRETKEY"), ""),
		Secure: false,
	})
	if err != nil {
		log.Print("error while initialising minio", err)
		return "", err
	}
	objectName := "images/" + req.ObjectName
	contentType := `image/jpeg`
	n, err := minioClient.PutObject(context.Background(), os.Getenv("BUCKET_NAME"), objectName, bytes.NewReader(req.ImageData), int64(len(req.ImageData)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Println("error while uploading to minio", err)
		return "", err
	}
	log.Printf("Successfully uploaded %s of size %v\n", objectName, n)
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), os.Getenv("BUCKET_NAME"), objectName, time.Second*24*60*60, nil)
	if err != nil {
		log.Println("error while generating presigned URL", err)
		return "", err
	}
	url, err := user.userAdapter.UploadProfileImage(presignedURL.String(), profileId)
	return url, err
}
func (user *UserUseCase) JobApply(req entities.Jobs, jobskills []*pb.JobSkillResponse, jobExp, userExp int) error {
	jobSkillMap := make(map[int]bool)
	for _, jobskill := range jobskills {
		jobSkillMap[int(jobskill.SkillId)] = true
	}
	profile, err := user.userAdapter.GetProfileIdByUserId(req.UserId.String())
	if err != nil {
		return err
	}
	userSkills, err := user.userAdapter.UserGetAllSkills(profile)
	if err != nil {
		return err
	}
	count := 0
	skillPresent := false
	for _, userSkill := range userSkills {
		if jobSkillMap[userSkill.SkillId] {
			skillPresent = true
			count++
		}
	}
	if !skillPresent && len(jobSkillMap) > 0 {
		return fmt.Errorf("you do not have enough skills to apply for this job")
	}
	if userExp == 0 && jobExp != 0 || userExp < jobExp {
		return fmt.Errorf("you don't have enough experience to apply to this job")
	}
	var expPoint float64
	if count <= 0 {
		req.Weightage = float64(userExp) * 60 / 100
	} else {
		if userExp != 0 && jobExp != 0 {
			expPoint = float64(userExp) / float64(jobExp)
		} else if userExp != 0 && jobExp == 0 {
			expPoint = float64(userExp)
		}
		req.Weightage = (float64(count)/float64(len(jobSkillMap)))*40/100 + expPoint*60/100
	}
	err = user.userAdapter.JobApply(req)
	if err != nil {
		return err
	}
	return nil
}
func (user *UserUseCase) GetAppliedUsersByJobId(jobId string) ([]entities.Jobs, error) {
	res, err := user.userAdapter.GetAppliedUsersByJobId(jobId)
	if err != nil {
		return []entities.Jobs{}, err
	}
	return res, nil
}
