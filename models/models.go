package models

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var db *mgo.Database

func Init(mdb *mgo.Database) {
	db = mdb

	c := DeviceLastPostCollection()

	// Unique Index
	index := mgo.Index{
		Key:        []string{"Imei", "line"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	if err := c.EnsureIndex(index); err != nil {
		panic(err)
	}
}

type GoodBad struct {
	Line      string `json:"line,omitempty"`
	Imei      string `json:"imei,omitempty"`
	Status    string `json:"status,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func GoodBadCollection() *mgo.Collection {
	return db.C("good_bad")
}

func GoodBadCount() int {
	count, err := GoodBadCollection().Count()
	if err != nil {
		panic(err)
	}

	return count
}

func (gb *GoodBad) Save() {
	gb.Timestamp = time.Now().Unix()

	if err := GoodBadCollection().Insert(gb); err != nil {
		panic(err)
	}
}

func Count(status string, line string) int {
	val, err := GoodBadCollection().Find(bson.M{"status": status, "line": line, "timestamp": bson.M{"$gt": time.Now().Unix() - 1800}}).Count()

	if err != nil {
		panic(err)
	}

	if val == 0 {
		return 0
	} else {
		return val
	}
}

func Decision(line string) string {
	var goods int = Count("good", line)
	var bads int = Count("bad", line)

	log.Println("Good: ", goods, " Bad: ", bads)

	if goods >= bads {
		return "good"
	} else {
		return "bad"
	}
}

type LineStatus struct {
	Line   string `json:"line,omitempty"`
	Goods  int    `json:"goods,omitempty"`
	Bads   int    `json:"bads,omitempty"`
	Status string `json:"status,omitempty"`
}

type DeviceLastPost struct {
	Imei      string
	Line      string
	Timestamp int64
}

func DeviceLastPostCollection() *mgo.Collection {
	return db.C("device_last_post")
}

func FindDeviceLastPost(imei, line string) *DeviceLastPost {
	dlp := DeviceLastPost{Imei: imei, Line: line}

	DeviceLastPostCollection().Find(bson.M{"imei": imei, "line": line}).One(&dlp)

	return &dlp
}

func (dlp *DeviceLastPost) Save() error {
	_, err := DeviceLastPostCollection().Upsert(bson.M{"imei": dlp.Imei, "line": dlp.Line}, dlp)

	return err
}
