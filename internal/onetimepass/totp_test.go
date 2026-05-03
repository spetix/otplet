package onetimepass

import (
	"errors"
	"testing"
	"time"
)

var errTest = errors.New("test error")

func TestPassCodeIsValid(t *testing.T) {
	valid := PassCode{
		Code:       "123456",
		ValidUntil: time.Now().Add(10 * time.Second),
		Remaining:  10 * time.Second,
	}

	if !valid.IsValid() {
		t.Fatal("expected valid PassCode to be valid")
	}
}

func TestPassCodeIsInvalidWhenExpiredOrErrored(t *testing.T) {
	expired := PassCode{
		Code:       "123456",
		ValidUntil: time.Now().Add(-1 * time.Second),
		Remaining:  -1 * time.Second,
	}

	if expired.IsValid() {
		t.Fatal("expected expired PassCode to be invalid")
	}

	errorCode := PassCode{
		Code:       "123456",
		ValidUntil: time.Now().Add(10 * time.Second),
		Remaining:  10 * time.Second,
		Error:      errTest,
	}

	if errorCode.IsValid() {
		t.Fatal("expected PassCode with error to be invalid")
	}
}
