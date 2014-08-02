package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Init(db *mgo.Database) {
	c := DeviceLastPostCollection(db)

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

func GoodBadCountCollection(db *mgo.Database) *mgo.Collection {
	return db.C("good_bad")
}

func GoodBadCount(db *mgo.Database) int {
	count, err := GoodBadCountCollection(db).Count()
	if err != nil {
		panic(err)
	}

	return count
}

func (gb *GoodBad) Save(db *mgo.Database) {
	gb.Timestamp = time.Now().Unix()

	if err := GoodBadCountCollection(db).Insert(gb); err != nil {
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

func DeviceLastPostCollection(db *mgo.Database) *mgo.Collection {
	return db.C("device_last_post")
}

func FindDeviceLastPost(db *mgo.Database, imei, line string) *DeviceLastPost {
	dlp := DeviceLastPost{Imei: imei, Line: line}

	DeviceLastPostCollection(db).Find(bson.M{"imei": imei, "line": line}).One(&dlp)

	return &dlp
}

func (dlp *DeviceLastPost) Save(db *mgo.Database) error {
	_, err := DeviceLastPostCollection(db).Upsert(bson.M{"imei": dlp.Imei, "line": dlp.Line}, dlp)

	return err
}
