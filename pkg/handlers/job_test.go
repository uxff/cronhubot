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

func TestCronjobsIndex(t *testing.T) {
	authMock := mocks.NewAuth()
	repoMock := mocks.NewJobRepo()
	schedulerMock := mocks.NewScheduler()
	h := NewJobsHandler(authMock, repoMock, schedulerMock)

	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/cronjobs", nil)
	if err != nil {
		t.Fail()
	}

	r := violetear.New()
	r.HandleFunc("/cronjobs", h.JobsIndex, "GET")
	r.ServeHTTP(res, req)

	cronjobs := []models.CronJob{}
	if err := json.NewDecoder(res.Body).Decode(&cronjobs); err != nil {
		t.Fail()
	}

	if len(cronjobs) == 0 {
		t.Errorf("Expected response to not be empty %s", strconv.Itoa(len(cronjobs)))
	}

	if res.Code != http.StatusOK {
		t.Errorf("Expected status %d to be equal %d", res.Code, http.StatusOK)
	}
}

func TestCronjobsIndexByStatus(t *testing.T) {
	authMock := mocks.NewAuth()
	schedulerMock := mocks.NewScheduler()
	repoMock := mocks.NewJobRepo()
	h := NewJobsHandler(authMock, repoMock, schedulerMock)

	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/cronjobs?status=active", nil)
	if err != nil {
		t.Fail()
	}

	r := violetear.New()
	r.HandleFunc("/cronjobs", h.JobsIndex, "GET")
	r.ServeHTTP(res, req)

	if !repoMock.ByStatus {
		t.Errorf("Expected to search by status")
	}

	if res.Code != http.StatusOK {
		t.Errorf("Expected status %d to be equal %d", res.Code, http.StatusOK)
	}
}

func TestCronjobsIndexByExpression(t *testing.T) {
	authMock := mocks.NewAuth()
	schedulerMock := mocks.NewScheduler()
	repoMock := mocks.NewJobRepo()
	h := NewJobsHandler(authMock, repoMock, schedulerMock)

	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/cronjobs?expression=* * * * *", nil)
	if err != nil {
		t.Fail()
	}

	r := violetear.New()
	r.HandleFunc("/cronjobs", h.JobsIndex, "GET")
	r.ServeHTTP(res, req)

	if !repoMock.ByExpression {
		t.Errorf("Expected to search by expression")
	}

	if res.Code != http.StatusOK {
		t.Errorf("Expected status %d to be equal %d", res.Code, http.StatusOK)
	}
}

func TestCronjobCreate(t *testing.T) {
	authMock := mocks.NewAuth()
	schedulerMock := mocks.NewScheduler()
	repoMock := mocks.NewJobRepo()
	h := NewJobsHandler(authMock, repoMock, schedulerMock)

	res := httptest.NewRecorder()
	body := strings.NewReader(`{"url":"http://foo.com","expression":"* * * * *"}`)
	req, err := http.NewRequest("POST", "/cronjobs", body)
	if err != nil {
		t.Fail()
	}

	r := violetear.New()
	r.HandleFunc("/cronjobs", h.JobsCreate, "POST")
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

func TestCronjobShow(t *testing.T) {
	authMock := mocks.NewAuth()
	schedulerMock := mocks.NewScheduler()
	repoMock := mocks.NewJobRepo()
	h := NewJobsHandler(authMock, repoMock, schedulerMock)

	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/cronjobs/1", nil)
	if err != nil {
		t.Fail()
	}

	r := violetear.New()
	r.AddRegex(":id", `^\d+$`)
	r.HandleFunc("/cronjobs/:id", h.JobsDetail, "GET")
	r.ServeHTTP(res, req)

	if !repoMock.Found {
		t.Error("Expected repo findCronjobById to be called")
	}

	if res.Code != http.StatusOK {
		t.Errorf("Expected status %d to be equal %d", res.Code, http.StatusOK)
	}
}

func TestCronjobsUpdate(t *testing.T) {
	authMock := mocks.NewAuth()
	schedulerMock := mocks.NewScheduler()
	repoMock := mocks.NewJobRepo()
	h := NewJobsHandler(authMock, repoMock, schedulerMock)

	res := httptest.NewRecorder()
	body := strings.NewReader(`{"url":"http://foo.com","expression":"* * * * *"}`)
	req, err := http.NewRequest("PUT", "/cronjobs/1", body)
	if err != nil {
		t.Fail()
	}

	r := violetear.New()
	r.AddRegex(":id", `^\d+$`)
	r.HandleFunc("/cronjobs/:id", h.JobsUpdate, "PUT")
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

func TestCronjobsDelete(t *testing.T) {
	authMock := mocks.NewAuth()
	schedulerMock := mocks.NewScheduler()
	repoMock := mocks.NewJobRepo()
	h := NewJobsHandler(authMock, repoMock, schedulerMock)

	res := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/cronjobs/1", nil)
	if err != nil {
		t.Fail()
	}

	r := violetear.New()
	r.AddRegex(":id", `^\d+$`)
	r.HandleFunc("/cronjobs/:id", h.JobsDelete, "DELETE")
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
