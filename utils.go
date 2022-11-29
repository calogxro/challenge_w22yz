package main

import (
	"fmt"
)

func TypeOf(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

// func TypeOf(v interface{}) string {
// 	return reflect.TypeOf(v).String()
// }
