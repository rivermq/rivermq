package model_test

import (
	"github.com/rivermq/rivermq/model"
	"testing"
)

var (
	validMsg = model.Message{
		Type: "messageType",
		Body: []byte("{}"),
	}
	invalidMsg = model.Message{
		Body: []byte("{}"),
	}
)

func TestValidateMessage(t *testing.T) {
	err := model.ValidateMessage(validMsg)
	if err != nil {
		t.Error(err)
	}
	err = model.ValidateMessage(invalidMsg)
	if err == nil {
		t.Error("Expected falure of invalid Message")
	}
}
