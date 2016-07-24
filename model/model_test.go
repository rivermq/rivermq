// +build integration

package model_test

import (
	. "github.com/rivermq/rivermq/model"
	"strconv"
	"testing"
)

var (
	sub = Subscription{
		Type:        "messageType",
		CallbackURL: "http://localhost:1234/msg",
	}

	msg = Message{
		Type: "messageType",
		Body: []byte("{}"),
	}
)

func TestSaveMessage(t *testing.T) {
	_, err := SaveMessage(msg)
	if err != nil {
		t.Error(err)
	}
	dropMessageCollection()
}

func TestSaveSubscription(t *testing.T) {
	_, err := SaveSubscription(sub)
	if err != nil {
		t.Error(err)
	}
	dropSubscriptionCollection()
}

func TestFindSubscriptionByID(t *testing.T) {
	savedSub, err := SaveSubscription(sub)
	if err != nil {
		t.Error(err)
	}
	if &savedSub == nil {
		t.Error("Saved subscription is nil")
	}
	foundSub, err := FindSubscriptionByID(savedSub.ID)
	if err != nil {
		t.Error(err)
	}
	if &foundSub == nil {
		t.Error("Unable to find subscription for id: ", savedSub.ID)
	}
	dropSubscriptionCollection()
}

func TestFindAllSubscriptionsByType(t *testing.T) {
	for port := 9080; port < 9090; port++ {
		SaveSubscription(
			Subscription{
				Type:        "integratTestType",
				CallbackURL: "http://localhost:" + strconv.Itoa(port) + "/msg",
			},
		)
	}
	for port := 8080; port < 8090; port++ {
		SaveSubscription(
			Subscription{
				Type:        "otherType",
				CallbackURL: "http://localhost:" + strconv.Itoa(port) + "/msg",
			},
		)
	}

	results, err := FindAllSubscriptionsByType("integratTestType")
	if err != nil {
		t.Error(err)
	}
	if len(results) != 10 {
		t.Errorf("TestFindAllSubscriptionsByType:\n\tExpected 10, Found %d", len(results))
	}
	dropSubscriptionCollection()
}

func TestFindAllSubscriptions(t *testing.T) {
	for port := 9080; port < 9090; port++ {
		SaveSubscription(
			Subscription{
				Type:        "integratTestType",
				CallbackURL: "http://localhost:" + strconv.Itoa(port) + "/msg",
			},
		)
	}
	results, err := FindAllSubscriptions()
	if err != nil {
		t.Error(err)
	}
	if len(results) != 10 {
		t.Errorf("TestFindAllSubscriptions:\n\tExpected 10, Found %d", len(results))
	}
}

func TestDeleteSubscriptionByID(t *testing.T) {
	sub, err := SaveSubscription(sub)
	err = DeleteSubscriptionByID(sub.ID)
	if err != nil {
		t.Error(err)
	}
}

func dropMessageCollection() {
	session := DBSession.Copy()
	defer session.Close()
	c := session.DB(DBName).C(MessageCollection)

	c.DropCollection()
}

func dropSubscriptionCollection() {
	session := DBSession.Copy()
	defer session.Close()
	c := session.DB(DBName).C(SubscriptionCollection)

	c.DropCollection()
}
