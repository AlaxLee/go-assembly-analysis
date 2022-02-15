package main

import (
	"fmt"
)

func lele(i int, j uint, k ...string) (string, []interface{}) {
	return fmt.Sprintf("%d%d%s", i, j, k), []interface{}{i, j, k}
}

func main() {
	var a func(int, uint, ...string) (string, []interface{})
	a = lele
	fmt.Println(a)
}
