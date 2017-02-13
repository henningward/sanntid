package main

import (
	"./network"
	"fmt"
	"time"
)


func main(){
	fmt.Printf("Starting...! \n \n") //her må vi fortelle systemet at heisen er i live...
	controllCh := make(chan int)
	broadcastCh := make(chan int)

	go network.Network(controllCh, broadcastCh)
	var test int
	test = 0
	for {
		test++
		network.SendMsg(broadcastCh, test)
		time.Sleep(1000*time.Millisecond)
	}

	time.Sleep(100*time.Second)
	


	
	

}