package models

import (
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

func GoodBadCountCollection() *mgo.Collection {
	return db.C("good_bad")
}

func GoodBadCount() int {
	count, err := GoodBadCountCollection().Count()
	if err != nil {
		panic(err)
	}

	return count
}

func (gb *GoodBad) Save() {
	gb.Timestamp = time.Now().Unix()

	if err := GoodBadCountCollection().Insert(gb); err != nil {
		panic(err)
	}
}

type LineStatus struct {
	Line   string `json:"line,omitempty"`
	Goods  uint32 `json:"goods,omitempty"`
	Bads   uint32 `json:"bads,omitempty"`
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
