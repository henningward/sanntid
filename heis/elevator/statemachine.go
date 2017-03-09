package elevator

import (
	"../driver"
	"time"
	//"fmt"
)

func Statemachine(floorChan chan driver.FloorStatus, executeOrderChan chan Order, motorDir *driver.Direction, elev *ElevState, orderCostList *OrderList, newOrders *OrderList) {
	//startTime := time.Now().UnixNano()
	var orderToExecute Order
	elev.STATE = "IDLE" // må gjøres et annet sted..
	startTimeInState := time.Now().UnixNano()
	currentTime := time.Now().UnixNano()

	for {
		switch elev.STATE {
		case "IDLE":
			currentTime = time.Now().UnixNano()
			elev.TimeInState = (currentTime - startTimeInState) / 1000000
			elev.Dir = NONE
			select {
			case orderToExecute = <-executeOrderChan:
				elev.FloorStatus = driver.GetFloor(floorChan)
				elev.STATE = checkDirection(elev.FloorStatus, orderToExecute, motorDir)
				startTimeInState = time.Now().UnixNano()
			case <-time.After(100 * time.Millisecond):
				elev.FloorStatus = driver.GetFloor(floorChan)
				/*
					if stopAtFloor(elev.FloorStatus, orderToExecute) {
						DeleteOrder(orderToExecute, orderCostList, newOrders)
						//åpne dører osv.....'
					}
				*/
				driver.MotorIDLE()
				elev.STATE = "IDLE"

			}

		case "UP":
			elev.Dir = UP
			currentTime = time.Now().UnixNano()
			elev.TimeInState = (currentTime - startTimeInState) / 1000000
			select {
			case orderToExecute = <-executeOrderChan:

			case <-time.After(10 * time.Millisecond):
				driver.MotorUP()

				elev.FloorStatus = driver.GetFloor(floorChan)
				if stopAtFloor(elev.FloorStatus, orderToExecute) {
					DeleteOrder(orderToExecute, orderCostList, newOrders)
					//åpne dører osv.....'
					driver.MotorIDLE()
					elev.STATE = "DOORS OPEN"
					startTimeInState = time.Now().UnixNano()
				}
			}
		case "DOWN":
			elev.Dir = DOWN
			currentTime = time.Now().UnixNano()
			elev.TimeInState = (currentTime - startTimeInState) / 1000000
			select {
			case orderToExecute = <-executeOrderChan:

			case <-time.After(10 * time.Millisecond):
				driver.MotorDOWN()
				elev.FloorStatus = driver.GetFloor(floorChan)
				if stopAtFloor(elev.FloorStatus, orderToExecute) {
					DeleteOrder(orderToExecute, orderCostList, newOrders)
					//åpne dører osv.....
					driver.MotorIDLE()
					elev.STATE = "DOORS OPEN"
					startTimeInState = time.Now().UnixNano()
				}
			}
		case "DOORS OPEN":
			driver.SetDoorLamp(1)
			time.Sleep(2 * time.Second)
			driver.SetDoorLamp(0)
			DeleteOrder(orderToExecute, orderCostList, newOrders)
			elev.STATE = "IDLE"
			startTimeInState = time.Now().UnixNano()
		}
		//curTime := time.Now().UnixNano()
		//fmt.Println((curTime- startTime) / 1000000)
		//println(elev.Dir)
		//time.Sleep(10 * time.Millisecond)
	}
}
