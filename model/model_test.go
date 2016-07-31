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
		Type:   "messageType",
		Status: "ACCEPTED",
		Body:   []byte("{}"),
	}
)

func TestSaveMessage(t *testing.T) {
	defer dropMessageCollection()
	_, err := SaveMessage(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestFindMessageByID(t *testing.T) {
	defer dropMessageCollection()
	savedMsg, err := SaveMessage(msg)
	if err != nil {
		t.Error(err)
	}
	if &savedMsg == nil {
		t.Error("SavedMsg is nil")
	}
	foundMsg, err := FindMessageByID(savedMsg.ID)
	if err != nil {
		t.Error(err)
	}
	if &foundMsg == nil {
		t.Error("Unable to find message for id: ", savedMsg.ID)
	}
}

func TestUpdateMessage(t *testing.T) {
	defer dropMessageCollection()
	savedMsg, err := SaveMessage(msg)
	if savedMsg.Status != StatusAccepted {
		t.Errorf("Invalid Status. Expected %s, found %s", StatusAccepted, savedMsg.Status)
	}
	savedMsg.Status = StatusInspected
	UpdateMessage(savedMsg)

	foundMsg, err := FindMessageByID(savedMsg.ID)
	if err != nil {
		t.Error(err)
	}
	if &foundMsg == nil {
		t.Error("Unable to find message for id: ", savedMsg.ID)
	}
	if foundMsg.Status != StatusInspected {
		t.Errorf("Invalid Status.  Expected %s, found %s", StatusInspected, foundMsg.Status)
	}
}

func TestFindMessageByStatus(t *testing.T) {
	defer dropMessageCollection()
	SaveMessage(Message{
		Type:   "messageType",
		Status: StatusFailedDelivery,
		Body:   []byte("{}"),
	})
	SaveMessage(Message{
		Type:   "messageType",
		Status: StatusAccepted,
		Body:   []byte("{}"),
	})

	messages, err := FindMessageByStatus(StatusFailedDelivery)
	if err != nil {
		t.Error(err)
	}
	if len(messages) != 1 {
		t.Errorf("Error [TestFIndMessageByStatus] Expected %s, found %s", 1, len(messages))
	}
}

func TestSaveSubscription(t *testing.T) {
	defer dropSubscriptionCollection()
	_, err := SaveSubscription(sub)
	if err != nil {
		t.Error(err)
	}
}

func TestFindSubscriptionByID(t *testing.T) {
	defer dropSubscriptionCollection()
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
}

func TestFindAllSubscriptionsByType(t *testing.T) {
	defer dropSubscriptionCollection()
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
}

func TestFindAllSubscriptions(t *testing.T) {
	defer dropSubscriptionCollection()
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
	defer dropSubscriptionCollection()
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
