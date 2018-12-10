package struct2sql_test

import (
	"fmt"
	"testing"

	s2s "fallen.world/mingchoi/struct2sql"
)

func TestCreateTableQuery(t *testing.T) {
	sqlTable := s2s.AnalyzeStruct(&Monster{})
	expected := "CREATE TABLE monster (id int NOT NULL AUTO_INCREMENT, name varchar(255), level int, health double, boss bool, spawn datetime, birth date, mid int, PRIMARY KEY (id), FOREIGN KEY (mid) REFERENCES map(id));"
	if s2s.CreateTableQuery(sqlTable) != expected {
		fmt.Println(s2s.CreateTableQuery(sqlTable))
		t.Error("Query for create table incorrect")
	}
}

func TestInsertQuery(t *testing.T) {
	model := Monster{}
	expected := "INSERT INTO monster(name, level, health, boss, spawn, birth, mid) VALUES(?, ?, ?, ?, ?, ?, ?);"
	if s2s.InsertQuery(model) != expected {
		t.Error("Query for insert row incorrect")
	}
}

func TestSelectQuery(t *testing.T) {
	model := Monster{}
	expected := "SELECT id, name, level, health, boss, spawn, birth, mid FROM monster WHERE id = 0;"
	if s2s.SelectQuery(model, "id = 0") != expected {
		t.Error("Query for insert row incorrect")
	}
}

func TestUpdateQuery(t *testing.T) {
	model := Monster{}
	expected := "UPDATE monster SET name=?, level=?, health=?, boss=?, spawn=?, birth=?, mid=? WHERE id = 0;"

	if s2s.UpdateQuery(model, "id = 0") != expected {
		t.Error("Query for insert row incorrect")
	}
}

func TestDeleteQuery(t *testing.T) {
	model := Monster{}
	expected := "DELETE FROM monster WHERE id = 0;"
	if s2s.DeleteQuery(model, "id = 0") != expected {
		t.Error("Query for insert row incorrect")
	}
}
