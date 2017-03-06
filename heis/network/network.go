package network

import (
	"../elevator"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

const SPAMTIME = 1000 //milliseconds

func Network(controllCh chan elevator.OrderMsg, BroadcastCh chan elevator.OrderMsg, msgRecCh chan elevator.OrderMsg) {

	//port:= "20013" 149
	//ip:= "255.255.255.255"
	//service :=  fmt.Sprintf("%d:%d", ip, port)
	service := "129.241.187.255:34798"
	//service := "255.255.255.255:34899"
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
	ownIP := localAddr[:15]
	elevator.Test.IP = ownIP
	ownID := ownIP[len(ownIP)-1:]
	ownIDint, _ := strconv.Atoi(ownID)
	if ownIDint > 5 {
		ownIDint = ownIDint - 5
	}
	if ownIDint == 0 {

	}
	elevator.Test.ID = ownIDint

	connRec, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("Net.ListenUDP failed!\n")
		return
	}

	go Receive(connRec, msgRecCh, localAddr)

	go Broadcast(conn, BroadcastCh)

	for {

		time.Sleep(100 * time.Second)

	}

}

func SendMsg(msgChan chan elevator.OrderMsg, msg elevator.OrderMsg) {
	msgChan <- msg
}

func Broadcast(conn net.Conn, broadcastChan chan elevator.OrderMsg) {
	// skal sende meldingen vår med et intervall tilsvarende SPAMTIME
	var msg elevator.OrderMsg

	//var delay time.Time
	for {
		msg = <-broadcastChan
		//fmt.Printf("message ready! \n") //her må vi fortelle systemet at heisen er i live...

		//fmt.Printf("message ready! \n") //her må vi fortelle systemet at heisen er i live...

		//if time.Since(delay) > SPAMTIME*time.Millisecond { // her kan vi også sjekke om meldingen er valid...
		//delay = time.Now()

		jsonMsg, _ := json.Marshal(msg)
		conn.Write(jsonMsg)

		//}

	}
}

func Receive(connRec *net.UDPConn, MsgRecCh chan elevator.OrderMsg, localAddr string) {
	var msg elevator.OrderMsg
	var buf [1024]byte
	for {
		//fmt.Printf("message ready! \n") //her må vi fortelle systemet at heisen er i live...
		n, receivedAddr, _ := connRec.ReadFrom(buf[0:])

		//n, receivedAddr, _ := connRec.ReadFrom(buf[0:])
		json.Unmarshal(buf[0:n], &msg)

		//receivedAddr.String() = " " //fjerne denne for å forhindre at meldinger mottas på samme maskin

		//MsgRecCh <- msg
		//printOrdersRec(msg)
		if receivedAddr.String() != localAddr {
			MsgRecCh <- msg

		}

	}
}

func printOrdersRec(Test elevator.OrderMsg) {
	fmt.Printf("|FLOOR|   |UP|  |DOWN|  |INSIDE|  |COST|\n")
	temp := 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			temp++
			if temp%3 == 0 {
				fmt.Printf("   %v                X                %v \n", i+1,
					Test.Orders[j][i].Cost)
			}
			if temp%3 == 1 {
				fmt.Printf("   %v                        X        %v \n", i+1,
					Test.Orders[j][i].Cost)
			}
			if temp%3 == 2 {
				fmt.Printf("   %v       X                         %v \n", i+1,
					Test.Orders[j][i].Cost)
			}

		}
		time.Sleep(1 * time.Millisecond)
		fmt.Printf("\n")
	}
}
