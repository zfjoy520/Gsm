package controllers

import (
	"fmt"
	"github.com/robfig/revel"
	"gsm_test/app/models"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	fmt.Printf("Index...\n")

	// gsm.GetDb().Incr("hello")
	// s, _ := gsm.GetDb().Get("hello")
	// log.Printf("> hello = %s\n", s)
	new(models.Test).Save()

	return c.Render()
}

// type Test struct {
// 	Name string
// }

// func (t *Test) Save() error {
// 	t.Name = "nnnn"

// 	rv := reflect.ValueOf(t)
// 	dataStruct := reflect.Indirect(rv)
// 	if dataStruct.Kind() != reflect.Struct {
// 		return errors.New("expected a pointer to a struct")
// 	}

// 	var params string
// 	// struct to string
// 	dataStructType := dataStruct.Type()
// 	for i := 0; i < dataStructType.NumField(); i++ {
// 		fieldv := dataStruct.Field(i)

// 		fmt.Println(fieldv)
// 		if reflect.String == fieldv.Kind() {
// 			// reflect.ValueOf(params).Set("fieldv")
// 			// params = fieldv
// 		}
// 		fmt.Println(params)
// 	}
// 	return nil
// }
