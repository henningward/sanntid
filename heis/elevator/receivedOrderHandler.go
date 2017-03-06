package elevator

import (
	"../driver"
	"fmt"
	"time"
)

//lage en slice med backups.. som skal erstatte recOrders
//var recOrdersBackup[] OrderMsg

const TIMEOUT = 500 * time.Millisecond

var timerRecOrders time.Time

func ReceiveOrder(msgRecCh chan OrderMsg, elev *ElevState, executeOrderCh chan Order, motorDir *driver.Direction, orderCostList *OrderList, newOrders *OrderList, ConnList *[]Connection) {, 
	var msgRec OrderMsg
	var recOrders OrderMsg
	for {
		orderCostListMerged := *orderCostList
		msgRec = <-msgRecCh
		recOrders = msgRec

		ignoreInternalOrders(&recOrders)

		setAllLamps(recOrders)

		if isNewMessage(recOrders, ConnList) {
			timerRecOrders = time.Now()
		}
		updateConnections(recOrders, ConnList)
		//printOrdersRec(msgRec)
		//recOrdersOwnCost = msgRec
		ComputeCost(elev, motorDir, &orderCostListMerged, &recOrders.Orders, recOrders.ID)
		compareCost(orderCostList, &recOrders, &orderCostListMerged, newOrders)
		time.Sleep(10 * time.Millisecond)
	}
}

func printOrdersRec(Test OrderMsg) {
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
		time.Sleep(10 * time.Millisecond)
		fmt.Printf("\n")
	}
}

func compareCost(orderCostList *OrderList, recOrders *OrderMsg, orderCostListMerged *OrderList, newOrders *OrderList) {
	for i := 0; i < 3; i++ {

		for j := 0; j < N_FLOORS; j++ {
			if orderCostListMerged[i][j].Cost < recOrders.Orders[i][j].Cost && orderCostListMerged[i][j].Cost != 0 {
				orderCostList[i][j] = orderCostListMerged[i][j]
				println("setting..")
				driver.SetButtonLamp(orderCostListMerged[i][j].Button, 1)
			}
			//println(recOrders.Orders[i][j])
			//printOrderss(orderCostListMerged)
			//println(recOrders.Orders[i][j].Cost)
			//if (recOrders.Orders[i][j].Cost != 0){
			//	fmt.Printf("rec: %v    order cost:%v \n", recOrders.Orders[i][j].Cost, orderCostList[i][j].Cost)
			//}

			if recOrders.Orders[i][j].Cost < orderCostList[i][j].Cost && recOrders.Orders[i][j].Cost != 0 {

				DeleteOrder(orderCostList[i][j], orderCostList, newOrders)
				println("deleting...")

			}
			//printOrderss(orderCostList)
		}
	}
}

func printOrderss(test *OrderList) {
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
		fmt.Printf("\n")
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

func setAllLamps(recOrders OrderMsg) {
	for i := 0; i < 3; i++ {
		for j := 0; j < N_FLOORS; j++ {
			if recOrders.Orders[i][j].Cost != 0 {
				driver.SetButtonLamp(recOrders.Orders[i][j].Button, 1)
			}
		}
	}

}

//Connections burde strent tatt være i nettverk, men utrolig knot å få det til grunnet at man ikke kan ha pakker i sykler
func updateConnections(recOrders OrderMsg, ConnList *[]Connection) {
	tempIP := recOrders.IP
	tempOrders := recOrders.Orders
	inList := false

	for i := 0; i < 10; i++ {
		if (*ConnList)[i].IP == tempIP {
			//println("updating existing connection \n")
			inList = true
			(*ConnList)[i].LastMsgTime = time.Now()
			(*ConnList)[i].Orders = tempOrders
		}
	}

	if inList == false {
		//println("new connection \n")
		for i := 0; i < 10; i++ {
			if (*ConnList)[i].IP == "" {
				newConn := Connection{IP: tempIP, LastMsgTime: time.Now(), Alive: true}
				(*ConnList)[i] = newConn
				println((*ConnList)[i].IP)
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
			} else {
				(*ConnList)[i].Alive = true
			}

			if ((*ConnList)[i].IP != "") && ((*ConnList)[i].Alive == false) {
				*newOrder = (*ConnList)[i].Orders
				(*ConnList)[i].IP = ""

			}
		}
	}
}
