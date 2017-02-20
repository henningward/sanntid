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
var OrderCostList [3][N_FLOORS]Order
var toExecute Order

var Test OrderMsg

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
					OrderCostList[i][j].Button = (orderlist[i][j]).Button
					orderFloor = int(orderlist[i][j].Button.Floor)
					orderDir = orderlist[i][j].Button.Dir
					cost := &OrderCostList[i][j].Cost

					if elev.Dir == UP {
						if orderFloor >= currentFloor {
							if orderDir != DOWN {
								*cost = 10 * (orderFloor - currentFloor)
							} else {
								if highestFloorOrder(currentFloor)>= orderFloor {
									*cost = 10 * (2*highestFloorOrder(currentFloor) - orderFloor - currentFloor)
								} else {
									*cost = 10 * (orderFloor - currentFloor)
								}

							}
						} else {
							if orderDir != UP {
								*cost = 10 * (2*highestFloorOrder(currentFloor)- orderFloor - currentFloor)
							} else {
								if lowestFloorOrder(currentFloor)<= orderFloor {

									*cost = 10 * (2*highestFloorOrder(currentFloor) - 2*lowestFloorOrder(currentFloor)+ orderFloor - currentFloor)
								} else {
									*cost = 10 * (2*highestFloorOrder(currentFloor)- currentFloor - orderFloor)
								}
							}
						}
					} else if elev.Dir == DOWN {

						if orderFloor <= currentFloor {
							if orderDir != UP {
								*cost = 10 * (currentFloor - orderFloor)
							} else {
								if lowestFloorOrder(currentFloor)<= orderFloor {
									*cost = 10 * (-2*lowestFloorOrder(currentFloor)+ orderFloor + currentFloor)
								} else {
									*cost = 10 * (currentFloor - orderFloor)
								}

							}
						} else {
							if orderDir != DOWN {
								*cost = 10 * (-2*lowestFloorOrder(currentFloor) + orderFloor + currentFloor)
							} else {
								if highestFloorOrder(currentFloor)>= orderFloor {
									*cost = 10 * (2*highestFloorOrder(currentFloor)- 2*lowestFloorOrder(currentFloor) - orderFloor + currentFloor)
								} else {
									*cost = 10 * (-2*lowestFloorOrder(currentFloor)+ currentFloor + orderFloor)
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

				}
			}

		}
		//printOrders()

		Test.Orders = OrderCostList
		//fmt.Printf("\n\n\n\n\n\n")
		time.Sleep(1 * time.Second)
	}
}
/*
func ExecuteOrder(executeOrderChan chan Order) {
	toExecute.Cost = 10000
	for {
		for i := 0; i < 3; i++ {
			for j := 0; j < N_FLOORS; j++ {
				if (OrderCostList[i][j].Cost < toExecute.Cost) && OrderCostList[i][j].Cost != 0 {
					toExecute = OrderCostList[i][j]
					executeOrderChan <- toExecute
				}
			}
		}
	}

}
*/

func ExecuteRecOrder(executeOrderChan chan Order) {
	toExecute.Cost = 10000
	for {
		for i := 0; i < 3; i++ {
			for j := 0; j < N_FLOORS; j++ {
				if (recorders[i][j].Cost < toExecute.Cost) && recorders[i][j].Cost != 0 {
					toExecute = recorders[i][j]
					executeOrderChan <- toExecute
				}
			}
		}
	}

}
func DeleteOrder(order Order) {
	emptyButton := driver.Button{0, 0}
	emptyOrder := Order{emptyButton, 0}
	driver.SetButtonLamp(order.Button, 0)
	OrderCostList[order.Button.Dir][int(order.Button.Floor)-1] = emptyOrder
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
					OrderCostList[j][i].Cost)
			}
			if temp%3 == 1 {
				fmt.Printf("   %v                        X        %v \n", i+1,
					OrderCostList[j][i].Cost)
			}
			if temp%3 == 2 {
				fmt.Printf("   %v       X                         %v \n", i+1,
					OrderCostList[j][i].Cost)
			}

		}
		fmt.Printf("\n")
	}
}

func highestFloorOrder(current int) int {
	highest := current
	for i := 0; i < 3; i++ {
		for j := 0; j < N_FLOORS; j++ {
			if orderlist[i][j].Cost != 0 {
				highest = N_FLOORS
			}
		}
	}
	return highest
}

func lowestFloorOrder(current int) int {
	lowest := current
	for i := 0; i < 3; i++ {
		for j := 0; j < N_FLOORS; j++ {
			if orderlist[i][j].Cost != 0 && j+1 < lowest {
				lowest = j+1
			}
		}
	}
	return lowest
}
