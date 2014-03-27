package gsm

import (
	"errors"
	"fmt"
	"menteslibres.net/gosexy/redis"
	"reflect"
	"strconv"
	"strings"
)

// const RedisMap
// const RedisMap map[string]*redis.Client
var RedisMap = make(map[string]*redis.Client)

func Init() error {
	var client *redis.Client
	client = redis.New()
	err := client.Connect("127.0.0.1", 1024)

	if err != nil {
		// fmt.Println("Connect failed: %s\n", err.Error())
		return errors.New("Connect failed:" + err.Error())
	}
	// 添加到常量中
	RedisMap["r1"] = client
	// fmt.Println("....%s", reflect.Indirect(RedisMap["r1"]))
	// client.Quit()
	return nil
}

func GetDb() *redis.Client {
	return RedisMap["r1"]
}

// func Get(key string) string {
// 	return GetDb().Get(key)
// }

func Save(obj interface{}) error {
	rv := reflect.ValueOf(obj)
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

		case reflect.String:
			fieldv = dfiled.String()

		case reflect.Bool:
			fieldv = dfiled.Bool()

		case reflect.Array, reflect.Slice:
			fieldv = dfiled.Interface().([]string)
			arrfield := reflect.ValueOf(fieldv)
			var str []string
			for i := 0; i < arrfield.Len(); i++ {
				str = append(str, arrfield.Index(i).String())
				// fmt.Println("Array://", str) //strings.Join(fieldv," "))
			}
			fieldv = strings.Join(str, " ")

		case reflect.Float32, reflect.Float64:
			fieldv = strconv.FormatFloat(dfiled.Float(), 'G', 30, 32)
			// fmt.Println("Float32://", fieldv)
		default:
			// return errors.New("unsupported type in Scan: ")
		}
		// params...
		params = append(params, field.Name, fieldv)
	}
	// fmt.Println(params)
	str, err := GetDb().HMSet(Key(obj), params...)
	fmt.Println(str)
	if err != nil {
		return err
	}
	return nil
}

func GetObj(obj interface{}) error {
	// 取值，得到指针
	rv := reflect.ValueOf(obj)
	// 取值描述*->取地址
	dataStruct := reflect.Indirect(rv)
	if dataStruct.Kind() != reflect.Struct {
		return errors.New("expected a pointer to a struct")
	}

	structmap := GetHash(obj)
	if len(structmap) == 0 {
		return errors.New("data is not exists")
	}

	dataStructType := dataStruct.Type()
	for i := 0; i < dataStructType.NumField(); i++ {
		field := dataStructType.Field(i)
		fieldv := dataStruct.Field(i)
		var thev interface{}
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
			// fieldv.Set(reflect.ValueOf(i))
			thev, _ = strconv.Atoi(structmap[field.Name])
			fieldv.Set(reflect.ValueOf(thev))
			fmt.Println(field.Name, "=>", thev)

		case reflect.Int64:

		}
		// if field.Type.Kind() == reflect.Int {

		// }
		// } else {
		// 	fieldv.Set(reflect.ValueOf(strconv.Itoa(i)))
		// }
		// fieldv.Set(reflect.ValueOf(thev))
		fmt.Println(structmap[field.Name])

	}

	return nil
}

func GetHash(obj interface{}) map[string]string {
	strs, _ := GetDb().HGetAll(Key(obj))
	slen := len(strs)
	var structmap = make(map[string]string)
	if slen > 0 && slen%2 == 0 {
		for i := 0; i < slen/2; i++ {
			structmap[strs[2*i]] = strs[2*i+1]
		}
	}
	return structmap
}

func NewKey(obj interface{}) string {
	sName := StructName(obj)
	key := sName + ":Id"

	GetDb().Incr(key)
	idstr, _ := GetDb().Get(key)

	return sName + ":" + idstr
}

func Key(obj interface{}) string {
	sName := StructName(obj)
	rv := reflect.ValueOf(obj)
	dataStruct := reflect.Indirect(rv)

	id := dataStruct.FieldByName("Id").Int()
	return sName + ":" + strconv.FormatInt(id, 10)
}

func StructName(obj interface{}) (name string) {
	typ := reflect.TypeOf(obj)
	tmp := strings.Split(typ.String(), ".")
	name = tmp[len(tmp)-1]
	return
}
