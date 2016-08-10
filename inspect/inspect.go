package inspect

import (
	"encoding/gob"
	"log"
	"os"

	"github.com/rivermq/rivermq/deliver"
	"github.com/rivermq/rivermq/model"
	"github.com/zeromq/goczmq"
)

var (
	pullSocket      *goczmq.Sock
	decoder         *gob.Decoder
	shutdownChannel chan bool
)

func init() {
	log.SetOutput(os.Stderr)
	pullSocket, err := goczmq.NewPull("inproc://handler")
	if err != nil {
		panic(err)
	}
	decoder = gob.NewDecoder(pullSocket)
	shutdownChannel = make(chan bool, 1)

	go func() {
		shutdownMsg := <-shutdownChannel
		if shutdownMsg {
			pullSocket.Destroy()
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

// ShutdownPullSocket provides a method to shutdown the blocked subSocket
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
			go func(msg model.Message, sub model.Subscription) {
				deliver.AttemptMessageDistribution(msg, sub)
			}(msg, sub)
		}
	}

	return nil
}
