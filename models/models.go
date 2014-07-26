package models

import (
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

type GoodBad struct {
	Line      string `json:"line,omitempty"`
	Imei      string `json:"imei,omitempty"`
	Status    string `json:"status,omitempty"`
	Timestamp uint32 `json:"timestamp,omitempty"`
}

func (gb *GoodBad) Collection(db *mgo.Database) *mgo.Collection {
	return db.C("good_bad")
}

func (gb *GoodBad) Save(db *mgo.Database) {
	if err := gb.Collection(db).Insert(gb); err != nil {
		panic(err)
	}
}

type LineStatus struct {
	Line   string `json:"line,omitempty"`
	Goods  uint32 `json:"goods,omitempty"`
	Bads   uint32 `json:"bads,omitempty"`
	Status string `json:"status,omitempty"`
}
