package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/nbari/violetear"
	"github.com/uxff/cronhubot/pkg/mocks"
	"github.com/uxff/cronhubot/pkg/models"
)

func TestEventsIndex(t *testing.T) {
	repoMock := mocks.NewJobRepo()
	schedulerMock := mocks.NewScheduler()
	h := NewJobsHandler(repoMock, schedulerMock)

	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/events", nil)
	if err != nil {
		t.Fail()
	}

	r := violetear.New()
	r.HandleFunc("/events", h.JobsIndex, "GET")
	r.ServeHTTP(res, req)

	events := []models.CronJob{}
	if err := json.NewDecoder(res.Body).Decode(&events); err != nil {
		t.Fail()
	}

	if len(events) == 0 {
		t.Errorf("Expected response to not be empty %s", strconv.Itoa(len(events)))
	}

	if res.Code != http.StatusOK {
		t.Errorf("Expected status %d to be equal %d", res.Code, http.StatusOK)
	}
}

func TestEventsIndexByStatus(t *testing.T) {
	schedulerMock := mocks.NewScheduler()
	repoMock := mocks.NewJobRepo()
	h := NewJobsHandler(repoMock, schedulerMock)

	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/events?status=active", nil)
	if err != nil {
		t.Fail()
	}

	r := violetear.New()
	r.HandleFunc("/events", h.JobsIndex, "GET")
	r.ServeHTTP(res, req)

	if !repoMock.ByStatus {
		t.Errorf("Expected to search by status")
	}

	if res.Code != http.StatusOK {
		t.Errorf("Expected status %d to be equal %d", res.Code, http.StatusOK)
	}
}

func TestEventsIndexByExpression(t *testing.T) {
	schedulerMock := mocks.NewScheduler()
	repoMock := mocks.NewJobRepo()
	h := NewJobsHandler(repoMock, schedulerMock)

	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/events?expression=* * * * *", nil)
	if err != nil {
		t.Fail()
	}

	r := violetear.New()
	r.HandleFunc("/events", h.JobsIndex, "GET")
	r.ServeHTTP(res, req)

	if !repoMock.ByExpression {
		t.Errorf("Expected to search by expression")
	}

	if res.Code != http.StatusOK {
		t.Errorf("Expected status %d to be equal %d", res.Code, http.StatusOK)
	}
}

func TestEventCreate(t *testing.T) {
	schedulerMock := mocks.NewScheduler()
	repoMock := mocks.NewJobRepo()
	h := NewJobsHandler(repoMock, schedulerMock)

	res := httptest.NewRecorder()
	body := strings.NewReader(`{"url":"http://foo.com","expression":"* * * * *"}`)
	req, err := http.NewRequest("POST", "/events", body)
	if err != nil {
		t.Fail()
	}

	r := violetear.New()
	r.HandleFunc("/events", h.JobsCreate, "POST")
	r.ServeHTTP(res, req)

	if !repoMock.Created {
		t.Error("Expected repo create to be called")
	}

	if !schedulerMock.Created {
		t.Error("Expected scheduler create to be called")
	}

	if res.Code != http.StatusCreated {
		t.Errorf("Expected status %d to be equal %d", res.Code, http.StatusOK)
	}
}

func TestEventShow(t *testing.T) {
	schedulerMock := mocks.NewScheduler()
	repoMock := mocks.NewJobRepo()
	h := NewJobsHandler(repoMock, schedulerMock)

	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/events/1", nil)
	if err != nil {
		t.Fail()
	}

	r := violetear.New()
	r.AddRegex(":id", `^\d+$`)
	r.HandleFunc("/events/:id", h.JobsDetail, "GET")
	r.ServeHTTP(res, req)

	if !repoMock.Found {
		t.Error("Expected repo findEventById to be called")
	}

	if res.Code != http.StatusOK {
		t.Errorf("Expected status %d to be equal %d", res.Code, http.StatusOK)
	}
}

func TestEventsUpdate(t *testing.T) {
	schedulerMock := mocks.NewScheduler()
	repoMock := mocks.NewJobRepo()
	h := NewJobsHandler(repoMock, schedulerMock)

	res := httptest.NewRecorder()
	body := strings.NewReader(`{"url":"http://foo.com","expression":"* * * * *"}`)
	req, err := http.NewRequest("PUT", "/events/1", body)
	if err != nil {
		t.Fail()
	}

	r := violetear.New()
	r.AddRegex(":id", `^\d+$`)
	r.HandleFunc("/events/:id", h.JobsUpdate, "PUT")
	r.ServeHTTP(res, req)

	if !repoMock.Updated {
		t.Error("Expected repo update to be called")
	}

	if !schedulerMock.Updated {
		t.Error("Expected scheduler update to be called")
	}

	if res.Code != http.StatusOK {
		t.Errorf("Expected status %d to be equal %d", res.Code, http.StatusOK)
	}
}

func TestEventsDelete(t *testing.T) {
	schedulerMock := mocks.NewScheduler()
	repoMock := mocks.NewJobRepo()
	h := NewJobsHandler(repoMock, schedulerMock)

	res := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/events/1", nil)
	if err != nil {
		t.Fail()
	}

	r := violetear.New()
	r.AddRegex(":id", `^\d+$`)
	r.HandleFunc("/events/:id", h.JobsDelete, "DELETE")
	r.ServeHTTP(res, req)

	if !repoMock.Deleted {
		t.Error("Expected repo update to be called")
	}

	if !schedulerMock.Deleted {
		t.Error("Expected scheduler update to be called")
	}

	if res.Code != http.StatusNoContent {
		t.Errorf("Expected status %d to be equal %d", res.Code, http.StatusNoContent)
	}
}
