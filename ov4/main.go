package main

import (
	"./network"
)


type TestMsg struct{
	text string
	number int
	cost int
	id int
}

func main(){
	var test TestMsg
	test.text = "hei"
	test.number = 2
	test.cost = 3
	test.id = 2
	Network
	
	testChan := make(chan TestMsg)
	testChan <- test

	network.Network(testChan)


}