package main

import "fmt"

type Foo struct{}

func (Foo) F() {
	D()
}

func A1() {
	B1()
}

func B1() {
	C1()
}

func C1() {
	D()
}

func D() {
}

func main() {
	A1()
	fmt.Println("ok")
}
