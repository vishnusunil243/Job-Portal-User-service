package concurrency

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/vishnusunil243/Job-Portal-User-service/entities"
	"github.com/vishnusunil243/Job-Portal-User-service/internal/adapters"
	"github.com/vishnusunil243/Job-Portal-User-service/internal/service"
	"github.com/vishnusunil243/Job-Portal-User-service/kafka"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
	"gorm.io/gorm"
)

type Concurrency struct {
	DB       *gorm.DB
	adapters adapters.AdapterInterface
	mu       sync.Mutex
}

func NewConcurrency(DB *gorm.DB, adapters adapters.AdapterInterface) *Concurrency {
	return &Concurrency{
		DB:       DB,
		adapters: adapters,
	}
}
func (c *Concurrency) Concurrency() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			c.mu.Lock()
			if err := c.DB.Exec(`
			UPDATE users SET is_blocked=true WHERE
			report_count>50
			`).Error; err != nil {
				log.Print("error while performing concurrency", err)
				break
			}
			if err := c.DB.Exec(`
			UPDATE jobs SET job_status_id=5 WHERE 
			NOW()>interview_date+INTERVAL '48 hours'
			AND job_status_id NOT IN (4)
			`).Error; err != nil {
				log.Print("error while performing concurrency ", err)
				break
			}
			if err := c.DB.Exec(`
			UPDATE shortlists SET status='rejected' 
			WHERE NOW()>interview_date+INTERVAL '48 hours'
			AND status NOT IN ('hired')
			`).Error; err != nil {
				log.Print(err.Error())
				break
			}
			if err := c.sendWarningNotifications(); err != nil {
				log.Print("error sending warning notifications ", err)
				break
			}
		}
		c.mu.Unlock()
	}()
}
func (c *Concurrency) sendWarningNotifications() error {
	var shortlists []entities.Shortlist
	err := c.DB.Raw(`
 SELECT *
 FROM shortlists 
 WHERE interview_date BETWEEN NOW() AND NOW()+INTERVAL '20 minutes'
 AND warning_sent=false
 `).Scan(&shortlists).Error
	if err != nil {
		return err
	}

	for _, shortlist := range shortlists {

		jobData, err := service.CompanyClient.GetJob(context.Background(), &pb.GetJobById{
			Id: shortlist.JobId.String(),
		})
		if err != nil {
			return err
		}
		if _, err := service.NotificationClient.AddNotification(context.Background(), &pb.AddNotificationRequest{
			UserId:  shortlist.UserId.String(),
			Message: fmt.Sprintf(`{"message":"Please be aware of your interview scheduled at %v by the company %s for the position %s with roomId %s"}`, shortlist.InterviewDate.String(), jobData.Company, jobData.Designation, shortlist.RoomId),
		}); err != nil {
			return err
		}
		userData, err := c.adapters.GetUserById(shortlist.UserId.String())
		if err != nil {
			log.Print("error retrieving user info ", err)
		}
		cmpny, err := service.CompanyClient.GetCompanyIdFromJobId(context.Background(), &pb.GetJobById{
			Id: shortlist.JobId.String(),
		})
		if err != nil {
			log.Println("error retrieving company id", err)
		}
		companyData, err := service.CompanyClient.GetCompanyById(context.Background(), &pb.GetJobByCompanyId{
			Id: cmpny.Id,
		})
		if err != nil {
			log.Println("error retrieving company info ", err)
		}
		if err := kafka.WarningEmail(userData.Name, jobData.Company, jobData.Designation, shortlist.InterviewDate.String(), companyData.Email, userData.Email, shortlist.RoomId); err != nil {
			log.Print("error sending messages ", err)
		}
		if err := c.DB.Exec(`
		UPDATE shortlists SET warning_sent=true WHERE id=?
		`, shortlist.ID).Error; err != nil {
			return err
		}
	}
	return nil
}
