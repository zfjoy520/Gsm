package models

import (
	"fmt"
	"gsm/app/gsm"
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
	t.Id = 12
	// t.Ip = "1111"
	// t.Int64 = 1213
	// t.Fat32 = 3.1415926
	// t.Name = "nnnn"
	// t.Ok = true
	// t.Arr = []string{"1", "2"}
	// t.IsAdmin = 1
	// gsm.Save(t)
	// fmt.Println(gsm.NewKey(t))
	fmt.Println(gsm.GetObj(t))

	fmt.Println("Struct:Id => ", t.Id)
	fmt.Println("Struct:Ip => ", t.Ip)
	fmt.Println("Struct:Int64 => ", t.Int64)
	fmt.Println("Struct:Name => ", t.Name)
	fmt.Println("Struct:Ok => ", t.Ok)
	fmt.Println("Struct:Arr => ", t.Arr)
	fmt.Println("Struct:Fat32 => ", t.Fat32)
	fmt.Println("Struct:IsAdmin => ", t.IsAdmin)
	return nil
}
