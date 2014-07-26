package models

import (
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

type GoodBad struct {
	Line      string `json:"line,omitempty"`
	Imei      string `json:"imei,omitempty"`
	Status    string `json:"status,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func Collection(db *mgo.Database) *mgo.Collection {
	return db.C("good_bad")
}

func GoodBadCount(db *mgo.Database) int {
	count, err := Collection(db).Count()
	if err != nil {
		panic(err)
	}

	return count
}

func (gb *GoodBad) Save(db *mgo.Database) {
	if err := Collection(db).Insert(gb); err != nil {
		panic(err)
	}
}

type LineStatus struct {
	Line   string `json:"line,omitempty"`
	Goods  uint32 `json:"goods,omitempty"`
	Bads   uint32 `json:"bads,omitempty"`
	Status string `json:"status,omitempty"`
}
