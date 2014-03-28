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

func valueToInterface(v reflect.Value) (inter interface{}, err error) {
	var fieldv interface{}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fieldv = v.Int()
		// fmt.Println("int://", field.Name, fieldv)

	case reflect.String:
		fieldv = v.String()

	case reflect.Bool:
		fieldv = v.Bool()

	case reflect.Float32, reflect.Float64:
		fieldv = strconv.FormatFloat(v.Float(), 'G', 30, 32)

	case reflect.Slice, reflect.Array:
		fmt.Println("///===", reflect.ValueOf(v.Interface()).Type())

		fieldv = v.Interface()
		arrfield := reflect.ValueOf(fieldv)
		var str []string
		for i := 0; i < arrfield.Len(); i++ {
			str = append(str, valueToString(arrfield.Index(i)))
		}
		fieldv = strings.Join(str, " ")

	default:
		return nil, errors.New("unsupported type in Scan: " + v.Kind().String())
	}
	return fieldv, nil
}

func valueToString(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.Itoa(int(v.Int()))
	case reflect.String:
		return v.String()

	case reflect.Bool:
		return strconv.FormatBool(v.Bool())

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'G', 30, 32)

	default:
		return ""
	}
	return ""
}

// func vauleToArray(v reflect.Value, tmp string) interface{} {
// 	_thefv = strings.Split(tmp, " ")
// 	var _inter []interface{}
// 	_inter
// }

func stringToValue(v reflect.Value, tmp string) (interface{}, error) {
	var thefv interface{}
	switch v.Kind() {
	case reflect.Int:
		thefv, _ = strconv.Atoi(tmp)

	case reflect.Int8:
		in8, _ := strconv.Atoi(tmp)
		thefv = int8(in8)

	case reflect.Int16:
		in16, _ := strconv.Atoi(tmp)
		thefv = int16(in16)

	case reflect.Int32:
		in32, _ := strconv.Atoi(tmp)
		thefv = int32(in32)

	case reflect.Int64:
		in64, _ := strconv.ParseInt(tmp, 10, 64)
		thefv = int64(in64)

	case reflect.Float32:
		fl32, _ := strconv.ParseFloat(tmp, 32)
		thefv = float32(fl32)

	case reflect.Float64:
		thefv, _ = strconv.ParseFloat(tmp, 64)

	case reflect.String:
		thefv = tmp

	case reflect.Bool:
		thefv, _ = strconv.ParseBool(tmp)

	case reflect.Slice, reflect.Array:
		// _vrr := strings.Split(tmp, " ")
		// fmt.Println("///===", v.Interface().([2]int), _vrr, len(_vrr))
		// _v := reflect.ValueOf(v.Interface()).Index(0).Kind()
		// fmt.Println(reflect.TypeOf("aaa"))
		// n := reflect.TypeOf("aaa")
		// var _arr []n

		// var _thefv []interface{}
		// for i := 0; i < len(_vrr); i++ {
		// 	fmt.Println("FOR:", i, _vrr[i])
		// 	_inter, err := stringToValue(_v.Index(i), _vrr[i])
		// 	fmt.Println("FOR:", i, _inter)

		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	_thefv = append(_thefv, _inter)
		// }
		stringToSlice(v, tmp)
		// thefv = []int{1, 2}
		return v.Interface(), nil
		// return nil, errors.New("unsupported type in Scan: " + v.Kind().String())
	default:
		return nil, errors.New("unsupported type in Scan: " + v.Kind().String())
	}
	return thefv, nil
}

func stringToSlice(v reflect.Value, tmp string) error {
	// var slice []int
	// a := 1
	// slice = append(slice, reflect.ValueOf(a))
	fmt.Println()
	_vrr := strings.Split(tmp, " ")
	switch v.Type().Elem().Kind() {
	case reflect.Int:
		var slice []int
		for i := 0; i < len(_vrr); i++ {
			_int, _ := strconv.Atoi(_vrr[i])
			slice = append(slice, _int)
		}
		v.Set(reflect.ValueOf(slice))

	case reflect.Int8:
		var slice []int8
		for i := 0; i < len(_vrr); i++ {
			_int, _ := strconv.Atoi(_vrr[i])
			_int8 := int8(_int)
			slice = append(slice, _int8)
		}
		v.Set(reflect.ValueOf(slice))

	case reflect.Int16:
		var slice []int16
		for i := 0; i < len(_vrr); i++ {
			_int, _ := strconv.Atoi(_vrr[i])
			_int16 := int16(_int)
			slice = append(slice, _int16)
		}
		v.Set(reflect.ValueOf(slice))

	case reflect.Int32:
		var slice []int32
		for i := 0; i < len(_vrr); i++ {
			_int, _ := strconv.Atoi(_vrr[i])
			_int32 := int32(_int)
			slice = append(slice, _int32)
		}
		v.Set(reflect.ValueOf(slice))

	case reflect.Int64:
		var slice []int64
		for i := 0; i < len(_vrr); i++ {
			_int, _ := strconv.Atoi(_vrr[i])
			_int64 := int64(_int)
			slice = append(slice, _int64)
		}
		v.Set(reflect.ValueOf(slice))

	case reflect.Float32:
		var slice []float32
		for i := 0; i < len(_vrr); i++ {
			fl32, _ := strconv.ParseFloat(_vrr[i], 32)
			_float32 := float32(fl32)
			slice = append(slice, _float32)
		}
		v.Set(reflect.ValueOf(slice))

	case reflect.Float64:
		var slice []float64
		for i := 0; i < len(_vrr); i++ {
			fl64, _ := strconv.ParseFloat(_vrr[i], 64)
			slice = append(slice, fl64)
		}
		v.Set(reflect.ValueOf(slice))

	case reflect.String:
		v.Set(reflect.ValueOf(_vrr))

	case reflect.Bool:
		var slice []bool
		for i := 0; i < len(_vrr); i++ {
			_bool, _ := strconv.ParseBool(_vrr[i])
			slice = append(slice, _bool)
		}
		v.Set(reflect.ValueOf(slice))

	default:
		return errors.New("unsupported type in Scan: " + v.Kind().String())
	}
	return nil
}
