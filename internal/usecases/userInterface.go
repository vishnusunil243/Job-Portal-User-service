package usecases

import (
	"github.com/vishnusunil243/Job-Portal-User-service/entities"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
)

type Usecases interface {
	UploadImage(*pb.UserImageRequest, string) (string, error)
	JobApply(entities.Jobs, []*pb.JobSkillResponse, int, int) error
	GetAppliedUsersByJobId(jobId string) ([]entities.Jobs, error)
}
