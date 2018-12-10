package struct2sql

import (
	"reflect"
	"strings"
)

//
//
//	Generate Query
//
//

func CreateTableQuery(table SQLTable) string {
	tableName := strings.ToLower(table.Name)
	query := "CREATE TABLE " + tableName + " ("
	detail := []string{}

	// Fields
	for _, f := range table.Fields {
		str := strings.ToLower(f.Name) + " "
		if f.CustomType != "" {
			str += f.CustomType
		} else {
			str += GoType2SQLType(f.Type)
		}
		if f.Keyword != "" {
			str += " " + f.Keyword
		}
		detail = append(detail, str)
	}

	// Keys
	if table.PrimaryKey != -1 {
		detail = append(detail, "PRIMARY KEY ("+strings.ToLower(table.Fields[table.PrimaryKey].Name)+")")
	}
	for _, f := range table.ForeignKey {
		str := "FOREIGN KEY (" + strings.ToLower(table.Fields[f.FieldIndex].Name) + ") REFERENCES " + f.Reference
		detail = append(detail, str)
	}

	// Combine the query
	query = query + strings.Join(detail, ", ")
	query += ");"
	return query
}

func DropTableQuery(tables ...SQLTable) string {
	tablesName := make([]string, len(tables))
	for i, table := range tables {
		tablesName[i] = strings.ToLower(table.Name)
	}
	query := "DROP TABLE " + strings.Join(tablesName, ", ") + ";"
	return query
}

func InsertQuery(model interface{}) string {
	t := reflect.Indirect(reflect.ValueOf(model)).Type()
	table := getTable(t)
	fields := []string{}
	questionMark := []string{}
	for i, f := range table.Fields {
		if i != table.PrimaryKey {
			fields = append(fields, strings.ToLower(f.Name))
			questionMark = append(questionMark, "?")
		}
	}
	query := "INSERT INTO " + strings.ToLower(table.Name) + "(" + strings.Join(fields, ", ") + ") VALUES(" + strings.Join(questionMark, ", ") + ");"
	return query
}

func SelectQuery(model interface{}, condition string) string {
	t := reflect.Indirect(reflect.ValueOf(model)).Type()
	var table SQLTable
	if t.Kind() == reflect.Slice {
		table = getTable(t.Elem())
	} else {
		table = getTable(t)
	}

	fields := []string{}
	for _, f := range table.Fields {
		fields = append(fields, strings.ToLower(f.Name))
	}
	return "SELECT " + strings.Join(fields, ", ") + " FROM " + strings.ToLower(table.Name) + " WHERE " + condition + ";"
}

func UpdateQuery(model interface{}, condition string) string {
	t := reflect.Indirect(reflect.ValueOf(model)).Type()
	table := getTable(t)
	fields := []string{}
	for i, f := range table.Fields {
		if i != table.PrimaryKey {
			fields = append(fields, strings.ToLower(f.Name)+"=?")
		}
	}
	return "UPDATE " + strings.ToLower(table.Name) + " SET " + strings.Join(fields, ", ") + " WHERE " + condition + ";"
}

func DeleteQuery(model interface{}, condition string) string {
	t := reflect.Indirect(reflect.ValueOf(model)).Type()
	table := getTable(t)
	return "DELETE FROM " + strings.ToLower(table.Name) + " WHERE " + condition + ";"
}

func Escape(str string) string {
	return strings.Replace(str, "'", "''", -1)
}
