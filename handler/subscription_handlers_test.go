package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	. "github.com/rivermq/rivermq/model"
	. "github.com/rivermq/rivermq/route"
)

const (
	validSub    = `{"type":"msgType","callbackUrl":"http://localhost/endpoint"}`
	responseSub = `{"timestamp":"","type":"msgType","callbackUrl":"http://localhost/endpoint"}`
)

func TestCreateSubscriptionHandler(t *testing.T) {
	buf := bytes.NewBufferString(validSub)
	req, _ := http.NewRequest("POST", "/subscriptions", buf)
	res := httptest.NewRecorder()
	NewRiverMQRouter().ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Errorf("CreateSubscription failure.  Expected %d to be %d", res.Code, http.StatusCreated)
	}
	var sub *Subscription
	json.NewDecoder(res.Body).Decode(&sub)
	if sub.Timestamp <= 0 {
		t.Error("Subscription save failure")
	}
	dropSubscriptionCollection()
}

func TestFindAllSubscriptionsHandler(t *testing.T) {
	for port := 9080; port < 9090; port++ {
		SaveSubscription(
			Subscription{
				Type:        "integratTestType",
				CallbackURL: "http://localhost:" + strconv.Itoa(port) + "/msg",
			},
		)
	}
	req, _ := http.NewRequest("GET", "/subscriptions", nil)
	res := httptest.NewRecorder()
	NewRiverMQRouter().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Errorf("FindAllSubscriptions Failure. Expected %d to be %d", res.Code, http.StatusOK)
	}
	// TODO: parse response json and validate size
}

func TestFindSubscriptionByIDHandler(t *testing.T) {
	sub, _ := SaveSubscription(Subscription{
		Type:        "messageType",
		CallbackURL: "http://localhost:123/site",
	})
	req, _ := http.NewRequest("GET", "/subscriptions/"+sub.ID, nil)
	res := httptest.NewRecorder()
	NewRiverMQRouter().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Errorf("FindSubscription failure.  Expected %d to be %d", res.Code, http.StatusOK)
	}
	dropSubscriptionCollection()
}

func TestDeleteSubscriptionByIDHandler(t *testing.T) {
	sub, _ := SaveSubscription(Subscription{
		Type:        "messageType",
		CallbackURL: "http://localhost:123/site",
	})
	req, _ := http.NewRequest("DELETE", "/subscriptions/"+sub.ID, nil)
	res := httptest.NewRecorder()
	NewRiverMQRouter().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Errorf("DeleteSubscription failure.  Expected %d to be %d", res.Code, http.StatusOK)
	}
	dropSubscriptionCollection()
}

func dropSubscriptionCollection() {
	session := DBSession.Copy()
	defer session.Close()
	c := session.DB(DBName).C(SubscriptionCollection)
	c.DropCollection()
}
