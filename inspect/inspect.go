package inspect

import (
	"fmt"

	"github.com/rivermq/rivermq/deliver"
	"github.com/rivermq/rivermq/model"
)

// HandleMessage checks for Subscriptions to the supplied Message and
// either drops it, if no Subscriptions are found, or passes it on
// for delivery
func HandleMessage(msg model.Message) (err error) {
	fmt.Println("handle message")
	subs, err := model.FindAllSubscriptionsByType(msg.Type)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	if len(subs) > 0 {
		for _, sub := range subs {
			deliver.AttemptMessageDistrobution(msg, sub)
		}
	}
	return nil
}
