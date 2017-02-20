package network

import (
	"../elevator"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const SPAMTIME = 1000 //milliseconds

func Network(controllCh chan elevator.OrderMsg, BroadcastCh chan elevator.OrderMsg, msgRecCh chan elevator.OrderMsg) {

	//port:= "20013" 149
	//ip:= "255.255.255.255"
	//service :=  fmt.Sprintf("%d:%d", ip, port)
	service := "129.241.187.255:34779"
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

	//hvorfor kan ikke denne deklareres "globalt", slik at receive ikke trenger å ta inn recChan?
	
	go Receive(connRec, msgRecCh, localAddr)

	go Broadcast(conn, BroadcastCh)
	time.Sleep(100 * time.Millisecond)
	
	for {
/*
		select {
		case dummy:= <-msgRecCh:
			_= dummy
			//printOrders(MsgRec)
			println("msg received")
		case <-time.After(100 * time.Millisecond):

		}
*/
		time.Sleep(100 * time.Millisecond)

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
		select {
		case msg = <-broadcastChan:
			//fmt.Printf("message ready! \n") //her må vi fortelle systemet at heisen er i live...
		}

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
		n, _, _ := connRec.ReadFrom(buf[0:])

		//n, receivedAddr, _ := connRec.ReadFrom(buf[0:])
		json.Unmarshal(buf[0:n], &msg)
		//receivedAddr.String() = " " //fjerne denne for å forhindre at meldinger mottas på samme maskin

		select {
		case MsgRecCh <- msg:
		case <-time.After(100 * time.Millisecond):
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

