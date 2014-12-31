package models

import (
	"fmt"
	"time"

	"github.com/sfreiberg/mongo"
	"labix.org/v2/mgo/bson"
)

func Init(server, database string) {
	if err := mongo.SetServers(server, database); err != nil {
		panic(err)
	}

	// TODO add unique index
}

type GoodBad struct {
	Line      string `json:"line,omitempty"`
	Imei      string `json:"imei,omitempty"`
	Status    string `json:"status,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func (self *GoodBad) Save() error {
	self.Timestamp = time.Now().Unix()

	if GetDeviceLastPost(self.Imei, self.Line).Timestamp > time.Now().Unix()-300 {
		return fmt.Errorf("Rate Limit Error")
	}

	err := mongo.Insert(self)

	return err
}

type LineStatus struct {
	Line   string `json:"line,omitempty"`
	Goods  int    `json:"goods,omitempty"`
	Bads   int    `json:"bads,omitempty"`
	Status string `json:"status,omitempty"`
}

func GetLineStatus(line string) LineStatus {
	goodCount := count("good", line)
	badCount := count("bad", line)
	decision := decision(goodCount, badCount)

	return LineStatus{
		Line:   line,
		Goods:  goodCount,
		Bads:   badCount,
		Status: decision,
	}
}

func decision(goods, bads int) string {
	if goods >= bads {
		return "good"
	}

	return "bad"
}

func count(status string, line string) int {
	result := &[]GoodBad{}
	err := mongo.Find(result, bson.M{"status": status, "line": line, "timestamp": bson.M{"$gt": time.Now().Unix() - 1800}})

	if err != nil {
		panic(err)
	}

	return len(*result)
}

type DeviceLastPost struct {
	Imei      string
	Line      string
	Timestamp int64
}

func GetDeviceLastPost(imei, line string) *DeviceLastPost {
	lastGb := &GoodBad{}
	mongo.Find(lastGb, bson.M{"$query": bson.M{"imei": imei, "line": line}, "$orderby": bson.M{"timestamp": -1}})

	return &DeviceLastPost{Imei: imei, Line: line, Timestamp: lastGb.Timestamp}
}
