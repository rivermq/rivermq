package model

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Subscription provides the data required to create a Subscription for clients
type Subscription struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Timestamp   int64  `json:"timestamp"`
	Type        string `json:"type"`
	CallbackURL string `json:"callbackUrl"`
}

// ValidateSubscription performs sanity checks on a Subscription instance
func ValidateSubscription(sub Subscription) (err error) {
	if typeLen := len(sub.Type); typeLen == 0 {
		log.Printf("Error parsing supplied subscription Type[ %v ]\n ", sub.Type)
		return fmt.Errorf("Error parsing supplied subscription Type[ %v ]\n ", sub.Type)
	}
	_, err = url.ParseRequestURI(sub.CallbackURL)
	if err != nil {
		return err
	}
	return nil
}

// SaveSubscription saves the supplied subscription to MongoDB
func SaveSubscription(sub Subscription) (resultSub Subscription, err error) {
	session := DBSession.Copy()
	defer session.Close()
	c := session.DB(DBName).C(SubscriptionCollection)

	sub.Timestamp = time.Now().UnixNano()
	err = c.Insert(sub)
	if err != nil {
		log.Fatal(err)
	}
	return sub, err
}

// FindSubscriptionByID finds a given subscription by its id
func FindSubscriptionByID(id string) (sub Subscription, err error) {
	session := DBSession.Copy()
	defer session.Close()
	c := session.DB(DBName).C(SubscriptionCollection)

	err = c.FindId(bson.ObjectIdHex(id)).One(&sub)
	if err != nil {
		return sub, err
	}
	return sub, nil
}
