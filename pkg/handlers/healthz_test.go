package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nbari/violetear"
	"github.com/uxff/cronhubot/pkg/checker"
	"github.com/uxff/cronhubot/pkg/mocks"
)

func TestHealthzIndex(t *testing.T) {
	checkers := map[string]checker.Checker{
		"api":      mocks.NewCheckerMock(),
		"postgres": mocks.NewCheckerMock(),
	}
	h := NewHealthzHandler(checkers)

	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fail()
	}

	r := violetear.New()
	r.HandleFunc("/healthz", h.HealthzIndex, "GET")
	r.ServeHTTP(res, req)

	response := make(map[string]bool)
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fail()
	}

	if response["api"] != true {
		t.Fail()
	}

	if response["postgres"] != true {
		t.Fail()
	}
}
