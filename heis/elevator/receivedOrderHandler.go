package elevator

import (
	"../driver"
	"fmt"
	"time"
)

//lage en slice med backups.. som skal erstatte recOrders
//var recOrdersBackup[] OrderMsg

const TIMEOUT = 500 * time.Millisecond

type Connection struct {
	IP          string
	LastMsgTime time.Time
	Alive       bool
}

func ReceiveOrder(msgRecCh chan OrderMsg, elev *ElevState, executeOrderCh chan Order, motorDir *driver.Direction, orderCostList *OrderList, newOrders *OrderList) {
	ConnList := make([]Connection, 10)
	var msgRec OrderMsg
	var recOrders OrderMsg
	for {
		orderCostListMerged := *orderCostList
		msgRec = <-msgRecCh
		recOrders = msgRec
		connectionsStatus(recOrders, &ConnList)
		//printOrdersRec(msgRec)
		//recOrdersOwnCost = msgRec
		ComputeCost(elev, motorDir, &orderCostListMerged, &recOrders.Orders)
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

func connectionsStatus(recOrders OrderMsg, ConnList *[]Connection) {
	tempIP := recOrders.IP
	inList := false

	for i := 0; i < 10; i++ {
		if tempIP == (*ConnList)[i].IP {
			inList = true
			(*ConnList)[i].LastMsgTime = time.Now()

			if inList == false {
				for i := 0; i < 10; i++ {
					if (*ConnList)[i].IP == "" {
						newConn := Connection{IP: tempIP, LastMsgTime: time.Now(), Alive: true}
						(*ConnList)[i] = newConn
						break
					}
				}
			}

			if (time.Since((*ConnList)[i].LastMsgTime) > TIMEOUT) && ((*ConnList)[i].IP != "") {
				(*ConnList)[i].Alive = false
				//DO SOME SHIT
			} else {
				(*ConnList)[i].Alive = true
			}
		}

	}
}
