package struct2sql_test

import (
	"database/sql/driver"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	s2s "fallen.world/mingchoi/struct2sql"
)

type Telephone string

func (t Telephone) Type() string {
	return "int"
}

func (t *Telephone) Scan(src interface{}) error {
	var str string
	switch src.(type) {
	case []uint8:
		str = string(src.([]uint8))
	case string:
		str = src.(string)
	case int64:
		str = strconv.FormatInt(src.(int64), 10)
	default:
		panic("Unable to scan unknown type")
	}
	*t = Telephone(str[:4] + "-" + str[4:])
	return nil
}

func (t Telephone) Value() (driver.Value, error) {
	str := string(t[:4] + t[5:])
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func TestGoType2SQLType(t *testing.T) {
	sqlTable := s2s.AnalyzeStruct(&Monster{})
	if s2s.GoType2SQLType(sqlTable.Fields[4].Type) != "bool" {
		t.Error("Type convert error: bool")
	}
	if s2s.GoType2SQLType(sqlTable.Fields[0].Type) != "int" {
		t.Error("Type convert error: int")
	}
	if s2s.GoType2SQLType(sqlTable.Fields[3].Type) != "double" {
		t.Error("Type convert error: float")
	}
	if s2s.GoType2SQLType(sqlTable.Fields[1].Type) != "varchar(255)" {
		t.Error("Type convert error: string")
	}
	if s2s.GoType2SQLType(sqlTable.Fields[5].Type) != "datetime" {
		t.Error("Type convert error: time")
	}
	if s2s.GoType2SQLType(reflect.TypeOf(sqlTable)) != "TYPE_NOT_SUPPORTED" {
		t.Error("Wait, what?")
	}

	// IParsableType
	telType := reflect.TypeOf((Telephone)(""))
	if s2s.GoType2SQLType(telType) != "int" {
		t.Error("Type convert error: IParsableType")
	}

}

func TestConvertSQLResult(t *testing.T) {
	// Ready variable
	monster := Monster{}
	s2s.AnalyzeStruct(&monster)
	nowTime := time.Now()
	oldTime, _ := time.Parse("2006-01-02", "1950-04-12")

	// Test convert
	input := make([]interface{}, 8)
	input[0] = int64(1)
	input[1] = "Monster 怪物"
	input[2] = int64(18)
	input[3] = 1029.8
	input[4] = true
	input[5] = nowTime
	input[6] = oldTime
	input[7] = int64(1)
	result := s2s.ConvertSQLResult(&monster, input)
	for i := range input {
		if input[i] != result[i] {
			t.Error("Element convert failed: ", input[i], "->", result[i])
		}
	}

	// Test convert From raw data
	input = make([]interface{}, 8)
	input[0] = []uint8("1")
	input[1] = []uint8("Monster 怪物")
	input[2] = []uint8("18")
	input[3] = []uint8("1029.8")
	input[4] = []uint8("TRUE")
	input[5] = time.Time{}
	input[6] = time.Time{}
	input[7] = []uint8("1")
	result = s2s.ConvertSQLResult(&monster, input)
	if result[4] != true {
		t.Error("Element convert failed: bool")
	}
	if result[0] != int64(1) || result[2] != int64(18) || result[7] != int64(1) {
		t.Error("Element convert failed: int")
	}
	if result[3] != 1029.8 {
		t.Error("Element convert failed: float")
	}
	if result[1] != "Monster 怪物" {
		t.Error("Element convert failed: string")
	}

	// Test converting IParsableType
	p := Parsable{}
	s2s.AnalyzeStruct(&p)
	input = make([]interface{}, 3)
	input[0] = 99
	input[1] = int64(12345678)
	input[2] = "10080.72"
	result = s2s.ConvertSQLResult(&p, input)
	t.Log(result)
	if result[0] != int64(99) {
		t.Error("Converting SQL result error: int64")
	}
	if result[1] != Telephone("1234-5678") {
		t.Error("Converting SQL result error: Custom type (Telephone)")
	}
	if !result[2].(decimal.Decimal).Equal(decimal.NewFromFloat(10080.72)) {
		t.Error("Converting SQL result error: Custom type (Decimal)")
	}

}
