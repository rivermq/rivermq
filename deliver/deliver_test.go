// +build integration

package deliver_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/rivermq/rivermq/model"
	. "github.com/rivermq/rivermq/route"
)

const (
	validMsg = `{"type":"msgType","body":{"state":"CO"}}`
)

// TestAttemptMessageDistrobution attempts to send a message to a test server
func TestAttemptMessageDistrobution(t *testing.T) {
	defer dropMessageCollection()
	defer dropSubscriptionCollection()
	resultChannel := make(chan bool, 1)
	// Server that will subscribe to a given message
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resultChannel <- true
		fmt.Fprintln(w, "Hello, client")
	}))
	defer testServer.Close()

	// Create subscription in RiverMQ
	sub := Subscription{
		Type:        "msgType",
		CallbackURL: testServer.URL,
	}
	subJSON, err := json.Marshal(sub)
	if err != nil {
		t.Error(err)
	}
	buf := bytes.NewBuffer(subJSON)
	req, _ := http.NewRequest("POST", "/subscriptions", buf)
	res := httptest.NewRecorder()
	NewRiverMQRouter().ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Errorf("CreateSubscription failure.  Expected %d to be %d", res.Code, http.StatusCreated)
	}

	// Send message to RiverMQ
	msgBuf := bytes.NewBufferString(validMsg)
	msgReq, _ := http.NewRequest("POST", "/messages", msgBuf)
	msgRes := httptest.NewRecorder()
	NewRiverMQRouter().ServeHTTP(msgRes, msgReq)
	if msgRes.Code != http.StatusAccepted {
		t.Errorf("CreateMessage failure.  Expected %d to be %d", msgRes.Code, http.StatusAccepted)
	}

	// Channel blocks while waiting for a response from the testServer to verify
	// that the message was delivered
	result := <-resultChannel
	if !result {
		t.Error("Failed to receive message")
	}

}

func dropSubscriptionCollection() {
	session := DBSession.Copy()
	defer session.Close()
	c := session.DB(DBName).C(SubscriptionCollection)
	c.DropCollection()
}

func dropMessageCollection() {
	session := DBSession.Copy()
	defer session.Close()
	c := session.DB(DBName).C(MessageCollection)
	c.DropCollection()
}
