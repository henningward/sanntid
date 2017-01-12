package main

import (
	. "fmt"
	"runtime"
	"time"
)

var i = 0

func goRoutine1() {
	for j := 0; j < 1000000; j++ {
		i++
		time.Sleep(10 * time.Millisecond)
	}
}

func goRoutine2() {
	for j := 0; j < 1000000; j++ {
		i--
		time.Sleep(10 * time.Millisecond)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // I guess this is a hint to what GOMAXPROCS does...
	// Try doing the exercise both with and without it!
	go goRoutine1()
	go goRoutine2()
	// We have no way to wait for the completion of a goroutine (without additional syncronization of some sort)
	// We'll come back to using channels in Exercise 2. For now: Sleep.
	time.Sleep(100 * time.Millisecond)
	Println(i)
}
