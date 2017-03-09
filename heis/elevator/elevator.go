package elevator

import (
	"../driver"
	//"fmt"
	"fmt"
	"os"
	"strconv"
	"time"
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
	StateTimer      time.Time
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

	doneImporting := make(chan bool)

	if _, err := os.Stat("./backup"); os.IsNotExist(err) {
		os.Create("./backup")
	}
	file, err := os.Open("./backup")
	defer file.Close()
	if err != nil {
		fmt.Println("failed to open backup file")
	}

	go importOrders(file, buttonChan, doneImporting, &elev)
	go updateBackup(doneImporting, &orderCostList)
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
	go updateButtonLights(&orderCostList, &ConnList)
	for {
		time.Sleep(100 * time.Second)
	}
}

func importOrders(file *os.File, buttonChan chan driver.Button, doneImporting chan bool, elev *ElevState) {
	time.Sleep(1 * time.Second)
	data := make([]byte, 100)
	count, err := file.Read(data)
	if err != nil {
		fmt.Println("Backup file is empty.")
		println("\n")
	} else {
		fmt.Printf("backup file contains %d orders at floor: %d", count, data[0]-48)
		for i := 1; i < count; i++ {
			fmt.Printf(",")
			fmt.Printf(" %d", data[i]-48)
		}

		fmt.Printf("\n...importing \n")
		println("\n")
	}

	newButton := driver.Button{0, 0}
	for i := 0; i < count; i++ {
		newButton.Floor = driver.Floor(data[i] - 48)
		if elev.FloorStatus.CurrentFloor != newButton.Floor {
			buttonChan <- newButton
			driver.SetButtonLamp(newButton, 1)
			time.Sleep(10 * time.Millisecond)
		}

	}
	doneImporting <- true
}
func updateBackup(doneImporting chan bool, orderCostList *OrderList) {
	<-doneImporting
	for {
		file, err := os.Create("./backup")
		if err != nil {
			fmt.Println("failed to create new backup file")
		}
		for j := 0; j < N_FLOORS; j++ {
			if orderCostList[0][j].Cost != 0 {
				_, err = file.WriteString(strconv.Itoa(j + 1))
				if err != nil {
					fmt.Println("error writing to backupfile")
					time.Sleep(20 * time.Millisecond)
				}
			}
		}

		time.Sleep(100 * time.Millisecond)

	}

}

func updateButtonLights(orderCostList *OrderList, ConnList *[]Connection) {
	time.Sleep(1 * time.Second)
	driver.ClearButtonLights()

	for {
		driver.ClearButtonLights()

		for k := 0; k < 10; k++ {
			if (*ConnList)[k].IP != "" {
				for i := 0; i < 3; i++ {
					for j := 0; j < N_FLOORS; j++ {
						if (*ConnList)[k].Orders[i][j].Cost != 0 {
							driver.SetButtonLamp((*ConnList)[k].Orders[i][j].Button, 1)
						}
					}
				}
			}
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < N_FLOORS; j++ {
				if orderCostList[i][j].Cost != 0 {
					driver.SetButtonLamp(orderCostList[i][j].Button, 1)
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

/*for {
        noRecOrders := true
        for k := 0; k < 10; k++ {
            if (*ConnList)[k].IP != "" {
                recOrders = (*ConnList)[k].Orders
                noRecOrders = false
            }
        }
        _ = recOrders
        _ = noRecOrders
        for i := 0; i < 3; i++ {
            for j := 0; j < N_FLOORS; j++ {
                if (orderCostList[i][j].Cost != 0) || (recOrders[i][j].Cost != 0 && !noRecOrders) {
                    driver.SetButtonLamp(orderCostList[i][j].Button, 1)
                } else {
                    driver.SetButtonLamp(orderCostList[i][j].Button, 0)
                }
                time.Sleep(10 * time.Millisecond)

            }
            time.Sleep(10 * time.Millisecond)
        }

        time.Sleep(10 * time.Millisecond)
    }
}*/
