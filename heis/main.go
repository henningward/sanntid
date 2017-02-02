package main

import (
	"./network"
	"./elevator"
	"fmt"
	"time"
)


func main(){
	fmt.Printf("Starting...! \n \n") //her må vi fortelle systemet at heisen er i live...
	var test elevator.TestMsg
	test.Text = "hei"
	test.Number = 2
	test.Cost = 3
	test.Id = 2

	testChan := make(chan elevator.TestMsg)

	go network.Network(testChan)
	time.Sleep(100*time.Second)
	


	
	

}