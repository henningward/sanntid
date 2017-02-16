package main

import (
	"./driver"
	"./elevator"
	"./network"
	"fmt"
	"time"
)

func main() {
	fmt.Printf("Starting...! \n \n") //her må vi fortelle systemet at heisen er i live...
	controllCh := make(chan elevator.TestMsg)
	broadcastCh := make(chan elevator.TestMsg)
	buttonChan := make(chan driver.Button)
	floorChan := make(chan driver.FloorStatus)

	var motorDir driver.Direction
	go network.Network(controllCh, broadcastCh)

	var test elevator.TestMsg
	test.Text = "hei :("
	test.Number = 2
	test.Cost = 3
	test.Id = 2

	go driver.Init(buttonChan, floorChan, &motorDir)
	go elevator.SetOrder(buttonChan)
	go elevator.ComputeCost(floorChan, &motorDir)

	go statemachine(floorChan, &motorDir)

	for {
		network.SendMsg(broadcastCh, test)
		time.Sleep(2000 * time.Millisecond)
	}

	time.Sleep(100 * time.Second)

}
