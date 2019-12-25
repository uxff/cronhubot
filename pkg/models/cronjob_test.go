package models

import (
	"encoding/json"
	"testing"
)

func TestValidate_Defaults(t *testing.T) {
	e := NewCronJob()
	errors, ok := e.Validate()

	if ok {
		t.Fail()
	}

	if len(errors) == 0 {
		t.Fail()
	}

	if _, ok := errors["url"]; !ok {
		t.Fail()
	}

	if _, ok := errors["expression"]; !ok {
		t.Fail()
	}

	if _, ok := errors["status"]; ok {
		t.Fail()
	}

	if _, ok := errors["retries"]; ok {
		t.Fail()
	}
}

func TestValidate_Status(t *testing.T) {
	e := NewCronJob()
	e.Status = "invalid"
	errors, ok := e.Validate()

	if ok {
		t.Fail()
	}

	if len(errors) == 0 {
		t.Fail()
	}

	if _, ok := errors["status"]; !ok {
		t.Fail()
	}
}

func TestValidate_Retries(t *testing.T) {
	e := NewCronJob()
	e.Retries = 0
	errors, ok := e.Validate()

	if ok {
		t.Fail()
	}

	if len(errors) == 0 {
		t.Fail()
	}

	if _, ok := errors["retries"]; !ok {
		t.Fail()
	}

	e = NewCronJob()
	e.Retries = 11
	errors, ok = e.Validate()

	if _, ok = errors["retries"]; !ok {
		t.Fail()
	}
}

func TestSetAttributes(t *testing.T) {
	e := NewCronJob()
	newEvent := &CronJob{
		Url:            "http://newapi.io",
		Expression:     "1 1 1 1 1",
		Status:         Inactive,
		Retries:        5,
		RequestTimeout: 10,
	}
	e.SetAttributes(newEvent)

	if e.Url != newEvent.Url {
		t.Fail()
	}

	if e.Expression != newEvent.Expression {
		t.Fail()
	}

	if e.Status != newEvent.Status {
		t.Fail()
	}

	if e.Retries != newEvent.Retries {
		t.Fail()
	}

	if e.RequestTimeout != newEvent.RequestTimeout {
		t.Fail()
	}
}

func TestCronJobs_CheckExpression(t *testing.T) {
	expression := "* * * * * *"

	e := NewCronJob()
	if err := e.CheckExpression(expression); err != nil {
		t.Fail()
	}
}

func TestNewEvent(t *testing.T) {
	event := NewCronJob()

	b := []byte(`{"expire_time":"2019-01-01 11:11:11"}`)
	err := json.Unmarshal(b, event)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", event)
}
