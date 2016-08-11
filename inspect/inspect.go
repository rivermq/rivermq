package inspect

import (
	"encoding/gob"
	"log"

	"github.com/rivermq/rivermq/model"
	"github.com/zeromq/goczmq"
)

var (
	handlerPullSocket *goczmq.Sock
	inspectPushSocket *goczmq.Sock
	decoder           *gob.Decoder
	encoder           *gob.Encoder
	shutdownChannel   chan bool
)

func init() {
	handlerPullSocket, err := goczmq.NewPull("inproc://handler")
	inspectPushSocket, err := goczmq.NewPush("inproc://inspect")
	if err != nil {
		panic(err)
	}
	decoder = gob.NewDecoder(handlerPullSocket)
	encoder = gob.NewEncoder(inspectPushSocket)
	shutdownChannel = make(chan bool, 1)

	go func() {
		shutdownMsg := <-shutdownChannel
		if shutdownMsg {
			handlerPullSocket.Destroy()
		}
	}()

	go func() {
		for {
			var msg model.Message
			decoder.Decode(&msg)
			HandleMessage(msg)
		}
	}()

}

// ShutdownPullSocket provides a method to shutdown the blocked inspectPullSocket
func ShutdownPullSocket() {
	shutdownChannel <- true
}

// HandleMessage checks for Subscriptions to the supplied Message and
// either drops it, if no Subscriptions are found, or passes it on
// for delivery
func HandleMessage(msg model.Message) (err error) {
	subs, err := model.FindAllSubscriptionsByType(msg.Type)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return err
	}
	if len(subs) > 0 {
		msg.Status = model.StatusInspected
		model.UpdateMessage(msg)
		for _, sub := range subs {
			msg.RetryEndpoint = sub.CallbackURL
			go func(msg model.Message) {
				encoder.Encode(msg)
			}(msg)
		}
	}

	return nil
}
