package mocks

import (
	"github.com/uxff/cronhubot/pkg/models"
)

type JobRepoMock struct {
	Created      bool
	Updated      bool
	Deleted      bool
	Found        bool
	Searched     bool
	ByStatus     bool
	ByExpression bool
}

func NewJobRepo() *JobRepoMock {
	return &JobRepoMock{
		Created:      false,
		Updated:      false,
		Deleted:      false,
		Found:        false,
		Searched:     false,
		ByStatus:     false,
		ByExpression: false,
	}
}

func (repo *JobRepoMock) Create(ent *models.CronJob) (err error) {
	repo.Created = true
	return
}

func (repo *JobRepoMock) FindById(id int) (ent *models.CronJob, err error) {
	repo.Found = true
	ent = &models.CronJob{Id: 1}
	return
}

func (repo *JobRepoMock) Update(ent *models.CronJob) (err error) {
	repo.Updated = true
	return
}

func (repo *JobRepoMock) Delete(ent *models.CronJob) (err error) {
	repo.Deleted = true
	return
}

func (repo *JobRepoMock) Search(sc *models.Query) (ents []models.CronJob, err error) {
	ents = append(ents, models.CronJob{Expression: "* * * * * *"})

	switch true {
	case sc.Status != "":
		repo.ByStatus = true
	case sc.Expression != "":
		repo.ByExpression = true
	default:
		repo.Searched = true
	}

	return
}
