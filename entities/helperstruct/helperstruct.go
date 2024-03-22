package helperstruct

import "github.com/google/uuid"

type SkillHelper struct {
	CategoryId   int
	CategoryName string
	SkillId      int
	SkillName    string
}
type JobHelper struct {
	JobId  uuid.UUID
	Status string
}
