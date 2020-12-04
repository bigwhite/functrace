package main

//go:generate ../../gen -w main.go

import "sync"

func A1() {
	B1()
}

func B1() {
	C1()
}

func C1() {
	D()
}

func A2() {
	B2()
}

func B2() {
	C2()
}

func C2() {
	D()
}

func D() {
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		A1()
		wg.Done()
	}()

	A2()
	wg.Wait()
}
