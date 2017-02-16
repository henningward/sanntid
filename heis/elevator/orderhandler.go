package elevator

import (
	"../driver"
	"fmt"
	"math/rand"
	"time"
)

const N_FLOORS = 4

const (
	NONE driver.Direction = iota
	UP
	DOWN
)

var orderlist [3][N_FLOORS]order
var orderCostList [3][N_FLOORS]order

type order struct {
	button driver.Button
	cost   int
}

func SetOrder(buttonChan chan driver.Button) {
	var newButton driver.Button

	for {
		select {
		case newButton = <-buttonChan:
			dir, floor := newButton.Dir, newButton.Floor
			orderlist[dir][floor-1] = order{newButton, 100000}

			fmt.Println("\n\n\n\n")
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func ComputeCost(floorChan chan driver.FloorStatus, motorDir *driver.Direction) {
	var newFloorStatus driver.FloorStatus

	for {
		select {
		case newFloorStatus = <-floorChan:
			current, prev, atFloor := newFloorStatus.CurrentFloor, newFloorStatus.PrevFloor, newFloorStatus.AtFloor
			_ = current
			_ = prev
			_ = atFloor
			for i := 0; i < 3; i++ {
				for j := 0; j < N_FLOORS; j++ {
					orderCostList[i][j].button = (orderlist[i][j]).button
					orderCostList[i][j].cost = rand.Intn(10)
				}
			}

			fmt.Printf("moving %v", int(*motorDir))

		case <-time.After(10 * time.Millisecond):
		}

	}

}

func (order) executeOrder() {
	var toExecute order
	for {
		toExecute.cost = 100000
		for i := 0; i < 3; i++ {
			for j := 0; j < N_FLOORS; j++ {
				if orderCostList[i][j].cost < toExecute.cost {
					toExecute = orderCostList[i][j]
				}
			}
		}
	}

}
