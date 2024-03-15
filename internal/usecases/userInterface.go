package usecases

import "github.com/vishnusunil243/Job-Portal-proto-files/pb"

type Usecases interface {
	UploadImage(*pb.UserImageRequest, string) (string, error)
}
