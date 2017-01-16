package main

import (
	. "fmt"
	"runtime"
	//"time"
)

var i = 0

func goRoutine1(buffer(chan int), done(chan bool)) {
	for j := 0; j < 1000000; j++ {
		buffer <- 1
		i++
		<- buffer	
	}
	done <- true
}

func goRoutine2(buffer(chan int), done(chan bool)) {
	for j := 0; j < 1000000; j++ {
		buffer <- 1
		i--
		<- buffer
	}
	done <- true
}

func main() {
	buffer := make(chan int, 1)
	done := make(chan bool, 2)
	runtime.GOMAXPROCS(runtime.NumCPU()) // I guess this is a hint to what GOMAXPROCS does...
	// Try doing the exercise both with and without it!
	go goRoutine1(buffer, done)
	go goRoutine2(buffer, done)
	// We have no way to wait for the completion of a goroutine (without additional syncronization of some sort)
	// We'll come back to using channels in Exercise 2. For now: Sleep.
	

	for i := 0; i < 2; i++ {
		<- done
	}

	Println(i)
}

//other ways to to wait for goroutines to finish, go to
//https://nathanleclaire.com/blog/2014/02/15/how-to-wait-for-all-goroutines-to-finish-executing-before-continuing/


