// +build integration

package deliver_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/rivermq/rivermq/inspect"
	. "github.com/rivermq/rivermq/model"
	. "github.com/rivermq/rivermq/route"
)

const (
	validMsg = `{"type":"msgType","body":{"state":"CO"}}`
)

var (
	// wg is used to indicate that all tests are complete
	wg sync.WaitGroup
)

func init() {
	wg.Add(2)
	go func() {
		wg.Wait()
		inspect.ShutdownPullSocket()
	}()
}

// TestAttemptMessageDistribution attempts to send a message to a test server
func TestAttemptMessageDistribution(t *testing.T) {
	defer dropMessageCollection()
	defer dropSubscriptionCollection()
	defer wg.Done()

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
	// that RiverMQ delivered the message to the testServer
	result := <-resultChannel
	if !result {
		t.Error("Failed to receive message")
	}

}

func TestAttemptMessageDistributionFailure(t *testing.T) {
	defer dropMessageCollection()
	defer dropSubscriptionCollection()
	defer wg.Done()

	// Create subscription in RiverMQ
	sub := Subscription{
		Type:        "msgType",
		CallbackURL: "http://localhost:3234234234",
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
		t.Errorf("CreateSubscription failure.  Expected %d to be %d\n", res.Code, http.StatusCreated)
	}

	// Send message to RiverMQ
	msgBuf := bytes.NewBufferString(validMsg)
	msgReq, _ := http.NewRequest("POST", "/messages", msgBuf)
	msgRes := httptest.NewRecorder()
	NewRiverMQRouter().ServeHTTP(msgRes, msgReq)
	if msgRes.Code != http.StatusAccepted {
		t.Errorf("CreateMessage failure.  Expected %d to be %d\n", msgRes.Code, http.StatusAccepted)
	}
	time.Sleep(500 * time.Millisecond)
	messages, err := FindMessageByStatus(StatusFailedDelivery)
	if err != nil {
		t.Errorf("%v\n", err)
	}
	if len(messages) != 1 {
		t.Errorf("FindMessageByStatus failure: Expected %d, found %d\n", 1, len(messages))
	}
	if messages[0].Status != StatusFailedDelivery {
		t.Errorf("FindMessageByStatus failure: Expected %s, found %s\n", StatusFailedDelivery, messages[0].Status)
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
	c.DropIndex("timestamp")
}
