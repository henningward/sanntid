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



var toExecute Order

var Test OrderMsg

type Order struct {
	Button driver.Button
	Cost   int
}

func SetOrder(buttonChan chan driver.Button, newOrders *OrderList ) {
	var newButton driver.Button
	for {
		newButton = <-buttonChan
		dir, floor := newButton.Dir, newButton.Floor
		newOrders[dir][floor-1] = Order{newButton, 100000}		
}	
}


func ComputeCost(elev *ElevState, motorDir *driver.Direction, orderCostList *OrderList, newOrders *OrderList){
	currentFloor := 0
	orderFloor := 0
	orderDir := NONE
	for {
		currentFloor = int(elev.FloorStatus.CurrentFloor)
		for i := 0; i < 3; i++ {
			for j := 0; j < N_FLOORS; j++ {
				if newOrders[i][j].Cost != 0 {
					orderCostList[i][j].Button = (newOrders[i][j]).Button
					orderFloor = int(newOrders[i][j].Button.Floor)
					orderDir = newOrders[i][j].Button.Dir
					cost := &orderCostList[i][j].Cost
					
					if elev.Dir == UP {
						if orderFloor >= currentFloor {
							if orderDir != DOWN {
								*cost = 10 * (orderFloor - currentFloor)
							} else {
								if highestFloorOrder(currentFloor, orderCostList)>= orderFloor {
									*cost = 10 * (2*highestFloorOrder(currentFloor, orderCostList) - orderFloor - currentFloor)
								} else {
									*cost = 10 * (orderFloor - currentFloor)
								}

							}
						} else {
							if orderDir != UP {
								*cost = 10 * (2*highestFloorOrder(currentFloor, orderCostList)- orderFloor - currentFloor)
							} else {
								if lowestFloorOrder(currentFloor, orderCostList)<= orderFloor {

									*cost = 10 * (2*highestFloorOrder(currentFloor, orderCostList) - 2*lowestFloorOrder(currentFloor, orderCostList)+ orderFloor - currentFloor)
								} else {
									*cost = 10 * (2*highestFloorOrder(currentFloor, orderCostList)- currentFloor - orderFloor)
								}
							}
						}
					} else if elev.Dir == DOWN {

						if orderFloor <= currentFloor {
							if orderDir != UP {
								*cost = 10 * (currentFloor - orderFloor)
							} else {
								if lowestFloorOrder(currentFloor, orderCostList)<= orderFloor {
									*cost = 10 * (-2*lowestFloorOrder(currentFloor, orderCostList)+ orderFloor + currentFloor)
								} else {
									*cost = 10 * (currentFloor - orderFloor)
								}

							}
						} else {
							if orderDir != DOWN {
								*cost = 10 * (-2*lowestFloorOrder(currentFloor, orderCostList) + orderFloor + currentFloor)
							} else {
								if highestFloorOrder(currentFloor, orderCostList)>= orderFloor {
									*cost = 10 * (2*highestFloorOrder(currentFloor, orderCostList)- 2*lowestFloorOrder(currentFloor, orderCostList) - orderFloor + currentFloor)
								} else {
									*cost = 10 * (-2*lowestFloorOrder(currentFloor, orderCostList)+ currentFloor + orderFloor)
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
		//printOrders(orderCostList)

		Test.Orders = *orderCostList
		//fmt.Printf("\n\n\n\n\n\n")
		time.Sleep(10 * time.Millisecond)
	}
}

func ExecuteOrder(executeOrderChan chan Order, orderCostList *OrderList) {
	toExecute.Cost = 10000
	for {
		for i := 0; i < 3; i++ {
			for j := 0; j < N_FLOORS; j++ {
				if (orderCostList[i][j].Cost < toExecute.Cost) && orderCostList[i][j].Cost != 0 {
					toExecute = orderCostList[i][j]
					executeOrderChan <- toExecute
					
				}
			}
		}
		time.Sleep(10*time.Millisecond)
	}
}


func DeleteOrder(order Order, orderCostList *OrderList, newOrders *OrderList) {
	emptyButton := driver.Button{0, 0}
	emptyOrder := Order{emptyButton, 0}
	driver.SetButtonLamp(order.Button, 0)
	orderCostList[order.Button.Dir][int(order.Button.Floor)-1] = emptyOrder
	newOrders[order.Button.Dir][int(order.Button.Floor)-1] = emptyOrder
	toExecute.Cost = 100000

}



func highestFloorOrder(current int, orderCostList *OrderList) int {
	highest := current
	for i := 0; i < 3; i++ {
		for j := 0; j < N_FLOORS; j++ {
			if orderCostList[i][j].Cost != 0 {
				highest = N_FLOORS
			}
		}
	}
	return highest
}

func lowestFloorOrder(current int, orderCostList *OrderList) int {
	lowest := current
	for i := 0; i < 3; i++ {
		for j := 0; j < N_FLOORS; j++ {
			if orderCostList[i][j].Cost != 0 && j+1 < lowest {
				lowest = j+1
			}
		}
	}
	return lowest
}




func printOrders(test *OrderList) {
	fmt.Printf("|FLOOR|   |UP|  |DOWN|  |INSIDE|  |COST|\n")
	temp := 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			temp++
			if temp%3 == 0 {
				fmt.Printf("   %v                X                %v \n", i+1,
					Test.Orders[j][i].Cost)
			}
			if temp%3 == 1 {
				fmt.Printf("   %v                        X        %v \n", i+1,
					Test.Orders[j][i].Cost)
			}
			if temp%3 == 2 {
				fmt.Printf("   %v       X                         %v \n", i+1,
					Test.Orders[j][i].Cost)
			}

		}
		fmt.Printf("\n")
	}
}

