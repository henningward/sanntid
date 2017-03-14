package elevator

import (
	"../driver"
	"time"
)

const TIMEOUT = 500 * time.Millisecond

var timerRecOrders time.Time

func ReceiveOrder(msgRecCh chan OrderMsg, elev *ElevState, executeOrderCh chan Order, motorDir *driver.Direction, orderCostList *OrderList, newOrders *OrderList, ConnList *[]Connection) {
	var msgRec OrderMsg
	var recOrders OrderMsg
	for {
		orderCostListMerged := *orderCostList
		msgRec = <-msgRecCh
		recOrders = msgRec
		ignoreInternalOrders(&recOrders)
		if isNewMessage(recOrders, ConnList) {
			timerRecOrders = time.Now()
		}

		updateConnections(recOrders, ConnList)
		ComputeCost(elev, motorDir, &orderCostListMerged, &recOrders.Orders, ElevatorMsg.ID)
		compareCost(orderCostList, &recOrders, &orderCostListMerged, newOrders)
		time.Sleep(10 * time.Millisecond)
	}
}

func compareCost(orderCostList *OrderList, recOrders *OrderMsg, orderCostListMerged *OrderList, newOrders *OrderList) {
	for i := 0; i < 3; i++ {
		for j := 0; j < N_FLOORS; j++ {
			if orderCostListMerged[i][j].Cost < recOrders.Orders[i][j].Cost && orderCostListMerged[i][j].Cost != 0 {
				orderCostList[i][j] = orderCostListMerged[i][j]
			}
			if recOrders.Orders[i][j].Cost < orderCostList[i][j].Cost && recOrders.Orders[i][j].Cost != 0 {
				DeleteOrder(orderCostList[i][j], orderCostList, newOrders)
			}
		}
	}
}

func isNewMessage(recOrders OrderMsg, ConnList *[]Connection) bool {
	newMessage := false
	for i := 0; i < 10; i++ {
		if ((*ConnList)[i].IP == recOrders.IP) && ((*ConnList)[i].Orders != recOrders.Orders) {
			newMessage = true
		}
	}
	return newMessage
}

func ignoreInternalOrders(recOrders *OrderMsg) {
	for i := 0; i < 3; i++ {
		for j := 0; j < N_FLOORS; j++ {
			if recOrders.Orders[i][j].Button.Dir == NONE {
				recOrders.Orders[i][j].Cost = 0
			}
		}
	}
}

func updateConnections(recOrders OrderMsg, ConnList *[]Connection) {
	tempIP := recOrders.IP
	tempOrders := recOrders.Orders
	inList := false
	for i := 0; i < 10; i++ {
		if (*ConnList)[i].IP == tempIP {
			inList = true
			(*ConnList)[i].LastMsgTime = time.Now()
			(*ConnList)[i].Orders = tempOrders
			(*ConnList)[i].Alive = true

		}
	}

	if inList == false {
		println("Connected to elevator:")
		for i := 0; i < 10; i++ {
			if (*ConnList)[i].IP == "" {
				newConn := Connection{IP: tempIP, LastMsgTime: time.Now(), Alive: true, Orders: tempOrders}
				(*ConnList)[i] = newConn
				println((*ConnList)[i].IP)
				println()
				break
			}
		}
	}

}

func checkConnections(ConnList *[]Connection, newOrders *OrderList) {

	for {
		for i := 0; i < 10; i++ {
			if ((*ConnList)[i].IP != "") && (time.Since((*ConnList)[i].LastMsgTime) > TIMEOUT) {
				(*ConnList)[i].Alive = false
				*newOrders = (*ConnList)[i].Orders
				println("lost connection to:")
				println((*ConnList)[i].IP)
				println()
				(*ConnList)[i].IP = ""
			}

		}
		time.Sleep(10 * time.Millisecond)
	}
}
