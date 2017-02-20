package main

import (
	"./elevator"
	"./network"
	"fmt"
	"time"
)

func main() {
	fmt.Printf("Starting...! \n \n") //her må vi fortelle systemet at heisen er i live...
	controllCh := make(chan elevator.OrderMsg)
	broadcastCh := make(chan elevator.OrderMsg)
	msgRecCh := make(chan elevator.OrderMsg)


	go network.Network(controllCh, broadcastCh, msgRecCh)
	go elevator.ElevatorInit(msgRecCh)


	for {
		network.SendMsg(broadcastCh, elevator.Test)
		time.Sleep(1000 * time.Millisecond)	
	}
	time.Sleep(100 * time.Second)

}




