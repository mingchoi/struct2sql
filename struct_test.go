package struct2sql_test

import (
	"reflect"
	"testing"
	"time"

	s2s "fallen.world/mingchoi/struct2sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
)

type EmptyStruct struct {
	ID int
}

type Map struct {
	ID int `keyword:"NOT NULL AUTO_INCREMENT" primarykey:"true"`
}

type Monster struct {
	ID     int `keyword:"NOT NULL AUTO_INCREMENT" primarykey:"true"`
	Name   string
	Level  int
	Health float64
	Boss   bool
	Spawn  time.Time
	Birth  time.Time `customtype:"date"`
	MID    int       `foreignkey:"map(id)"`
}

type Parsable struct {
	ID        int `keyword:"NOT NULL AUTO_INCREMENT" primarykey:"true"`
	Telephone Telephone
	Money     decimal.Decimal `customtype:"decimal(10,2)"`
}

func TestAnalyzeStruct(t *testing.T) {
	// Analyze a type
	tSqlTable := s2s.AnalyzeStruct(reflect.TypeOf(EmptyStruct{}))
	if tSqlTable.Name != "EmptyStruct" {
		t.Error("Table name incorrect")
	}
	if tSqlTable.NumField != 1 || len(tSqlTable.Fields) != 1 {
		t.Error("Field info incorrect")
	}
	if f := tSqlTable.Fields[0]; f.Name != "ID" || f.Type != reflect.TypeOf(0) || f.Keyword != "" || f.IsForeignKey != false {
		t.Error("Field info incorrect")
	}
	if tSqlTable.PrimaryKey != -1 {
		t.Error("Primary key wrong index")
	}
	if len(tSqlTable.ForeignKey) != 0 {
		t.Error("Foreign key incorrect")
	}
	// Analyze a struct
	sqlTable := s2s.AnalyzeStruct(&Monster{})
	if sqlTable.Name != "Monster" {
		t.Error("Table name incorrect")
	}
	if sqlTable.NumField != 8 || len(sqlTable.Fields) != 8 {
		t.Error("Number of fields incorrect")
	}
	if f := sqlTable.Fields[0]; f.Name != "ID" || f.Type != reflect.TypeOf(0) || f.Keyword != "NOT NULL AUTO_INCREMENT" {
		t.Error("Field info incorrect")
	}
	if sqlTable.Fields[6].CustomType != "date" {
		t.Error("Custom type info incorrect")
	}
	if sqlTable.PrimaryKey != 0 {
		t.Error("Primary key wrong index")
	}
	if sqlTable.Fields[7].IsForeignKey != true {
		t.Error("Foreign key tag missing")
	}
	if len(sqlTable.ForeignKey) != 1 {
		t.Error("Foreign key missing")
	}
	if fk := sqlTable.ForeignKey[0]; fk.FieldIndex != 7 || fk.Reference != "map(id)" {
		t.Error("Foreign key info incorrect")
	}

}
