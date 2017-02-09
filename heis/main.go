package main

import (
	"./network"
	"./elevator"
	"fmt"
	"time"
)


func main(){
	fmt.Printf("Starting...! \n \n") //her må vi fortelle systemet at heisen er i live...
	controllCh := make(chan elevator.TestMsg)
	broadcastCh := make(chan elevator.TestMsg)

	go network.Network(controllCh, broadcastCh)
	var test elevator.TestMsg
	test.Text = "hei :("
	test.Number = 2
	test.Cost = 3
	test.Id = 2

	for {
		network.SendMsg(broadcastCh, test)
		time.Sleep(1*time.Second)
	}

	time.Sleep(100*time.Second)
	


	
	

}