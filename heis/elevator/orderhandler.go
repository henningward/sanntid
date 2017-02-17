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

func ComputeCost(elev ElevState, motorDir *driver.Direction) {
	for {
			for i := 0; i < 3; i++ {
				for j := 0; j < N_FLOORS; j++ {
					if orderlist[i][j].Cost != 0{
						orderCostList[i][j].Button = (orderlist[i][j]).Button
						
						orderCostList[i][j].Cost = 2
						fmt.Println("computing..")
					}
					
				}
			}
			for i := 0; i < 3; i++ {
				for j := 0; j < N_FLOORS; j++ {
					//fmt.Printf("button: %v with cost: %v", orderCostList[i][j].Button,
					//orderCostList[i][j].Cost)
					//fmt.Printf("\n")
				}
			}
		//fmt.Printf("\n\n\n\n\n\n")
		time.Sleep(1*time.Second)
		}
}

func ExecuteOrder(executeOrderChan chan Order) {
	toExecute.Cost = 10000
	for {
		for i := 0; i < 3; i++ {
			for j := 0; j < N_FLOORS; j++ {
				if (orderCostList[i][j].Cost < toExecute.Cost) && orderCostList[i][j].Cost != 0{
					toExecute = orderCostList[i][j]
					executeOrderChan <- toExecute
					println("\n execute order \n")
				}
			}
		}
	}

}


func DeleteOrder(order Order){
	emptyButton := driver.Button{0, 0}
	emptyOrder := Order{emptyButton, 0}
	driver.SetButtonLamp(order.Button, 0)
	orderCostList[order.Button.Dir][int(order.Button.Floor)-1] = emptyOrder
	orderlist[order.Button.Dir][int(order.Button.Floor)-1] = emptyOrder
	toExecute.Cost = 100000

}