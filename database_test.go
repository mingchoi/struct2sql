package struct2sql_test

import (
	"testing"

	s2s "fallen.world/mingchoi/struct2sql"
)

func TestTranscation(t *testing.T) {
	// Database Connection
	db, err := s2s.Open("mysql", "gomysql:123456@/gomysql?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	// Prepare table for testing
	prepareDBTable(db, t)

	// Start Transcation
	tx, _ := db.Begin()

	// Test Insert Null model
	nullMonster := Monster{Name: "Test"}
	err = tx.Insert(&nullMonster)
	tx.Commit()
	if err != nil {
		t.Error(err)
	}
	if nullMonster.ID != 1 {
		t.Error("Null model insert failed")
	}
}

func TestPassbyInterface(t *testing.T) {
	// Database Connection
	db, err := s2s.Open("mysql", "gomysql:123456@/gomysql?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	// Prepare table for testing
	prepareDBTable(db, t)

	// Test Pass by Interface
	m1 := Monster{Name: "Pass By Interface(DB): OK"}
	passByInterface(db, &m1)
	m2 := Monster{Name: "Pass By Interface(Tx): OK"}
	tx, _ := db.Begin()
	passByInterface(tx, &m2)
	tx.Commit()
	if m1.ID != 1 {
		t.Error("Pass by Interface(DB) Failed")
	}
	if m2.ID != 2 {
		t.Error("Pass by Interface(Tx) Failed")
	}
}

func passByInterface(db s2s.IDB, m *Monster) {
	db.Insert(m)
}
