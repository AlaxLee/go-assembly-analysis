package main

import "fmt"

func main() {
	b := [3]int{1, 2, 3}
	var a *[3]int
	a = &b
	fmt.Println(a)
}
