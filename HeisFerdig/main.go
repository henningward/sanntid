package main

import (
	"./elevator"
	"./network"
	"time"
)

func main() {
	controllCh := make(chan elevator.OrderMsg)
	broadcastCh := make(chan elevator.OrderMsg)
	msgRecCh := make(chan elevator.OrderMsg)

	go network.Network(controllCh, broadcastCh, msgRecCh)
	go elevator.ElevatorInit(msgRecCh)

	for {
		network.SendMsg(broadcastCh, elevator.ElevatorMsg)
		time.Sleep(80 * time.Millisecond)
	}
	time.Sleep(100 * time.Second)

}
