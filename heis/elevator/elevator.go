package elevator

import (
	"../driver"
	//"fmt"
	"time"
	"fmt"
    "os"
)

type OrderList [3][N_FLOORS]Order

type OrderMsg struct {
	IP          string
	ID          int
	LastMsgTime time.Time
	Orders      OrderList
}

type ElevState struct {
	STATE           string
	PrevState       string
	TimeInState     int64
	FloorStatus     driver.FloorStatus
	Dir             driver.Direction
	ExectutingOrder Order
}

type Connection struct {
	IP          string
	LastMsgTime time.Time
	Alive       bool
	Orders      OrderList
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
		return "DOORS OPEN" //OPEN DOOR
	}

	return "IDLE"
}

func stopAtFloor(currentFloorStatus driver.FloorStatus, orderToExecute Order) bool {
	return (currentFloorStatus.CurrentFloor == orderToExecute.Button.Floor)
}

func ElevatorInit(msgRecCh chan OrderMsg) {
	buttonChan := make(chan driver.Button)
	floorChan := make(chan driver.FloorStatus)
	executeOrderChan := make(chan Order)

	ConnList := make([]Connection, 10)

	var newOrders OrderList
	var orderCostList OrderList
	var motorDir driver.Direction
	var elev ElevState

    
    
	if _, err := os.Stat("./backup"); os.IsNotExist(err) {
		os.Create("./backup")
}
	file, err := os.Open("./backup")
	if err != nil{
		fmt.Println("failed to open backup file")
	}
	importOrders(file)

	//d2 := []byte{115, 111, 109, 101, 10}
    //f.Write(d2)



	go ReceiveOrder(msgRecCh, &elev, executeOrderChan, &motorDir, &orderCostList, &newOrders, &ConnList)
	go SetOrder(buttonChan, &newOrders)
	go func() {
		for {
			Test.Orders = orderCostList
			ComputeCost(&elev, &motorDir, &orderCostList, &newOrders, Test.ID)
			time.Sleep(10 * time.Millisecond)
		}

	}()
	go ExecuteOrder(executeOrderChan, &orderCostList)
	go Statemachine(floorChan, executeOrderChan, &motorDir, &elev, &orderCostList, &newOrders)
	go driver.Init(buttonChan, floorChan, &motorDir)
	go checkConnections(&ConnList, &newOrders)
	for {
		time.Sleep(100 * time.Second)
	}
}


func importOrders(file *os.File){
	data := make([]byte, 100)
	count, err := file.Read(data)
	if err != nil {
		fmt.Println("backup file was empty")
	}
	fmt.Printf("backup file contains %d orders at floor: %q \n", count-1, data[:count-1])
}