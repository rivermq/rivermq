// +build integration

package model_test

import (
	"github.com/rivermq/rivermq/model"
	"testing"
)

var (
	sub = model.Subscription{
		Type:        "messageType",
		CallbackURL: "http://localhost:1234/msg",
	}

	msg = model.Message{
		Type: "messageType",
		Body: []byte("{}"),
	}
)

func TestSaveMessage(t *testing.T) {
	_, err := model.SaveMessage(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSaveSubscription(t *testing.T) {
	_, err := model.SaveSubscription(sub)
	if err != nil {
		t.Error(err)
	}
}
