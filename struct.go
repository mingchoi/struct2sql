package struct2sql

import (
	"reflect"
)

//
//
//	Struct
//
//

type SQLTable struct {
	Name       string
	NumField   int
	Fields     []SQLField
	PrimaryKey int
	ForeignKey []SQLForeignKey
}

type SQLField struct {
	Name         string
	Type         reflect.Type
	Keyword      string
	CustomType   string
	IsForeignKey bool
}

type SQLForeignKey struct {
	FieldIndex int
	Reference  string
}

//
//
//	Global Variable
//
//

var sqlTables map[reflect.Type]SQLTable = make(map[reflect.Type]SQLTable)

//
//
//	Analyze Struct
//
//

func AnalyzeStruct(model interface{}) SQLTable {
	// Get type of the struct
	var v reflect.Type
	var ok bool
	if v, ok = model.(reflect.Type); !ok {
		v = reflect.Indirect(reflect.ValueOf(model)).Type()
	}

	// Return if table already exist
	if t, ok := sqlTables[v]; ok {
		return t
	}

	// Get basic info from the struct
	sqlTable := SQLTable{
		Name:       v.Name(),
		NumField:   v.NumField(),
		Fields:     make([]SQLField, v.NumField()),
		PrimaryKey: -1,
		ForeignKey: make([]SQLForeignKey, 0),
	}
	if v.Field(0).Tag.Get("tablename") != "" {
		sqlTable.Name = v.Field(0).Tag.Get("tablename")
	}
	// Get struct field info
	for i := 0; i < sqlTable.NumField; i++ {
		sqlTable.Fields[i] = SQLField{
			Name:       v.Field(i).Name,
			Type:       v.Field(i).Type,
			Keyword:    v.Field(i).Tag.Get("keyword"),
			CustomType: v.Field(i).Tag.Get("customtype"),
		}

		// Detect tags like primary key and foreign key
		if v.Field(i).Tag.Get("primarykey") == "true" {
			sqlTable.PrimaryKey = i
		} else if v.Field(i).Tag.Get("foreignkey") != "" {
			sqlTable.ForeignKey = append(sqlTable.ForeignKey, SQLForeignKey{
				FieldIndex: i,
				Reference:  v.Field(i).Tag.Get("foreignkey"),
			})
			sqlTable.Fields[i].IsForeignKey = true
		}
	}
	sqlTables[v] = sqlTable
	return sqlTable
}

func getTable(t reflect.Type) SQLTable {
	if sqlTable, ok := sqlTables[t]; ok {
		return sqlTable
	}
	return AnalyzeStruct(t)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
