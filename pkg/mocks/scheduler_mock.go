package mocks

import (
	"github.com/robfig/cron"
	"github.com/uxff/cronhubot/pkg/models"
)

type SchedulerMock struct {
	Created bool
	Updated bool
	Deleted bool
}

func NewScheduler() *SchedulerMock {
	return &SchedulerMock{
		Created: false,
		Updated: false,
		Deleted: false,
	}
}

func (s *SchedulerMock) Create(*models.CronJob) (err error) {
	s.Created = true
	return
}

func (s *SchedulerMock) Update(event *models.CronJob) (err error) {
	s.Updated = true
	return
}

func (s *SchedulerMock) Delete(id uint) (err error) {
	s.Deleted = true
	return
}

func (s SchedulerMock) Find(id uint) (c *cron.Cron, err error) {
	return
}

func (s *SchedulerMock) ScheduleAll() {
	return
}
