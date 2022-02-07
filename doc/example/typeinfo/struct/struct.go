package main

import "fmt"

type A struct {
	B int32  `haha`
	c string `hehe`
}

func (a *A) Lala() {
	fmt.Println(a)
}

func (a A) lele() {
	fmt.Println(a)
}

func main() {
	a := A{2, "kaka"}
	fmt.Println(a)
}
