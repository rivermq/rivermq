package model

import (
	"flag"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	// DBName provides the name of the database
	DBName = "rivermq"
	// MessageCollection provides the name of the collection which will hold
	// messages.
	MessageCollection = "messages"
	// SubscriptionCollection provides the name of the collection which will
	// hold subscriptions.
	SubscriptionCollection = "subscriptions"
)

var (
	dbAddress string
	once      sync.Once
	// DBSession is the initial session used to communicate with MongoDB.
	// Whenever a new session is needed, one should use the
	// mgo.Session.Copy() function.
	//
	// For example:
	//
	//  DBSession.Copy()
	//
	DBSession *mgo.Session
)

func init() {
	flag.StringVar(&dbAddress, "dbAddress", "localhost", "MongoDB address")
	once.Do(func() {
		session, err := mgo.Dial(dbAddress)
		if err != nil {
			panic(err)
		}
		session.SetMode(mgo.Monotonic, true)
		DBSession = session

		msgCollection := session.DB(DBName).C(MessageCollection)

		err = msgCollection.EnsureIndex(mgo.Index{
			Key:         []string{"timestamp"},
			ExpireAfter: 30 * time.Minute,
		})
		if err != nil {
			panic(err)
		}
	})
}
