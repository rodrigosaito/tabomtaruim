package models

import (
	"testing"

	"github.com/sfreiberg/mongo"
	"github.com/stretchr/testify/assert"
)

func init() {
	dbName := "good_bad_test_models"

	Init("localhost", dbName)

	session, _ := mongo.GetSession()
	session.DB(dbName).DropDatabase()
}

func TestGoodBadSave(t *testing.T) {
	goodBad := &GoodBad{
		Line:   "cptm-9",
		Imei:   "123321",
		Status: "good",
	}

	if err := goodBad.Save(); err != nil {
		t.Error("Failed to save GoodBad")
	}

	if goodBad.Timestamp == 0 {
		t.Error("Timestamp should be set before saving")
	}
}

func TestGoodBadSaveTwice(t *testing.T) {
	gb1 := &GoodBad{Line: "metro-1", Imei: "123321", Status: "good"}

	if err := gb1.Save(); err != nil {
		t.Error("Failed to save GoodBad")
	}

	gb2 := &GoodBad{Line: "metro-1", Imei: "123321", Status: "bad"}

	err := gb2.Save()
	assert.EqualError(t, err, "Rate Limit Error")
}

func TestGetLineStatus(t *testing.T) {
	// Insert test data
	gb1 := &GoodBad{Line: "metro-4", Imei: "1234", Status: "good"}
	gb2 := &GoodBad{Line: "metro-4", Imei: "1235", Status: "bad"}
	gb3 := &GoodBad{Line: "metro-4", Imei: "1236", Status: "good"}

	gb1.Save()
	gb2.Save()
	gb3.Save()

	// Test
	status := GetLineStatus("metro-4")

	assert.Equal(t, "metro-4", status.Line, "should return the correct line status")
	assert.Equal(t, 2, status.Goods)
	assert.Equal(t, 1, status.Bads)
	assert.Equal(t, "good", status.Status)
}

func TestDecision(t *testing.T) {
	assert.Equal(t, "good", decision(10, 1), "should be 'good' when goods are greater")
	assert.Equal(t, "bad", decision(1, 10), "should be 'bad' when goods are smaller")
	assert.Equal(t, "good", decision(10, 10), "should be 'good' when goods are equal bads")
}

func TestGetDeviceLastPost(t *testing.T) {
	// When first post
	dlp1 := GetDeviceLastPost("1234", "cptm-8")

	assert.Equal(t, 0, dlp1.Timestamp)

	// When a previous post has been sent
	gb := &GoodBad{Line: "cptm-8", Imei: "1234", Status: "good"}
	gb.Save()

	dlp2 := GetDeviceLastPost("1234", "cptm-8")

	assert.NotNil(t, dlp2.Timestamp)
}

func TestGoodCountBefore30Minutes(t *testing.T) {
	t.SkipNow()
}
