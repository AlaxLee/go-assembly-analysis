package main

import "fmt"

type I interface {
	lala(int, string) bool
}

type E interface {
}

type A struct{}

func (a A) lala(num int, str string) bool {
	fmt.Println(num, str)
	return true
}

func main() {
	var e E = 1
	fmt.Println(e)

	var a A
	var i I = a
	fmt.Println(i)
}
