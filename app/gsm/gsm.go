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

		var fieldv interface{}
		fieldv = dataStruct.Field(i)
		if field.Type.Kind() == reflect.Int {
			fieldv = dataStruct.Field(i).Int()
		}
		// params...
		params = append(params, field.Name, fieldv)
	}
	fmt.Println(params)
	GetDb().HMSet(Key(obj), params...)
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
		// fieldv := dataStruct.Field(i)

		// if field.Type.Kind() == reflect.Int {
		// 	fieldv.Set(reflect.ValueOf(i))
		// } else {
		// 	fieldv.Set(reflect.ValueOf(strconv.Itoa(i)))
		// }
		fmt.Println(structmap[field.Name])

	}
	// fmt.Println(dataStructType)
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
