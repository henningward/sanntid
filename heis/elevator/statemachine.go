package elevator

import (
	"../driver"
	//"fmt"
	"os/exec"
	"time"
)

const STUCKTIME = N_FLOORS * 2200 * time.Millisecond

func Statemachine(floorChan chan driver.FloorStatus, executeOrderChan chan Order, motorDir *driver.Direction, elev *ElevState, orderCostList *OrderList, newOrders *OrderList) {
	//startTime := time.Now().UnixNano()
	var orderToExecute Order
	var tempFloor driver.Floor
	elev.STATE = "IDLE" // må gjøres et annet sted..
	elev.StateTimer = time.Now()
	for {
		switch elev.STATE {
		case "IDLE":

			elev.Dir = NONE
			select {
			case orderToExecute = <-executeOrderChan:
				elev.FloorStatus = driver.GetFloor(floorChan)
				elev.STATE = checkDirection(elev.FloorStatus, orderToExecute, motorDir)
				elev.StateTimer = time.Now()
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
					elev.StateTimer = time.Now()
				}
			}
		case "DOWN":
			elev.Dir = DOWN
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
					elev.StateTimer = time.Now()
				}
			}
		case "DOORS OPEN":
			driver.SetDoorLamp(1)
			time.Sleep(2 * time.Second)
			driver.SetDoorLamp(0)
			DeleteOrder(orderToExecute, orderCostList, newOrders)
			elev.STATE = "IDLE"
			elev.StateTimer = time.Now()
		case "STUCK":
			beep := exec.Command("beep", "-r", "1", "beep", "-f", "1000")
			beep.Run()
			elev.FloorStatus = driver.GetFloor(floorChan)
			for i := 0; i < 3; i++ {
				for j := 0; j < N_FLOORS; j++ {
					if (orderCostList[i][j].Cost) != 0 && orderCostList[i][j].Button.Dir != NONE {
						orderCostList[i][j].Cost = 100000
					}
				}
			}
			if elev.FloorStatus.CurrentFloor != tempFloor {
				elev.STATE = "IDLE"
				driver.MotorIDLE()
				elev.StateTimer = time.Now()
				println("unstuck")
			}
		}
		if time.Since(elev.StateTimer) > STUCKTIME && elev.STATE != "IDLE" && elev.STATE != "STUCK" {
			elev.STATE = "STUCK"
			tempFloor = elev.FloorStatus.CurrentFloor
			println("Elevator stuck!")
		}

	}

}
