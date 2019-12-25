package scheduler

import (
	"testing"

	"github.com/uxff/cronhubot/pkg/mocks"
	"github.com/uxff/cronhubot/pkg/models"
)

func TestScheduleAll(t *testing.T) {
	repoMock := mocks.NewJobRepo()
	s := New(repoMock)
	s.ScheduleAll()
}

func TestSchedulerCreate(t *testing.T) {
	s := New(mocks.NewJobRepo())
	c := &models.Event{Id: 1, Expression: "* * * * *"}
	if err := s.Create(c); err != nil {
		t.Fail()
	}
}

func TestSchedulerFind(t *testing.T) {
	s := New(mocks.NewJobRepo())
	c := &models.Event{Id: 1, Expression: "* * * * *"}
	if err := s.Create(c); err != nil {
		t.Fail()
	}

	_, err := s.Find(c.Id)
	if err != nil {
		t.Fail()
	}
}

func TestSchedulerUpdate(t *testing.T) {
	s := New(mocks.NewJobRepo())
	c := &models.Event{Id: 1, Expression: "* * * * *"}
	if err := s.Create(c); err != nil {
		t.Fail()
	}

	c.Status = "active"
	if err := s.Update(c); err != nil {
		t.Fail()
	}
}

func TestSchedulerDelete(t *testing.T) {
	s := New(mocks.NewJobRepo())
	c := &models.Event{Id: 1, Expression: "* * * * *"}
	if err := s.Create(c); err != nil {
		t.Fail()
	}

	if err := s.Delete(c.Id); err != nil {
		t.Fail()
	}
}
