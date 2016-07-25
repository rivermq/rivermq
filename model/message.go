package model

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pborman/uuid"
)

// Message provides a message to be distributed by RiverMQ
type Message struct {
	ID        string          `json:"id "bson:"_id,omitempty"`
	Timestamp int64           `json:"timestamp"`
	Type      string          `json:"type"`
	Status    string          `json:"status"`
	Body      json.RawMessage `json:"body"`
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
	msg.Timestamp = time.Now().UnixNano()
	err = c.Insert(msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg, err
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
