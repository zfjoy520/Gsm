package models

import (
	"fmt"
	"gsm/app/gsm"
)

type Test struct {
	Id          int `PK`
	Ip          string
	Int8        int8
	Int16       int16
	Int32       int32
	Int64       int64
	Name        string
	Ok          bool
	Fat32       float32
	Fat64       float64
	Arr_int     []int
	Arr_int8    []int8
	Arr_int16   []int16
	Arr_int32   []int32
	Arr_int64   []int64
	Arr_string  []string
	Arr_bool    []bool
	Arr_float32 []float32
	Arr_float64 []float64

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
	t.Ip = "1111"
	t.Int8 = 12
	t.Int16 = 1213
	t.Int32 = 121312
	t.Int64 = 1213123
	t.Fat32 = 3.1415
	t.Fat64 = 3.1415926
	t.Name = "nnnn"
	t.Ok = true

	t.Arr_int = []int{1, 2, 3, 4, 5}
	t.Arr_int8 = []int8{2, 3, 4, 5, 6}
	t.Arr_int16 = []int16{4, 5, 6, 7, 7}
	t.Arr_int32 = []int32{5, 6, 7, 8, 9}
	t.Arr_int64 = []int64{7, 6, 5, 4, 3}
	t.Arr_string = []string{"a", "b", "c"}
	t.Arr_bool = []bool{true, false, true, true}
	t.Arr_float32 = []float32{1.1111, 1.22222, 1.33333}
	t.Arr_float64 = []float64{2.33333333, 2.444444444, 2.12345}

	t.IsAdmin = 1
	gsm.Save(t)
	fmt.Println(gsm.NewKey(t))
	gsm.GetObj(t)

	fmt.Println("Struct:Id => ", t.Id)
	fmt.Println("Struct:Ip => ", t.Ip)
	fmt.Println("Struct:Int8 => ", t.Int8)
	fmt.Println("Struct:Int16 => ", t.Int16)
	fmt.Println("Struct:Int32 => ", t.Int32)
	fmt.Println("Struct:Int64 => ", t.Int64)
	fmt.Println("Struct:Name => ", t.Name)
	fmt.Println("Struct:Ok => ", t.Ok)
	fmt.Println("Struct:Fat32 => ", t.Fat32)
	fmt.Println("Struct:Fat64 => ", t.Fat64)
	fmt.Println("Struct:IsAdmin => ", t.IsAdmin)
	fmt.Println("Struct:Arr_int => ", t.Arr_int)
	fmt.Println("Struct:Arr_int8 => ", t.Arr_int8)
	fmt.Println("Struct:Arr_int16 => ", t.Arr_int16)
	fmt.Println("Struct:Arr_int32 => ", t.Arr_int32)
	fmt.Println("Struct:Arr_int64 => ", t.Arr_int64)
	fmt.Println("Struct:Arr_string => ", t.Arr_string)
	fmt.Println("Struct:Arr_bool => ", t.Arr_bool)
	fmt.Println("Struct:Arr_float32 => ", t.Arr_float32)
	fmt.Println("Struct:Arr_float64 => ", t.Arr_float64)

	return nil
}
