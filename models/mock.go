package models

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

type Mock struct {
	Id          bson.ObjectId     `json:"id" bson:"_id"`
	Method      string            `json:"method"`
	Endpoint    string            `json:"endpoint"`
	QueryParams string            `json:"query_params"`
	Headers     map[string]string `json:"headers"`
	Payload     *json.RawMessage  `json:"payload"`
}
