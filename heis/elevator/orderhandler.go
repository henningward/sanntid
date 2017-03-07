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

var timerOwnOrders time.Time

type Order struct {
	Button driver.Button
	Cost   int
}

func SetOrder(buttonChan chan driver.Button, newOrders *OrderList) {

	var newButton driver.Button
	for {
		newButton = <-buttonChan
		println("button pressed")
		timerOwnOrders = time.Now()
		dir, floor := newButton.Dir, newButton.Floor
		newOrders[dir][floor-1] = Order{newButton, 100000}
		time.Sleep(100 * time.Millisecond)
	}
}

func ComputeCost(elev *ElevState, motorDir *driver.Direction, orderCostList *OrderList, newOrders *OrderList, ID int) {
	currentFloor := 0
	orderFloor := 0
	orderDir := NONE

	currentFloor = int(elev.FloorStatus.CurrentFloor)
	for i := 0; i < 3; i++ {
		for j := 0; j < N_FLOORS; j++ {
			if newOrders[i][j].Cost != 0 {
				orderCostList[i][j].Button = (newOrders[i][j]).Button
				orderFloor = int(newOrders[i][j].Button.Floor)
				orderDir = newOrders[i][j].Button.Dir
				atFloor := elev.FloorStatus.AtFloor
				cost := &orderCostList[i][j].Cost

				if (orderFloor == currentFloor) && !atFloor {
					*cost = 100000
					break
				}

				if elev.Dir == UP {
					if orderFloor >= currentFloor {
						if orderDir != DOWN {
							*cost = 10*(orderFloor-currentFloor) + ID
						} else {
							if highestFloorOrder(currentFloor, orderCostList) >= orderFloor {
								*cost = 10*(2*highestFloorOrder(currentFloor, orderCostList)-orderFloor-currentFloor) + ID

							} else {
								*cost = 10*(orderFloor-currentFloor) + ID

							}

						}
					} else {
						if orderDir != UP {
							*cost = 10*(2*highestFloorOrder(currentFloor, orderCostList)-orderFloor-currentFloor) + ID

						} else {
							if lowestFloorOrder(currentFloor, orderCostList) <= orderFloor {

								*cost = 10*(2*highestFloorOrder(currentFloor, orderCostList)-2*lowestFloorOrder(currentFloor, orderCostList)+orderFloor-currentFloor) + ID

							} else {
								*cost = 10*(2*highestFloorOrder(currentFloor, orderCostList)-currentFloor-orderFloor) + ID

							}
						}
					}
				} else if elev.Dir == DOWN {

					if orderFloor <= currentFloor {
						if orderDir != UP {
							*cost = 10*(currentFloor-orderFloor) + ID

						} else {
							if lowestFloorOrder(currentFloor, orderCostList) <= orderFloor {
								*cost = 10*(-2*lowestFloorOrder(currentFloor, orderCostList)+orderFloor+currentFloor) + ID

							} else {
								*cost = 10*(currentFloor-orderFloor) + ID

							}

						}
					} else {
						if orderDir != DOWN {
							*cost = 10*(-2*lowestFloorOrder(currentFloor, orderCostList)+orderFloor+currentFloor) + ID

						} else {
							if highestFloorOrder(currentFloor, orderCostList) >= orderFloor {
								*cost = 10*(2*highestFloorOrder(currentFloor, orderCostList)-2*lowestFloorOrder(currentFloor, orderCostList)-orderFloor+currentFloor) + ID

							} else {
								*cost = 10*(-2*lowestFloorOrder(currentFloor, orderCostList)+currentFloor+orderFloor) + ID

							}
						}
					}
				} else if elev.Dir == NONE {
					if elev.STATE == "IDLE" {
						if elev.TimeInState > 500 {
							if currentFloor > orderFloor {
								*cost = 10*(currentFloor-orderFloor) + ID

							} else {
								*cost = 10*(orderFloor-currentFloor) + ID

							}

						}
					}
				} else {
					fmt.Print("failed to compute cost at state %v", elev.STATE)
				}
				//elev.FloorStatus.CurrentFloor

			}

		}
		//printOrders(orderCostList)
		//printOrders(&Test.Orders)
		//fmt.Printf("\n\n\n\n\n\n")

	}
}

func ExecuteOrder(executeOrderChan chan Order, orderCostList *OrderList) {
	toExecute.Cost = 100000

	for {
		if time.Since(timerOwnOrders) < 5000*time.Millisecond {
			println("Sleeping1...")
			time.Sleep(1000 * time.Millisecond)
		}
		if time.Since(timerRecOrders) < 5000*time.Millisecond {
			println("Sleeping2...")
			time.Sleep(1000 * time.Millisecond)
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < N_FLOORS; j++ {
				if orderCostList[i][j].Cost != 0 {
					//println(orderCostList[i][j].Cost)
					//println("\n")
				}
				if (orderCostList[i][j].Cost < toExecute.Cost) && orderCostList[i][j].Cost != 0 {
					toExecute = orderCostList[i][j]
					//executeOrderChan <- toExecute
				}

			}
			time.Sleep(10 * time.Millisecond)
		}

		if toExecute.Cost < 100000 {
			printOrders(orderCostList)
			println(toExecute.Cost)
			println(toExecute.Cost)
			println(toExecute.Cost)
			println(toExecute.Cost)
			println(toExecute.Cost)

			executeOrderChan <- toExecute
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func DeleteOrder(order Order, orderCostList *OrderList, newOrders *OrderList) {
	emptyButton := driver.Button{0, 0}
	emptyOrder := Order{emptyButton, 0}
	driver.SetButtonLamp(order.Button, 0)
	orderCostList[order.Button.Dir][int(order.Button.Floor)-1] = emptyOrder
	newOrders[order.Button.Dir][int(order.Button.Floor)-1] = emptyOrder
	toExecute.Cost = 100000
	//printOrders(orderCostList)

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
				lowest = j + 1
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
