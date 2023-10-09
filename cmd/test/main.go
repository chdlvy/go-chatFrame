package main

import (
	"fmt"
	"reflect"
)

func main() {
	type aint int
	//var t aint
	t := 1
	val := reflect.ValueOf(&t)
	fmt.Println(val.Elem())
}
