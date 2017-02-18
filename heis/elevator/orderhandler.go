package elevator

import (
	"../driver"
	"fmt"
	//"math/rand"
	"time"
)

const N_FLOORS = 4

const (
	NONE driver.Direction = iota
	UP
	DOWN
)

var orderlist [3][N_FLOORS]Order
var orderCostList [3][N_FLOORS]Order
var toExecute Order

type Order struct {
	Button driver.Button
	Cost   int
}

func SetOrder(buttonChan chan driver.Button) {
	var newButton driver.Button

	for {
		select {
		case newButton = <-buttonChan:
			dir, floor := newButton.Dir, newButton.Floor
			orderlist[dir][floor-1] = Order{newButton, 100000}
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func ComputeCost(elev *ElevState, motorDir *driver.Direction) {
	currentFloor := 0
	orderFloor := 0
	orderDir := NONE
	for {
		currentFloor = int(elev.FloorStatus.CurrentFloor)
		for i := 0; i < 3; i++ {
			for j := 0; j < N_FLOORS; j++ {
				if orderlist[i][j].Cost != 0 {
					orderCostList[i][j].Button = (orderlist[i][j]).Button
					orderFloor = int(orderlist[i][j].Button.Floor)
					orderDir = orderlist[i][j].Button.Dir
					cost := &orderCostList[i][j].Cost

					if elev.Dir == UP {
						if orderFloor >= currentFloor {
							if orderDir != DOWN {
								*cost = 10 * (orderFloor - currentFloor)
							} else {
								if highestFloorOrder() >= orderFloor {
									*cost = 10 * (2*highestFloorOrder() - orderFloor - currentFloor)
								} else {
									*cost = 10 * (orderFloor - currentFloor)
								}

							}
						} else {
							if orderDir != UP {
								*cost = 10 * (2*highestFloorOrder() - orderFloor - currentFloor)
							} else {
								if lowestFloorOrder() <= orderFloor {

									*cost = 10 * (2*highestFloorOrder() - 2*lowestFloorOrder() + orderFloor - currentFloor)
								} else {
									*cost = 10 * (2*highestFloorOrder() - currentFloor - orderFloor)
								}
							}
						}
					} else if elev.Dir == DOWN {

						if orderFloor <= currentFloor {
							if orderDir != UP {
								*cost = 10 * (currentFloor - orderFloor)
							} else {
								if lowestFloorOrder() <= orderFloor {
									*cost = 10 * (-2*lowestFloorOrder() + orderFloor + currentFloor)
								} else {
									*cost = 10 * (currentFloor - orderFloor)
								}

							}
						} else {
							if orderDir != DOWN {
								*cost = 10 * (-2*lowestFloorOrder() + orderFloor + currentFloor)
							} else {
								if highestFloorOrder() >= orderFloor {
									*cost = 10 * (2*highestFloorOrder() - 2*lowestFloorOrder() - orderFloor + currentFloor)
								} else {
									*cost = 10 * (-2*lowestFloorOrder() + currentFloor + orderFloor)
								}
							}
						}
					} else if elev.Dir == NONE {
						if elev.STATE == "IDLE" {
							if elev.TimeInState > 500 {
								if currentFloor > orderFloor {
									*cost = 10 * (currentFloor - orderFloor)
								} else {
									*cost = 10 * (orderFloor - currentFloor)
								}

							}
						}
					} else {
						fmt.Print("failed to compute cost at state %v", elev.STATE)
					}
					//elev.FloorStatus.CurrentFloor

					fmt.Println("computing..")
				}
			}

		}

		printOrders()
		//fmt.Printf("\n\n\n\n\n\n")
		time.Sleep(1 * time.Second)
	}
}

func ExecuteOrder(executeOrderChan chan Order) {
	toExecute.Cost = 10000
	for {
		for i := 0; i < 3; i++ {
			for j := 0; j < N_FLOORS; j++ {
				if (orderCostList[i][j].Cost < toExecute.Cost) && orderCostList[i][j].Cost != 0 {
					toExecute = orderCostList[i][j]
					executeOrderChan <- toExecute
					println("\n execute order \n")
				}
			}
		}
	}

}

func DeleteOrder(order Order) {
	emptyButton := driver.Button{0, 0}
	emptyOrder := Order{emptyButton, 0}
	driver.SetButtonLamp(order.Button, 0)
	orderCostList[order.Button.Dir][int(order.Button.Floor)-1] = emptyOrder
	orderlist[order.Button.Dir][int(order.Button.Floor)-1] = emptyOrder
	toExecute.Cost = 100000

}

func printOrders() {
	fmt.Printf("|FLOOR|   |UP|  |DOWN|  |INSIDE|  |COST|\n")
	temp := 0
	for i := 0; i < N_FLOORS; i++ {
		for j := 0; j < 3; j++ {
			temp++
			if temp%3 == 0 {
				fmt.Printf("   %v                X                %v \n", i+1,
					orderCostList[j][i].Cost)
			}
			if temp%3 == 1 {
				fmt.Printf("   %v                        X        %v \n", i+1,
					orderCostList[j][i].Cost)
			}
			if temp%3 == 2 {
				fmt.Printf("   %v       X                         %v \n", i+1,
					orderCostList[j][i].Cost)
			}

		}
		fmt.Printf("\n")
	}
}

func highestFloorOrder() int {
	highest := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < N_FLOORS; j++ {
			if orderlist[i][j].Cost != 0 {
				highest = N_FLOORS
			}
		}
	}
	return highest
}

func lowestFloorOrder() int {
	lowest := N_FLOORS + 1
	for i := 0; i < 3; i++ {
		for j := 0; j < N_FLOORS; j++ {
			if orderlist[i][j].Cost != 0 && j < lowest {
				lowest = j
			}
		}
	}
	return lowest
}
