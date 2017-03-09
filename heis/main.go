package main

import (
	"./elevator"
	"./network"
	"fmt"

	"time"
)

func main() {
	/*
		beep := exec.Command("beep", "-r", "2", "beep", "-f", "1000")
		for {
			beep = exec.Command("beep", "-r", "2", "beep", "-f", "5000")
			beep.Run()
			beep = exec.Command("beep", "-r", "2", "beep", "-f", "1000")
			beep.Run()
			time.Sleep(100 * time.Millisecond)
		}*/
	fmt.Printf("Starting...! \n \n") //her må vi fortelle systemet at heisen er i live...
	controllCh := make(chan elevator.OrderMsg)
	broadcastCh := make(chan elevator.OrderMsg)
	msgRecCh := make(chan elevator.OrderMsg)

	go network.Network(controllCh, broadcastCh, msgRecCh)
	go elevator.ElevatorInit(msgRecCh)

	for {
		network.SendMsg(broadcastCh, elevator.Test)
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(100 * time.Second)

}
