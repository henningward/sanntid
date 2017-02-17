package elevator

import(
	"../driver"
	//"fmt"
)

type TestMsg struct{
	Text string
	Number int
	Cost int
	Id int
}

type ElevState struct {
	STATE string
	PrevState string
	FloorStatus driver.FloorStatus
	Dir driver.Direction
	ExectutingOrder Order

}


func checkDirection(currentFloorStatus driver.FloorStatus, orderToExecute Order, motorDir *driver.Direction) string{
	floordif := currentFloorStatus.CurrentFloor - orderToExecute.Button.Floor
	if floordif > 0{
		return "DOWN"
	}
	if floordif < 0{
		return "UP"
	}
	if floordif == 0{
		return "IDLE" //OPEN DOOR
	}
	return "IDLE"
}

func stopAtFloor(currentFloorStatus driver.FloorStatus, orderToExecute Order) bool{
	return currentFloorStatus.CurrentFloor == orderToExecute.Button.Floor
}

