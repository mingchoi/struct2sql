package struct2sql

import (
	"database/sql"
	"reflect"
	"strconv"
	"time"
)

//
//
//	Convertion
//
//

type IParsableType interface {
	Type() string
}

func GoType2SQLType(fieldType reflect.Type) string {
	switch fieldType {
	case reflect.TypeOf(true):
		return "bool"
	case reflect.TypeOf(0):
		return "int"
	case reflect.TypeOf(0.0):
		return "double"
	case reflect.TypeOf(""):
		return "varchar(255)"
	case reflect.TypeOf(time.Time{}):
		return "datetime"
	}

	// IParsableType
	if fieldType.Implements(reflect.TypeOf((*IParsableType)(nil)).Elem()) {
		return reflect.New(fieldType).Elem().Interface().(IParsableType).Type()
	}

	return "TYPE_NOT_SUPPORTED"
}

func ConvertSQLResult(model interface{}, args []interface{}) []interface{} {
	m := reflect.Indirect(reflect.ValueOf(model))
	var table SQLTable
	if m.Kind() == reflect.Slice {
		table = getTable(m.Type().Elem())
	} else {
		table = getTable(m.Type())
	}

	result := make([]interface{}, len(table.Fields))

	for i, f := range table.Fields {
		switch f.Type {
		// Base Type
		case reflect.TypeOf(true):
			result[i] = conv2bool(args[i])
		case reflect.TypeOf(0):
			result[i] = conv2int(args[i])
		case reflect.TypeOf(0.0):
			result[i] = conv2float(args[i])
		case reflect.TypeOf(""):
			result[i] = conv2str(args[i])
		case reflect.TypeOf(time.Time{}):
			result[i] = args[i]
		default:
			val := reflect.New(f.Type)
			val.Elem().Addr().Interface().(sql.Scanner).Scan(args[i])
			result[i] = val.Elem().Interface()
		}
	}
	return result
}

func conv2bool(val interface{}) bool {
	var result bool
	switch val.(type) {
	case bool:
		return val.(bool)
	case int64:
		if val.(int64) == 1 {
			return true
		} else {
			return false
		}
	case []uint8:
		var err error
		result, err = strconv.ParseBool(string(val.([]byte)))
		check(err)
	}
	return result
}

func conv2int(val interface{}) int64 {
	var result int64
	switch val.(type) {
	case int64:
		return val.(int64)
	case int:
		return int64(val.(int))
	case []uint8:
		var err error
		result, err = strconv.ParseInt(string(val.([]byte)), 10, 64)
		check(err)
	}
	return result
}

func conv2float(val interface{}) float64 {
	var result float64
	switch val.(type) {
	case float64:
		return val.(float64)
	case []uint8:
		var err error
		result, err = strconv.ParseFloat(string(val.([]byte)), 64)
		check(err)
	}
	return result

}

func conv2str(val interface{}) string {
	var result string
	switch val.(type) {
	case string:
		return val.(string)
	case []uint8:
		result = string(val.([]byte))
	}
	return result
}
