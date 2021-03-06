package elevator

import (
	"../driver"
	"fmt"
	"time"
)

const N_FLOORS = 4

const (
	NONE driver.Direction = iota
	UP
	DOWN
)

var toExecute Order

var ElevatorMsg OrderMsg

var timerOwnOrders time.Time

type Order struct {
	Button driver.Button
	Cost   int
}

func SetOrder(buttonChan chan driver.Button, newOrders *OrderList) {

	var newButton driver.Button
	for {
		newButton = <-buttonChan
		timerOwnOrders = time.Now()
		dir, floor := newButton.Dir, newButton.Floor
		newOrders[dir][floor-1] = Order{newButton, 100000}
		time.Sleep(10 * time.Millisecond)
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
							*cost = 1000*(orderFloor-currentFloor) + ID
						} else {
							if highestFloorOrder(currentFloor, orderCostList) >= orderFloor {
								*cost = 1000*(2*highestFloorOrder(currentFloor, orderCostList)-orderFloor-currentFloor) + ID + 1

							} else {
								*cost = 1000*(orderFloor-currentFloor) + ID + 1

							}

						}
					} else {
						if orderDir != UP {
							*cost = 1000*(2*highestFloorOrder(currentFloor, orderCostList)-orderFloor-currentFloor) + ID + 1

						} else {
							if lowestFloorOrder(currentFloor, orderCostList) <= orderFloor {

								*cost = 1000*(2*highestFloorOrder(currentFloor, orderCostList)-2*lowestFloorOrder(currentFloor, orderCostList)+orderFloor-currentFloor) + ID

							} else {
								*cost = 1000*(2*highestFloorOrder(currentFloor, orderCostList)-currentFloor-orderFloor) + ID

							}
						}
					}
				} else if elev.Dir == DOWN {

					if orderFloor <= currentFloor {
						if orderDir != UP {
							*cost = 1000*(currentFloor-orderFloor) + ID

						} else {
							if lowestFloorOrder(currentFloor, orderCostList) <= orderFloor {
								*cost = 1000*(-2*lowestFloorOrder(currentFloor, orderCostList)+orderFloor+currentFloor) + ID + 1

							} else {
								*cost = 1000*(currentFloor-orderFloor) + ID + 1

							}

						}
					} else {
						if orderDir != DOWN {
							*cost = 1000*(-2*lowestFloorOrder(currentFloor, orderCostList)+orderFloor+currentFloor) + ID + 1

						} else {
							if highestFloorOrder(currentFloor, orderCostList) >= orderFloor {
								*cost = 1000*(2*highestFloorOrder(currentFloor, orderCostList)-2*lowestFloorOrder(currentFloor, orderCostList)-orderFloor+currentFloor) + ID

							} else {
								*cost = 1000*(-2*lowestFloorOrder(currentFloor, orderCostList)+currentFloor+orderFloor) + ID

							}
						}
					}
				} else if elev.Dir == NONE {
					if elev.STATE == "IDLE" {
						if time.Since(elev.StateTimer) > 500*time.Millisecond {
							if currentFloor > orderFloor {
								*cost = 1000*(currentFloor-orderFloor) + ID

							} else {
								*cost = 1000*(orderFloor-currentFloor) + ID
							}

						}
					}

				} else {
					fmt.Printf("failed to compute cost at state %v", elev.STATE)
				}

				if elev.STATE == "DOORS OPEN" {
					if currentFloor == orderFloor {
						clearOrdersSameFloor(currentFloor, orderCostList, newOrders)

					}
				}
			}

		}
	}
}

func clearOrdersSameFloor(currentFloor int, orderCostList *OrderList, newOrders *OrderList) {
	for i := 0; i < 3; i++ {
		for j := 0; j < N_FLOORS; j++ {
			if int(orderCostList[i][j].Button.Floor) == currentFloor {
				DeleteOrder(orderCostList[i][j], orderCostList, newOrders)
			}
		}
	}

}

func ExecuteOrder(executeOrderChan chan Order, orderCostList *OrderList) {
	toExecute.Cost = 100000
	lastExecute := toExecute
	for {
		if time.Since(timerOwnOrders) < 1000*time.Millisecond {
			time.Sleep(1000 * time.Millisecond)
		}
		if time.Since(timerRecOrders) < 1000*time.Millisecond {
			time.Sleep(1000 * time.Millisecond)
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < N_FLOORS; j++ {
				if (orderCostList[i][j].Cost < toExecute.Cost) && orderCostList[i][j].Cost != 0 {
					toExecute = orderCostList[i][j]
				}
				if toExecute.Button == orderCostList[i][j].Button && toExecute.Cost < orderCostList[i][j].Cost {
					toExecute.Cost = orderCostList[i][j].Cost
				}

			}
			time.Sleep(10 * time.Millisecond)
		}
		if toExecute.Cost < 100000 && toExecute != lastExecute {
			executeOrderChan <- toExecute
			lastExecute = toExecute
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func DeleteOrder(order Order, orderCostList *OrderList, newOrders *OrderList) {
	emptyButton := driver.Button{0, 0}
	emptyOrder := Order{emptyButton, 0}
	driver.ClearButtonPressed(order.Button)
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
				lowest = j + 1
			}
		}
	}
	return lowest
}
