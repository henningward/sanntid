package elevator

import (
	"../driver"
	"time"
)

func Statemachine(floorChan chan driver.FloorStatus, executeOrderChan chan Order, motorDir *driver.Direction, elev *ElevState) {

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
			case <-time.After(10 * time.Millisecond):
				elev.FloorStatus = driver.GetFloor(floorChan)
				if stopAtFloor(elev.FloorStatus, orderToExecute) {
					DeleteOrder(orderToExecute)
					//åpne dører osv.....'
				}

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
					DeleteOrder(orderToExecute)
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
					DeleteOrder(orderToExecute)
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
			elev.STATE = "IDLE"
			startTimeInState = time.Now().UnixNano()
		}

	}
}
