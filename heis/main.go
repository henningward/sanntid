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
	connDeadCh := make(chan elevator.OrderList)

	go network.Network(controllCh, broadcastCh, msgRecCh, connDead)
	go elevator.ElevatorInit(msgRecCh, connDead)

	for {
		network.SendMsg(broadcastCh, elevator.Test)
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(100 * time.Second)

}
