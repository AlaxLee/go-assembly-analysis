package main

import (
	"fmt"
)

func lala(i int, j uint, k string) (string, []interface{}) {
	return fmt.Sprintf("%d%d%s", i, j, k), []interface{}{i, j, k}
}

func main() {
	var a func(int, uint, string) (string, []interface{})
	a = lala
	fmt.Println(a)
}
