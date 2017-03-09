package network

import (
	"../elevator"
	"encoding/json"
	"net"
	"os"
	"strconv"
	"time"
)

func Network(controllCh chan elevator.OrderMsg, BroadcastCh chan elevator.OrderMsg, msgRecCh chan elevator.OrderMsg) {
	service := "129.241.187.255:34798"
	addr, err := net.ResolveUDPAddr("udp4", service)

	if err != nil {
		println("Net.ResolveUDPAddr failed!")
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		println("Net.DialUDP failed!")
		println("Elevator offline.")
		println("Elevator not in use.")
		os.Exit(1)
		return
	} else {
		println("Elevator online! \n")
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	ownIP := localAddr[:15]
	println("Own IP: ")
	println(ownIP)
	println("\n")
	elevator.ElevatorMsg.IP = ownIP
	elevator.ElevatorMsg.ID = makeID(ownIP)
	connRec, err := net.ListenUDP("udp", addr)
	if err != nil {
		println("Net.ListenUDP failed!")
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
	var msg elevator.OrderMsg
	firsterr := false
	for {
		msg = <-broadcastChan

		jsonMsg, _ := json.Marshal(msg)
		_, err := conn.Write(jsonMsg)

		if err != nil {
			if !firsterr {
				println("Elevator offline.")
				println("Single elevator mode.")
				println("trying to reconnect...\n")
				firsterr = true
			}
		} else {
			if firsterr {
				println("\nReconnected!\n")
				firsterr = false
			}
		}

	}
}

func Receive(connRec *net.UDPConn, MsgRecCh chan elevator.OrderMsg, localAddr string) {
	var msg elevator.OrderMsg
	var buf [1024]byte
	for {
		n, receivedAddr, _ := connRec.ReadFrom(buf[0:])
		_ = receivedAddr
		json.Unmarshal(buf[0:n], &msg)
		if receivedAddr.String() != localAddr {
			MsgRecCh <- msg
		}
	}
}

func makeID(IP string) int {
	ID1 := IP[len(IP)-3:]
	ID2 := IP[len(IP)-2:]
	ID3 := IP[len(IP)-1:]
	ID1_int, _ := strconv.Atoi(ID1)
	ID2_int, _ := strconv.Atoi(ID2)
	ID3_int, _ := strconv.Atoi(ID3)
	IDint := ((ID1_int + 1) * (ID2_int + 1) * (ID3_int + 1)) % 500
	println(IDint)
	return IDint
}
