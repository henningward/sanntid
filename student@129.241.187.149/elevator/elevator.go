package elevator

import (
	"../driver"
	//"fmt"
	"time"
)

type TestMsg struct {
	Text   string
	Number int
	Cost   int
	Id     int
}

type OrderMsg struct{
	Orders [3][N_FLOORS]Order
	Id int 

}

type ElevState struct {
	STATE           string
	PrevState       string
	TimeInState     int64
	FloorStatus     driver.FloorStatus
	Dir             driver.Direction
	ExectutingOrder Order
}

func checkDirection(currentFloorStatus driver.FloorStatus, orderToExecute Order, motorDir *driver.Direction) string {
	floordif := currentFloorStatus.CurrentFloor - orderToExecute.Button.Floor
	if floordif > 0 {
		return "DOWN"
	}
	if floordif < 0 {
		return "UP"
	}
	if floordif == 0 {
		return "IDLE" //OPEN DOOR
	}
	return "IDLE"
}

func stopAtFloor(currentFloorStatus driver.FloorStatus, orderToExecute Order) bool {
	return currentFloorStatus.CurrentFloor == orderToExecute.Button.Floor
}

func ElevatorInit(msgRecCh chan OrderMsg){
	buttonChan := make(chan driver.Button)
	floorChan := make(chan driver.FloorStatus)
	executeOrderChan := make(chan Order)

	var motorDir driver.Direction
	var elev ElevState

	Test.Id = 1
	// ser for meg at dette gjøres via nettverket på en eller annen måte.. iterere fra feks 1-20


	go ReceiveOrder(msgRecCh, &elev, executeOrderChan)
	go ExecuteRecOrder(executeOrderChan)
	go SetOrder(buttonChan)
	go ComputeCost(&elev, &motorDir)
	//	go ExecuteOrder(executeOrderChan)
	go Statemachine(floorChan, executeOrderChan, &motorDir, &elev)
	go driver.Init(buttonChan, floorChan, &motorDir)
	for {
		time.Sleep(100 * time.Second)
	}
}