package models

import (
	"errors"
	"fmt"
	"gsm_test/app/gsm"
	"reflect"
	"strconv"
	"strings"
)

type Test struct {
	Id      int `PK`
	Ip      string
	Int64   int64
	Name    string
	Ok      bool
	Arr     []string
	Fat32   float32
	IsAdmin int
}

func (t *Test) Get(id int) {
	t.Id = id
	fmt.Println("--", t)
	gsm.GetObj(t)
	fmt.Println(t.Id)
	fmt.Println(t.Ip)
	fmt.Println(t.Name)
	fmt.Println(t.IsAdmin)

	// return true
}

func (t *Test) Save() error {
	t.Id = 11
	t.Ip = "1111"
	t.Int64 = 1213
	t.Fat32 = 3.1415926
	t.Name = "nnnn"
	t.Ok = true
	t.Arr = []string{"1", "2"}
	t.IsAdmin = 1
	// gsm.Save(t)
	// typ := reflect.TypeOf(t)
	// typestr := typ.String()
	// nameArray := strings.Split(typestr, ".")
	// typestr = nameArray[len(nameArray)-1]
	// if lastDotIndex != -1 {
	// 	typestr = typestr[lastDotIndex+1:]
	// }
	// fmt.Println(len(lastDotIndex))
	// fmt.Println(typestr)

	fmt.Println(gsm.Key(t))
	// fmt.Println(gsm.NewKey(t))
	fmt.Println(gsm.GetObj(t))

	rv := reflect.ValueOf(t)
	dataStruct := reflect.Indirect(rv)
	if dataStruct.Kind() != reflect.Struct {
		return errors.New("expected a pointer to a struct")
	}

	var params []interface{}
	dataStructType := dataStruct.Type()
	for i := 0; i < dataStructType.NumField(); i++ {
		field := dataStructType.Field(i)
		dfiled := dataStruct.Field(i)

		var fieldv interface{}
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldv = dfiled.Int()
			fmt.Println("int://", field.Name, fieldv)

		// case reflect.Int64:
		// 	fieldv = strconv.FormatInt(dfiled.String(), 10)
		case reflect.String:
			fieldv = dfiled.String()

		case reflect.Bool:
			fieldv = dfiled.Bool()
			fmt.Println(fieldv)

		case reflect.Array, reflect.Slice:
			// fmt.Println("Array://", dfiled) //strings.Join(fieldv," "))
			fieldv = dfiled.Interface().([]string)
			// ss := reflect.ValueOf(fieldv)
			sss := reflect.ValueOf(fieldv)
			var str []string
			for i := 0; i < sss.Len(); i++ {
				// fmt.Println(sss.Index(i))
				str = append(str, sss.Index(i).String())
				fmt.Println("Array://", str) //strings.Join(fieldv," "))
			}
			fieldv = strings.Join(str, " ")
		case reflect.Float32, reflect.Float64:
			fieldv = strconv.FormatFloat(dfiled.Float(), 'G', 30, 32)
			fmt.Println("Float32://", fieldv)
		default:
			// return errors.New("unsupported type in Scan: ")
		}
		// if field.Type.Kind() == reflect.Int {
		// 	fieldv = dfiled.Int()
		// }
		// params...
		params = append(params, field.Name, fieldv)
	}
	// fmt.Println(params)
	gsm.GetDb().HMSet(gsm.Key(t), params...)

	arr, _ := gsm.GetDb().HGetAll(gsm.Key(t))
	fmt.Println("++++:", arr)
	fmt.Println(len(arr))

	return nil
}
