package concurrency

import (
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

type Concurrency struct {
	DB *gorm.DB
	mu sync.Mutex
}

func NewConcurrency(DB *gorm.DB) *Concurrency {
	return &Concurrency{
		DB: DB,
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
		}
		c.mu.Unlock()
	}()
}
