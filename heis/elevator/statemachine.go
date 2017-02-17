package elevator

import(
	"../driver"
	"time"
)



func Statemachine(floorChan chan driver.FloorStatus, executeOrderChan chan Order, motorDir *driver.Direction, elev *ElevState){

	var orderToExecute Order
	elev.STATE = "IDLE" // må gjøres et annet sted..
	for{
	switch elev.STATE{	
		case "IDLE":

			select{
				case orderToExecute = <- executeOrderChan:
						elev.FloorStatus = driver.GetFloor(floorChan)
						elev.STATE = checkDirection(elev.FloorStatus, orderToExecute, motorDir)
						

				case <- time.After(10 * time.Millisecond):
					elev.FloorStatus = driver.GetFloor(floorChan)
				if stopAtFloor(elev.FloorStatus, orderToExecute){
					DeleteOrder(orderToExecute)
					//åpne dører osv.....'
				}

					driver.MotorIDLE()
					elev.STATE = "IDLE" 

					
			}

		case "UP":
			select{
			case orderToExecute = <- executeOrderChan:

			case <-time.After(10 * time.Millisecond):
				driver.MotorUP()

				elev.FloorStatus = driver.GetFloor(floorChan)
				if stopAtFloor(elev.FloorStatus, orderToExecute){
					DeleteOrder(orderToExecute)
					//åpne dører osv.....'
					driver.MotorIDLE()
					elev.STATE = "DOORS OPEN"
				}
}
		case "DOWN":
			select{
			case orderToExecute = <- executeOrderChan:

			case <-time.After(10 * time.Millisecond):
				driver.MotorDOWN()
				elev.FloorStatus = driver.GetFloor(floorChan)
				if stopAtFloor(elev.FloorStatus, orderToExecute){
					DeleteOrder(orderToExecute)
					//åpne dører osv.....
					driver.MotorIDLE()
					elev.STATE = "DOORS OPEN"
				}
			}
		case "DOORS OPEN":
			driver.SetDoorLamp(1)
			time.Sleep(2* time.Second)			
			driver.SetDoorLamp(0)
			elev.STATE = "IDLE"
}

}}
