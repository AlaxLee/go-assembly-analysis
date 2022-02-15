package main

import "fmt"

func main() {
	var a chan int = make(chan int, 1)
	fmt.Println(a)
}
