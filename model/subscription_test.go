package model_test

import (
	"github.com/rivermq/rivermq/model"
	"testing"
)

var (
	validSub = model.Subscription{
		Type:        "messageType",
		CallbackURL: "http://localhost:1234/msg",
	}

	invalidSub = model.Subscription{
		Type:        "messageType",
		CallbackURL: "http//localhost:1234/msg",
	}
)

func TestValidateSubscription(t *testing.T) {
	err := model.ValidateSubscription(validSub)
	if err != nil {
		t.Error(err)
	}
	err = model.ValidateSubscription(invalidSub)
	if err == nil {
		t.Error("Expected falure of invalid Subscription")
	}
}
