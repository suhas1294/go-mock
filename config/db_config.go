package config

import (
	"gopkg.in/mgo.v2"
	"log"
)

var MongoSession *mgo.Session
var dbErr error

func init() {
	log.Println("Initialising database..")
	// Connect to our local mongo
	MongoSession, dbErr = mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if dbErr != nil {
		panic(dbErr)
	}
}
