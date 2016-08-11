package deliver

import (
	"bytes"
	"encoding/gob"
	"log"
	"net/http"

	"github.com/rivermq/rivermq/model"
	"github.com/zeromq/goczmq"
)

var (
	inspectPullSocket *goczmq.Sock
	decoder           *gob.Decoder
	shutdownChannel   chan bool
)

func init() {
	inspectPullSocket, err := goczmq.NewPull("inproc://inspect")
	if err != nil {
		panic(err)
	}
	decoder = gob.NewDecoder(inspectPullSocket)
	shutdownChannel = make(chan bool, 1)

	go func() {
		shutdownMsg := <-shutdownChannel
		if shutdownMsg {
			inspectPullSocket.Destroy()
		}
	}()

	go func() {
		for {
			var msg model.Message
			decoder.Decode(&msg)
			AttemptMessageDistribution(msg)
		}
	}()
}

// ShutdownPullSocket provides a method to destory the blocked inspectPullSocket
func ShutdownPullSocket() {
	shutdownChannel <- true
}

// AttemptMessageDistribution passes the provided Message to the provided
// Subscription.CallbackURL.  If delivery fails, the message is saved to the DB
// for redelivery at a later time.
func AttemptMessageDistribution(msg model.Message) (err error) {
	req, err := http.NewRequest("POST", msg.RetryEndpoint, bytes.NewBuffer(msg.Body))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		model.SaveFailedDeliveryMessage(msg)
	} else {
		defer resp.Body.Close() // TODO: determine if this is necessary
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			err := model.SaveFailedDeliveryMessage(msg)
			if err != nil {
				log.Printf("Error while saving Failed Delivery Message: %v", err)
			}
		}
	}
	return nil
}
