package elevator

import (
	"time"
	"fmt"
)

var recorders [3][N_FLOORS]Order
func ReceiveOrder(msgRecCh chan OrderMsg, elev *ElevState, executeOrderCh chan Order){
		var MsgRec OrderMsg
		for {
		select {
		case MsgRec = <-msgRecCh:
			recorders = MsgRec.Orders
			printOrdersRec(MsgRec)
		case <-time.After(100 * time.Millisecond):

		}

		time.Sleep(100 * time.Millisecond)

	}
}

func printOrdersRec(Test OrderMsg) {
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

