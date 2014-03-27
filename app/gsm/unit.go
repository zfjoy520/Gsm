package gsm

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func scanMapIntoStruct(obj interface{}, objMap map[string][]byte) error {
	dataStruct := reflect.Indirect(reflect.ValueOf(obj))
	if dataStruct.Kind() != reflect.Struct {
		return errors.New("expected a pointer to a struct")
	}

	dataStructType := dataStruct.Type()

	for i := 0; i < dataStructType.NumField(); i++ {
		field := dataStructType.Field(i)
		fieldv := dataStruct.Field(i)

		err := scanMapElement(fieldv, field, objMap)
		if err != nil {
			return err
		}
	}

	return nil
}

func scanMapElement(fieldv reflect.Value, field reflect.StructField, objMap map[string][]byte) error {

	objFieldName := field.Name
	bb := field.Tag
	sqlTag := bb.Get("sql")

	if bb.Get("beedb") == "-" || sqlTag == "-" || reflect.ValueOf(bb).String() == "-" {
		return nil
	}
	sqlTags := strings.Split(sqlTag, ",")
	sqlFieldName := objFieldName
	if len(sqlTags[0]) > 0 {
		sqlFieldName = sqlTags[0]
	}
	inline := false
	//omitempty := false //TODO!
	// CHECK INLINE
	if len(sqlTags) > 1 {
		// if stringArrayContains("inline", sqlTags[1:]) {
		// 	inline = true
		// }
	}
	if inline {
		if field.Type.Kind() == reflect.Struct && field.Type.String() != "time.Time" {
			for i := 0; i < field.Type.NumField(); i++ {
				err := scanMapElement(fieldv.Field(i), field.Type.Field(i), objMap)
				if err != nil {
					return err
				}
			}
		} else {
			return errors.New("A non struct type can't be inline.")
		}
	}

	// not inline

	data, ok := objMap[sqlFieldName]

	if !ok {
		return nil
	}

	var v interface{}

	switch field.Type.Kind() {

	case reflect.Slice:
		v = data
	case reflect.String:
		v = string(data)
	case reflect.Bool:
		v = string(data) == "1"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		x, err := strconv.Atoi(string(data))
		if err != nil {
			return errors.New("arg " + sqlFieldName + " as int: " + err.Error())
		}
		v = x
	case reflect.Int64:
		x, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return errors.New("arg " + sqlFieldName + " as int: " + err.Error())
		}
		v = x
	case reflect.Float32, reflect.Float64:
		x, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return errors.New("arg " + sqlFieldName + " as float64: " + err.Error())
		}
		v = x
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return errors.New("arg " + sqlFieldName + " as int: " + err.Error())
		}
		v = x
	//Supports Time type only (for now)
	case reflect.Struct:
		if fieldv.Type().String() != "time.Time" {
			return errors.New("unsupported struct type in Scan: " + fieldv.Type().String())
		}

		x, err := time.Parse("2006-01-02 15:04:05", string(data))
		if err != nil {
			x, err = time.Parse("2006-01-02 15:04:05.000 -0700", string(data))

			if err != nil {
				return errors.New("unsupported time format: " + string(data))
			}
		}

		v = x
	default:
		return errors.New("unsupported type in Scan: " + reflect.TypeOf(v).String())
	}

	fieldv.Set(reflect.ValueOf(v))

	return nil
}

func Atos(any interface{}, str string) error {
	// field = reflect.ValueOf(any)
	field := reflect.ValueOf(any)
	// fmt.Println("0:", field.Tag, field.Kind(), reflect.String, field)
	var v interface{}
	switch field.Kind() {

	case reflect.String:
		v = field
		// case reflect.Slice:
		// 	v = data

		// case reflect.Bool:
		// 	v = string(data) == "1"
		// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		// 	x, err := strconv.Atoi(string(data))
		// 	if err != nil {
		// 		return errors.New("arg " + sqlFieldName + " as int: " + err.Error())
		// 	}
		// 	v = x
		// case reflect.Int64:
		// 	x, err := strconv.ParseInt(string(data), 10, 64)
		// 	if err != nil {
		// 		return errors.New("arg " + sqlFieldName + " as int: " + err.Error())
		// 	}
		// 	v = x
		// case reflect.Float32, reflect.Float64:
		// 	x, err := strconv.ParseFloat(string(data), 64)
		// 	if err != nil {
		// 		return errors.New("arg " + sqlFieldName + " as float64: " + err.Error())
		// 	}
		// 	v = x
		// case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// 	x, err := strconv.ParseUint(string(data), 10, 64)
		// 	if err != nil {
		// 		return errors.New("arg " + sqlFieldName + " as int: " + err.Error())
		// 	}
		// 	v = x
		// //Supports Time type only (for now)
		// case reflect.Struct:
		// 	if fieldv.Type().String() != "time.Time" {
		// 		return errors.New("unsupported struct type in Scan: " + fieldv.Type().String())
		// 	}

		// 	x, err := time.Parse("2006-01-02 15:04:05", string(data))
		// 	if err != nil {
		// 		x, err = time.Parse("2006-01-02 15:04:05.000 -0700", string(data))

		// 		if err != nil {
		// 			return errors.New("unsupported time format: " + string(data))
		// 		}
		// 	}

		// 	v = x
		// default:
		// 	return errors.New("unsupported type in Scan: " + reflect.TypeOf(v).String())
	}
	strv := reflect.ValueOf(str)
	fmt.Println("1:", reflect.TypeOf(str))
	fmt.Println("2:", reflect.TypeOf(strv))
	// strv.Set(reflect.ValueOf(v))
	fmt.Println("3:", reflect.TypeOf(v))

	return nil
}
