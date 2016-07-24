package model

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/pborman/uuid"
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

	sub.ID = uuid.NewUUID().String()
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

	err = c.FindId(id).One(&sub)
	if err != nil {
		return sub, err
	}
	return sub, nil
}

// FindAllSubscriptionsByType returns a slice of subscriptions having a given type
func FindAllSubscriptionsByType(subType string) (subs []Subscription, err error) {
	session := DBSession.Copy()
	defer session.Close()
	c := session.DB(DBName).C(SubscriptionCollection)

	err = c.Find(bson.M{"type": subType}).All(&subs)
	return subs, err
}

// FindAllSubscriptions returns a slice of all subscriptions
func FindAllSubscriptions() (subs []Subscription, err error) {
	session := DBSession.Copy()
	defer session.Close()
	c := session.DB(DBName).C(SubscriptionCollection)

	err = c.Find(nil).All(&subs)
	return subs, err
}

// DeleteSubscriptionByID deletes a given subscription having a given id
func DeleteSubscriptionByID(subID string) (err error) {
	session := DBSession.Copy()
	defer session.Close()
	c := session.DB(DBName).C(SubscriptionCollection)

	err = c.RemoveId(subID)
	return err
}
