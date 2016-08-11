package model

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

const (
	// StatusAccepted indicateds that a Message has been accepted, but has not
	// been processed
	StatusAccepted = "ACCEPTED"

	// StatusInspected indicates that a Message has been inspected and there are
	// are current Subscriptions expecting to receive the message
	StatusInspected = "INSPECTED"

	// StatusFailedDelivery indicates that RiverMQ was unable to deliver the
	// message.  When this occurs the 'RetryEndpoint' property is set, a new ID
	// is generated and the Message is saved to MongoDB.  At a later time RiverMQ
	// will attempt to redeliver the message.
	StatusFailedDelivery = "FAILED_DELIVERY"

	// StatusDeliveryRetrySuccess indicates that the message has been successfully
	// delivered on retry and no further actions are required.
	StatusDeliveryRetrySuccess = "DELIVERY_RETRY_SUCCESS"
)

// Message provides a message to be distributed by RiverMQ
type Message struct {
	ID            string          `json:"id" bson:"_id,omitempty"`
	Timestamp     time.Time       `json:"timestamp"`
	Type          string          `json:"type"`
	Status        string          `json:"status"`
	RetryEndpoint string          `json:"retryEndpoint,omitempty"`
	Body          json.RawMessage `json:"body"`
}

// ValidateMessage ensures a message is valid
func ValidateMessage(msg Message) (err error) {
	if len(msg.Type) == 0 {
		log.Printf("Error parsing supplied message Type[ %v ]\n ", msg.Type)
		return fmt.Errorf("Error parsing supplied message Type[ %v ]\n ", msg.Type)
	}
	return nil
}

// SaveMessage saves the supplied message to MongoDB
func SaveMessage(msg Message) (resultMsg Message, err error) {
	session := DBSession.Copy()
	defer session.Close()
	c := session.DB(DBName).C(MessageCollection)

	msg.ID = uuid.NewUUID().String()
	msg.Timestamp = bson.Now()
	err = c.Insert(msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg, err
}

// SaveFailedDeliveryMessage prepares a message for redelivery and saves it
// to the DB
func SaveFailedDeliveryMessage(msg Message) (err error) {
	session := DBSession.Copy()
	defer session.Close()
	c := session.DB(DBName).C(MessageCollection)

	// prepare message for retry delivery
	msg.ID = uuid.NewUUID().String()
	msg.Status = StatusFailedDelivery
	msg.Timestamp = bson.Now()

	err = c.Insert(msg)
	return err
}

// UpdateMessage updates a given message in the db
func UpdateMessage(msg Message) (err error) {
	session := DBSession.Copy()
	defer session.Close()
	c := session.DB(DBName).C(MessageCollection)

	err = c.UpdateId(msg.ID, msg)
	return err
}

// FindMessageByID finds a given message by its id
func FindMessageByID(id string) (msg Message, err error) {
	session := DBSession.Copy()
	defer session.Close()
	c := session.DB(DBName).C(MessageCollection)

	err = c.FindId(id).One(&msg)
	if err != nil {
		return msg, err
	}
	return msg, nil
}

// FindMessageByStatus finds all messages which contain the supplied status
func FindMessageByStatus(status string) (messages []Message, err error) {
	session := DBSession.Copy()
	defer session.Close()
	c := session.DB(DBName).C(MessageCollection)

	err = c.Find(bson.M{"status": status}).All(&messages)
	return messages, err

}
