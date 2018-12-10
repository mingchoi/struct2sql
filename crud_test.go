package struct2sql_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	s2s "fallen.world/mingchoi/struct2sql"
)

func TestCRUD(t *testing.T) {
	// Database Connection
	db, err := s2s.Open("mysql", "gomysql:123456@/gomysql?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	// Prepare table for testing
	prepareDBTable(db, t)

	// Test Insert Null model
	nullMonster := Monster{}
	err = db.Insert(&nullMonster)
	if err != nil {
		t.Error(err)
	}
	if nullMonster.ID != 1 {
		t.Error("Null model insert failed")
	}

	// Test Insert with foreign key
	err = db.Insert(&Map{})
	if err != nil {
		t.Error(err)
	}
	nowTime := time.Now()
	oldTime, err := time.Parse("2006-01-02", "1950-04-12")
	monster := Monster{
		Name:   "Monster 怪物",
		Level:  18,
		Health: 1029.8,
		Boss:   true,
		Spawn:  nowTime,
		Birth:  oldTime,
		MID:    1,
	}
	err = db.Insert(&monster)
	if err != nil {
		t.Error(err)
	}
	if monster.ID != 2 {
		t.Error(err)
	}

	// Test Select
	sm := Monster{}
	err = db.Select(&sm, "id = ?", 1)
	if err != nil {
		t.Error(err)
	}
	if sm.ID != 1 || sm.Name != "" || sm.Level != 0 || sm.Health != 0 || sm.Boss != false || sm.MID != 0 {
		t.Error("Select result incorrect: Null value for bool/int/float/string")
		t.Log(sm)
	}
	if (sm.Spawn != time.Time{} || sm.Birth != time.Time{}) {
		t.Error("Select result incorrect: Null value for time")
	}
	sm = Monster{}
	err = db.Select(&sm, "id = ?", 2)
	if err != nil {
		t.Error(err)
	}
	if sm.ID != 2 || sm.Name != "Monster 怪物" || sm.Level != 18 || sm.Health != 1029.8 || sm.Boss != true {
		t.Error("Select result incorrect: bool/int/float/string")
		t.Log(sm)
	}
	if !(sm.Spawn.After(nowTime.Add(-1*time.Second)) && sm.Spawn.Before(nowTime.Add(1*time.Second))) {
		t.Error("Select result incorrect: time")
		t.Log(sm.Spawn)
		t.Log(nowTime)
	}
	if sm.Birth.Year() != oldTime.Year() || sm.Birth.Month() != oldTime.Month() || sm.Birth.Day() != oldTime.Day() {
		t.Error("Select result incorrect: date")
		t.Log(sm.Birth)
		t.Log(oldTime)
	}
	if sm.MID != 1 {
		t.Error("Select result incorrect: foreign key")
		t.Log(sm)
	}

	// Select multiple
	sms := []Monster{}
	err = db.Select(&sms, "id > ?", 0)
	if len(sms) != 2 {
		t.Error("Selete multiple result failed: No result")
	}
	if sms[0].Level != 0 || sms[1].Name != "Monster 怪物" {
		t.Error("Select multiple result failed: Element incorrect")
	}

	// Update
	sm.Level = 21
	err = db.Update(&sm, "id = ?", sm.ID)
	if err != nil {
		t.Error(err)
	}
	um := Monster{}
	err = db.Select(&um, "id = ?", sm.ID)
	if err != nil {
		t.Error(err)
	}
	if um.Level != 21 {
		t.Error("Update row failed")
	}

	// Delete
	err = db.Delete(&Monster{}, "id = ?", 1)
	if err != nil {
		t.Error(err)
	}
	dm := Monster{}
	err = db.Select(&dm, "id = ?", 1)
	if err != nil {
		t.Error(err)
	}
	if dm.ID != 0 {
		t.Error("Delete row failed")
	}
}

func TestParsableTable(t *testing.T) {
	// Database Connection
	db, err := s2s.Open("mysql", "gomysql:123456@/gomysql?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	// Prepare table for testing
	prepareDBTable(db, t)

	// Test create row
	amount, err := decimal.NewFromString("678.90")
	pObj := Parsable{
		Telephone: "7777-7777",
		Money:     amount,
	}
	err = db.Insert(&pObj)
	if err != nil {
		t.Error(err)
	}
	if pObj.ID != 1 {
		t.Error("Insert failed")
	}

	// Test create row
	amount, err = decimal.NewFromString("999.02")
	pObj2 := Parsable{
		Telephone: "1234-5678",
		Money:     amount,
	}
	err = db.Insert(&pObj2)
	if err != nil {
		t.Error(err)
	}
	if pObj2.ID != 2 {
		t.Error("Insert failed")
	}

	// Read
	pObjs := []Parsable{}
	db.Select(&pObjs, "id > ?", 0)
	/*
		// Regular error handling
		err = db.Select(&pObjs, "id > 0")
		if err != nil {
			t.Error(err)
		}
	*/
	if len(pObjs) != 2 {
		t.Error("Seleted objects with incorrect amount")
	}
	if !reflect.DeepEqual(pObj, pObjs[0]) || !reflect.DeepEqual(pObj2, pObjs[1]) {
		t.Error("Seleted objects with incorrect info")
	}

	// Update
	amount, err = decimal.NewFromString("7777.77")
	pObj.Telephone = "8888-8888"
	pObj.Money = amount
	err = db.Update(&pObj, "id = ?", 1)
	if err != nil {
		t.Error(err)
	}
	updatedResult := Parsable{}
	db.Select(&updatedResult, "id = ?", 1)
	if updatedResult.ID != 1 {
		t.Error("Select object failed")
	}
	if !reflect.DeepEqual(pObj, updatedResult) {
		t.Error("Update object failed")
	}

	// Delete
	err = db.Delete(&Parsable{}, "id = ?", 2)
	if err != nil {
		t.Error(err)
	}
	deletedResult := Parsable{}
	db.Select(&deletedResult, "id = ?", 2)
	if err != nil {
		t.Error(err)
	}
	if deletedResult.ID != 0 {
		t.Error("Delete object failed")
	}
}

func prepareDBTable(db s2s.IDB, t *testing.T) {
	// Drop Table
	db.DropTable(&Monster{}, &Parsable{})
	db.DropTable(&Map{})

	// Create Table
	db.CreateTable(&Map{})
	db.CreateTable(&Monster{})
	db.CreateTable(&Parsable{})
}

func BenchmarkInsert(b *testing.B) {
	// Database Connection
	db, _ := s2s.Open("mysql", "gomysql:123456@/gomysql?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()

	// Prepare table for testing
	prepareDBTable(db, nil)

	for n := 0; n < b.N; n++ {
		db.Insert(&Monster{Name: "Benchmark Monster"})
	}
}

func BenchmarkInsertTx(b *testing.B) {
	// Database Connection
	db, _ := s2s.Open("mysql", "gomysql:123456@/gomysql?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()

	// Prepare table for testing
	prepareDBTable(db, nil)

	// Start Transcation
	tx, _ := db.Begin()
	for n := 0; n < b.N; n++ {
		tx.Insert(&Monster{Name: "Benchmark Monster"})
	}
	tx.Commit()
}
