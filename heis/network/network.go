package network

import (
	"encoding/json"
	"fmt"
	"net"
	"../elevator"
	"time"
)

const SPAMTIME = 1000 //milliseconds


func Network(controllCh chan elevator.TestMsg){
	
	
	//port:= "20013"
	//ip:= "255.255.255.255"
	//service :=  fmt.Sprintf("%d:%d", ip, port)
	service := "129.241.187.255:3477"
	addr, err := net.ResolveUDPAddr("udp4", service)

	if err != nil {
		fmt.Printf("Net.ResolveUDPAddr failed!\n")
		return 
	}	

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Printf("Net.DialUDP failed!\n")
		return 
	}	

	

	//broadcastChan := make(chan elevator.TestMsg)
	//go Broadcast(conn, broadcastChan)

	defer conn.Close()

	localAddr := conn.LocalAddr().String()

	connRec, err := net.ListenUDP("udp", addr)
	if err != nil {
			fmt.Printf("Net.ListenUDP failed!\n")
			return 
		}	

	recChan := make(chan elevator.TestMsg)
	go Receive(connRec, recChan, localAddr)
	
	




	testChan := make(chan elevator.TestMsg)
	var test elevator.TestMsg
	test.Text = "hei"
	test.Number = 2
	test.Cost = 3
	test.Id = 2
	
	go func(){
		for{
			testChan <- test
		}
	}()
	

	go Broadcast(conn, testChan)
	

}













func Broadcast(conn net.Conn, broadcastChan chan elevator.TestMsg) {
	// skal sende meldingen vår med et intervall tilsvarende SPAMTIME
	var msg elevator.TestMsg
	var delay time.Time 
	for {
		select{
			case msg = <- broadcastChan:
				//fmt.Printf("message ready! \n") //her må vi fortelle systemet at heisen er i live...
		}

		if time.Since(delay) > SPAMTIME*time.Millisecond { // her kan vi også sjekke om meldingen er valid...
			delay = time.Now()
			jsonMsg, _ := json.Marshal(msg)
			conn.Write(jsonMsg)

		}

	}
}




func Receive(connRec *net.UDPConn, recChan chan elevator.TestMsg, localAddr string){
	var msg elevator.TestMsg
	var buf [1024]byte
	for {
		fmt.Printf("message ready! \n") //her må vi fortelle systemet at heisen er i live...
		n, _, _ := connRec.ReadFrom(buf[0:])

		//n, receivedAddr, _ := connRec.ReadFrom(buf[0:])
		json.Unmarshal(buf[0:n], &msg)
		//receivedAddr.String() = " " //fjerne denne for å forhindre at meldinger mottas på samme maskin
		
		select {
				case recChan <- msg:
					fmt.Println(recChan)
				case <-time.After(100*time.Millisecond):
			}
/*
		if (receivedAddr.String() != localAddr){
			select {
				case recChan <- msg:
					fmt.Println(recChan)
			}
		}
*/
	}
}