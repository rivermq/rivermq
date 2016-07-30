package deliver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rivermq/rivermq/model"
	"net/http"
)

// AttemptMessageDistrobution passes the provided Message to the provided Subscription
func AttemptMessageDistrobution(msg model.Message, sub model.Subscription) (err error) {
	fmt.Println("DeliverMessage")
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	req, err := http.NewRequest("POST", sub.CallbackURL, bytes.NewBuffer(msgJSON))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	return nil
}
