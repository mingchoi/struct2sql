package struct2sql

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

func Slice2CSV(path string, model interface{}, bom bool) error {
	// Create file
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write UTF-8 BOM header
	if bom == true {
		_, err = file.Write([]byte{0xef, 0xbb, 0xbf})
		if err != nil {
			return err
		}
	}

	// Convert
	list := Slice2Strings(model)

	// Write to file
	w := csv.NewWriter(file)
	w.WriteAll(list)

	if err = w.Error(); err != nil {
		return err
	}
	return nil
}

func Slice2Strings(model interface{}) (result [][]string) {
	s := reflect.ValueOf(model)
	t := s.Index(0).Type()
	if s.Kind() != reflect.Slice {
		panic("Only accept slice")
	}
	result = make([][]string, s.Len()+1)

	// Get column detail
	column := []string{}
	totalFields := t.NumField()
	for i := 0; i < totalFields; i++ {
		column = append(column, t.Field(i).Name)
	}
	result[0] = column

	// Convert
	for i := 0; i < s.Len(); i++ {
		result[i+1] = Struct2Strings(s.Index(i).Interface())
	}

	return
}

func Struct2Strings(model interface{}) (result []string) {
	// TODO: panic when not kind of struct
	result = []string{}
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Struct {
		panic("Only accept struct")
	}
	totalFields := v.NumField()
	for i := 0; i < totalFields; i++ {
		switch v.Field(i).Type() {
		case reflect.TypeOf(time.Time{}):
			val := v.Field(i).Interface().(time.Time)
			if (val == time.Time{}) {
				result = append(result, "")
			} else if v.Type().Field(i).Tag.Get("customtype") == "date" {
				result = append(result, val.Format("2006-01-02"))
			} else {
				result = append(result, val.Format("2006-01-02 15:04:05"))
			}
		default:
			result = append(result, fmt.Sprint(v.Field(i).Interface()))
		}
	}
	return
}

func ToCSVRowSQL(data interface{}) string {
	result := []string{}
	v := reflect.ValueOf(data)
	totalFields := v.NumField()
	for i := 1; i < totalFields; i++ {
		switch v.Field(i).Type() {
		case reflect.TypeOf(true):
			if v.Field(i).Interface().(bool) == true {
				result = append(result, "1")
			} else {
				result = append(result, "0")
			}
		case reflect.TypeOf(0):
			if v.Type().Field(i).Tag.Get("foreignkey") != "" && v.Field(i).Interface().(int) == 0 {
				result = append(result, "NULL")
			} else {
				result = append(result, "\""+fmt.Sprint(v.Field(i).Interface())+"\"")
			}
		case reflect.TypeOf(""):
			str := fmt.Sprint(v.Field(i).Interface())
			str = strings.Replace(str, "\\", "\\\\", -1)
			str = strings.Replace(str, "\"", "\\\"", -1)
			result = append(result, "\""+str+"\"")
		case reflect.TypeOf(time.Time{}):
			val := v.Field(i).Interface().(time.Time)
			if val.IsZero() {
				result = append(result, "NULL")
			} else {
				result = append(result, "\""+val.Format("2006-01-02 15:04:05")+"\"")
			}
		default:
			result = append(result, "\""+fmt.Sprint(v.Field(i).Interface())+"\"")
		}
	}
	return strings.Join(result, ",")
}
