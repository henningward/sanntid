package elevator

import (
	"../driver"
	"fmt"
	"time"
)

type OrderList [3][N_FLOORS]Order

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
	fmt.Println(floordif)
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

	var newOrders OrderList
	var orderCostList OrderList
	var motorDir driver.Direction
	var elev ElevState

	Test.Id = 1
	// ser for meg at dette gjøres via nettverket på en eller annen måte.. iterere fra feks 1-20


	go ReceiveOrder(msgRecCh, &elev, executeOrderChan)
	go SetOrder(buttonChan, &newOrders)
	go ComputeCost(&elev, &motorDir, &orderCostList, &newOrders)
	go ExecuteOrder(executeOrderChan, &orderCostList)
	go Statemachine(floorChan, executeOrderChan, &motorDir, &elev, &orderCostList, &newOrders)
	go driver.Init(buttonChan, floorChan, &motorDir)
	for {
		time.Sleep(1 * time.Second)
	}
}