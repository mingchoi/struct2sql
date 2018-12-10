package struct2sql

import (
	"database/sql"
	"errors"
	"reflect"
	"time"
)

func (db DB) Insert(model interface{}) error {
	return create(db, model)
}
func (db DB) Select(model interface{}, condition string, cargs ...interface{}) error {
	return read(db, model, condition, cargs...)
}
func (db DB) Update(model interface{}, condition string, cargs ...interface{}) error {
	return update(db, model, condition, cargs...)
}
func (db DB) Delete(model interface{}, condition string, cargs ...interface{}) error {
	return delete(db, model, condition, cargs...)
}
func (db DB) View(model interface{}, query string, cargs ...interface{}) error {
	return view(db, model, query, cargs...)
}

func (db Tx) Insert(model interface{}) error {
	return create(db, model)
}
func (db Tx) Select(model interface{}, condition string, cargs ...interface{}) error {
	return read(db, model, condition, cargs...)
}
func (db Tx) Update(model interface{}, condition string, cargs ...interface{}) error {
	return update(db, model, condition, cargs...)
}
func (db Tx) Delete(model interface{}, condition string, cargs ...interface{}) error {
	return delete(db, model, condition, cargs...)
}
func (db Tx) View(model interface{}, query string, cargs ...interface{}) error {
	return view(db, model, query, cargs...)
}

//
//
//	Create, Read, Update, Delete in Database
//
//

func create(db IDB, model interface{}) error {
	// Check pointer
	if reflect.ValueOf(model).Kind() != reflect.Ptr {
		return errors.New("Only pointer model excepted")
	}

	// Generate Query
	query := InsertQuery(model)

	// Extract variable
	v := reflect.Indirect(reflect.ValueOf(model))
	table := getTable(v.Type())

	args := make([]interface{}, len(table.Fields))
	for i, f := range table.Fields {
		if i != table.PrimaryKey {
			switch f.Type {
			case reflect.TypeOf(true):
				args[i] = v.Field(i).Bool()
			case reflect.TypeOf(0):
				args[i] = v.Field(i).Int()
				if table.Fields[i].IsForeignKey && v.Field(i).Int() == 0 {
					args[i] = nil
				}
			case reflect.TypeOf(0.0):
				args[i] = v.Field(i).Float()
			case reflect.TypeOf(""):
				args[i] = v.Field(i).String()
			case reflect.TypeOf(time.Time{}):
				args[i] = v.Field(i).Interface()
				if (args[i] == time.Time{}) {
					args[i] = sql.NullInt64{}
				}
			default:
				args[i] = v.Field(i).Interface()
			}
		}
	}
	if table.PrimaryKey != -1 {
		args = append(args[:table.PrimaryKey], args[table.PrimaryKey+1:]...)
	}

	// Check query
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Insert Data
	result, err := stmt.Exec(args...)
	if err != nil {
		return err
	}

	// Return Primary Key
	id, err := result.LastInsertId()
	check(err)
	if table.PrimaryKey != -1 {
		v.Field(table.PrimaryKey).SetInt(id)
	}

	return nil
}

func read(db IDB, model interface{}, condition string, cargs ...interface{}) error {
	// Check pointer
	if reflect.ValueOf(model).Kind() != reflect.Ptr {
		return errors.New("Only pointer model excepted")
	}

	// Generate Query
	query := SelectQuery(model, condition)

	// Get table
	m := reflect.Indirect(reflect.ValueOf(model))
	var table SQLTable
	if m.Kind() == reflect.Slice {
		table = getTable(m.Type().Elem())
	} else {
		table = getTable(m.Type())
	}

	// Detect Struct or Slice
	var sliceMode bool
	switch m.Kind() {
	case reflect.Slice:
		sliceMode = true
	case reflect.Struct:
		sliceMode = false
	}

	// Read result
	rows, err := db.Query(query, cargs...)
	check(err)
	defer rows.Close()
	for rows.Next() {
		store := make([]interface{}, len(table.Fields))
		args := make([]interface{}, len(table.Fields))
		for i := range args {
			args[i] = &store[i]
		}
		if err := rows.Scan(args...); err != nil {
			check(err)
		}

		var obj reflect.Value
		if sliceMode {
			obj = reflect.New(m.Type().Elem()).Elem()
		} else {
			obj = m
		}

		// Convert result & reflect to model
		converted := ConvertSQLResult(model, store)
		for i, f := range table.Fields {
			switch f.Type {
			case reflect.TypeOf(false):
				obj.Field(i).SetBool(converted[i].(bool))
			case reflect.TypeOf(0):
				obj.Field(i).SetInt(converted[i].(int64))
			case reflect.TypeOf(0.0):
				obj.Field(i).SetFloat(converted[i].(float64))
			case reflect.TypeOf(""):
				obj.Field(i).SetString(converted[i].(string))
			case reflect.TypeOf(time.Time{}):
				if converted[i] != nil {
					obj.Field(i).Set(reflect.ValueOf(converted[i]))
				}
			default:
				if converted[i] != nil {
					obj.Field(i).Set(reflect.ValueOf(converted[i]))
				}
			}

		}

		// Detect Struct or Slice
		if sliceMode {
			m.Set(reflect.Append(m, obj))
		} else {
			break
		}
	}
	return nil
}

func update(db IDB, model interface{}, condition string, cargs ...interface{}) error {
	// Check pointer
	if reflect.ValueOf(model).Kind() != reflect.Ptr {
		return errors.New("Only pointer model excepted")
	}

	// Generate Query
	query := UpdateQuery(model, condition)

	// Extract variable
	v := reflect.Indirect(reflect.ValueOf(model))
	table := getTable(v.Type())

	args := make([]interface{}, 0)
	for i, f := range table.Fields {
		if i != table.PrimaryKey {
			switch f.Type {
			case reflect.TypeOf(true):
				args = append(args, v.Field(i).Bool())
			case reflect.TypeOf(0):
				args = append(args, v.Field(i).Int())
			case reflect.TypeOf(0.0):
				args = append(args, v.Field(i).Float())
			case reflect.TypeOf(""):
				args = append(args, v.Field(i).String())
			case reflect.TypeOf(time.Time{}):
				if (v.Field(i).Interface().(time.Time) == time.Time{}) {
					args = append(args, sql.NullInt64{})
				} else {
					args = append(args, v.Field(i).Interface())
				}
			default:
				args = append(args, v.Field(i).Interface())
			}
		}
	}

	// Check query
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	args = append(args, cargs...)
	// Insert Data
	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}

func delete(db IDB, model interface{}, condition string, cargs ...interface{}) error {
	// Check pointer
	if reflect.ValueOf(model).Kind() != reflect.Ptr {
		return errors.New("Only pointer model excepted")
	}

	// Generate Query
	query := DeleteQuery(model, condition)

	// Check query
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Insert Data
	_, err = stmt.Exec(cargs...)
	if err != nil {
		return err
	}

	return nil
}

func view(db IDB, model interface{}, query string, cargs ...interface{}) error {
	// Check pointer
	if reflect.ValueOf(model).Kind() != reflect.Ptr {
		return errors.New("Only pointer model excepted")
	}

	// Get table
	m := reflect.Indirect(reflect.ValueOf(model))
	var table SQLTable
	if m.Kind() == reflect.Slice {
		table = getTable(m.Type().Elem())
	} else {
		table = getTable(m.Type())
	}

	// Detect Struct or Slice
	var sliceMode bool
	switch m.Kind() {
	case reflect.Slice:
		sliceMode = true
	case reflect.Struct:
		sliceMode = false
	}

	// Read result
	rows, err := db.Query(query, cargs...)
	check(err)
	defer rows.Close()
	for rows.Next() {
		store := make([]interface{}, len(table.Fields))
		args := make([]interface{}, len(table.Fields))
		for i := range args {
			args[i] = &store[i]
		}
		if err := rows.Scan(args...); err != nil {
			check(err)
		}

		var obj reflect.Value
		if sliceMode {
			obj = reflect.New(m.Type().Elem()).Elem()
		} else {
			obj = m
		}

		// Convert result & reflect to model
		converted := ConvertSQLResult(model, store)
		for i, f := range table.Fields {
			switch f.Type {
			case reflect.TypeOf(false):
				obj.Field(i).SetBool(converted[i].(bool))
			case reflect.TypeOf(0):
				obj.Field(i).SetInt(converted[i].(int64))
			case reflect.TypeOf(0.0):
				obj.Field(i).SetFloat(converted[i].(float64))
			case reflect.TypeOf(""):
				obj.Field(i).SetString(converted[i].(string))
			case reflect.TypeOf(time.Time{}):
				if converted[i] != nil {
					obj.Field(i).Set(reflect.ValueOf(converted[i]))
				}
			default:
				if converted[i] != nil {
					obj.Field(i).Set(reflect.ValueOf(converted[i]))
				}
			}

		}

		// Detect Struct or Slice
		if sliceMode {
			m.Set(reflect.Append(m, obj))
		} else {
			break
		}
	}
	return nil
}
