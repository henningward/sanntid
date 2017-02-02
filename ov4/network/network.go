package network

import (
	"encoding/json"
	"fmt"
	"net"
	"./main"
)

const SPAMTIME = 1000 //milliseconds


func Network(controllCh chan main.TestMsg){

	port:= "20013"
	ip:= "255.255.255.255"
	service =  fmt.Sprintf("%d:%d", ip, port)

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

	

	broadcastChan := make(chan main.Msg)
	go Broadcast(conn, chan broadcastChan)

	defer conn.Close()

	localAddr := conn.LocalAddr().String()

	connRec, err := net.ListenUDP("udp", addr)
	if err != nil {
			fmt.Printf("Net.ListenUDP failed!\n")
			return 
		}	

	recChan := make(chan main.Msg)
	go Receive(connRec, chan recChan, localAddr)
	

}













func Broadcast(conn net.Conn, broadcastChan chan main.TestMsg) {
	// skal sende meldingen vår med et intervall tilsvarende SPAMTIME
	var msg main.TestMsg
	var delay time.Time 
	
	for {
		select{
			case msg = <- broadcastChan:
				fmt.Printf("message ready! \n") //her må vi fortelle systemet at heisen er i live...
		}

		if time.Since(delay) > SPAMTIME*time.Milliseconds { // her kan vi også sjekke om meldingen er valid...
			delay = time.Now()
			jsonMsg, _ = json.Marshall(msg)
			conn.Write(jsonMsg)




		}

	}
}




func Receive(connRec *net.IPConn, recChan chan main.TestMsg, localAddr){
	var msg main.TestMsg
	var buf [1024]byte
	for {
		n, receivedAddr, _ := connRec.ReadFrom(buf[0:])
		json.Unmarshal()
		json.Unmarshal(buf[0:n], &msg)
		receivedAddr = 0 //fjerne denne for å forhindre at meldinger mottas på samme maskin
		if (receivedAddr.String() != localAddr){
			select {
				case recChan <- msg:
					fmt.Println(recChan)
			}
		}

	}
}