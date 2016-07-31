package inspect

import (
	"log"

	"github.com/rivermq/rivermq/deliver"
	"github.com/rivermq/rivermq/model"
)

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
				deliver.AttemptMessageDistrobution(msg, sub)
			}(msg, sub)
		}
	}

	return nil
}
