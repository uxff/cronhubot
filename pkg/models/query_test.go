package models

import (
	"testing"
)

func TestIsEmpty(t *testing.T) {
	q := NewQuery(StatusActive, "exp")

	if q.IsEmpty() {
		t.Fail()
	}
}
