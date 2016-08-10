package deliver

import (
	"bytes"
	"log"
	"net/http"

	"github.com/rivermq/rivermq/model"
)

// AttemptMessageDistribution passes the provided Message to the provided
// Subscription.CallbackURL.  If delivery fails, the message is saved to the DB
// for redelivery at a later time.
func AttemptMessageDistribution(msg model.Message, sub model.Subscription) (err error) {
	req, err := http.NewRequest("POST", sub.CallbackURL, bytes.NewBuffer(msg.Body))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		model.SaveFailedDeliveryMessage(msg, sub)
	} else {
		defer resp.Body.Close() // TODO: determine if this is necessary
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			err := model.SaveFailedDeliveryMessage(msg, sub)
			if err != nil {
				log.Printf("Error while saving Failed Delivery Message: %v", err)
			}
		}
	}
	return nil
}
