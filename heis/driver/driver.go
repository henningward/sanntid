package driver
import "fmt"
import "time"
const MOTOR_SPEED = 2800
const MOTOR_STOP = 0


func Init() {
	init_status := io_init()
	if (init_status == 0) {
		fmt.Printf("Unable to initialize elevator hardware! \n")
	}

	motorChan = make(chan Direction) //Hva er greia med make? og deklarere utenfor funksjonene?
	buttonChan = make(chan button)
	floorChan = make(chan floorStatus)
	for _, val := range setButtonLightMap{
		io_clear_bit(val)

	}

	go func(){ //Bare for å foreløpig lese ut av buttonChan
		for{
			dir, floor := GetButton()
			newButton := button{dir, floor}
			fmt.Println("\n New button pressed: ")
			fmt.Println(newButton)			
		}
	}()

	go func(){ //Bare for å foreløpig lese ut av buttonChan
		for{
			current, prev, at := GetFloor()
			newFloorStat := floorStatus{current, prev, at}
			fmt.Println("\n New floor status: ")
			fmt.Println(newFloorStat)
			
		}
	}()




	go setMotorDirection()

	go listenButton()
	//testing motor direction function
	go listenFloor()
	
	MotorDOWN()

	time.Sleep(1*time.Second)

	MotorIDLE()


	time.Sleep(5*time.Second)
	io_write_analog(MOTOR, int(NONE))

}

type Direction int 
type Floor int
var motorChan chan Direction
var motorDir Direction
var buttonChan chan button
var floorChan chan floorStatus

const (
	NONE Direction = iota
	UP
	DOWN
)
type button struct{
	dir Direction
	floor Floor
}

type floorStatus struct{
	currentFloor Floor
	prevFloor Floor
	atFloor bool
}

var buttonLightMap = map[int] button {
	LIGHT_UP1 : {UP, 1},
	LIGHT_UP2 : {UP, 2},
	LIGHT_UP3 : {UP, 3},

	LIGHT_DOWN2 : {DOWN, 2},
	LIGHT_DOWN3 : {DOWN, 3},
	LIGHT_DOWN4 : {DOWN, 4},

	LIGHT_COMMAND1 : {NONE, 1},
	LIGHT_COMMAND2 : {NONE, 2},
	LIGHT_COMMAND3 : {NONE, 3},
	LIGHT_COMMAND4 : {NONE, 4},
}

var setButtonLightMap = map[button] int {
	{UP, 1}   : LIGHT_UP1,
	{UP, 2}	  : LIGHT_UP2,
	{UP, 3}	  : LIGHT_UP3,

	{DOWN, 2} : LIGHT_DOWN2,
	{DOWN, 3} : LIGHT_DOWN3,
	{DOWN, 4} : LIGHT_DOWN4,

	{NONE, 1} : LIGHT_COMMAND1,
	{NONE, 2} : LIGHT_COMMAND2,
	{NONE, 3} : LIGHT_COMMAND3,
	{NONE, 4} : LIGHT_COMMAND4,
}




func setMotorDirection(){

	//hva gjør clear/setbit av MOTORDIR?
	for{
		motorDir := <- motorChan
		if (motorDir == NONE){
		io_write_analog(MOTOR, int(NONE))
		} else if (motorDir == UP){
		io_clear_bit(MOTORDIR)
		io_write_analog(MOTOR, int(MOTOR_SPEED))
		} else if (motorDir == DOWN){
		io_set_bit(MOTORDIR)
		io_write_analog(MOTOR, int(MOTOR_SPEED))

	}
	}
}

func setButtonLamp(btn button, value int){
	if (value != 0){
		io_set_bit(setButtonLightMap[btn])
	} else {
		io_clear_bit(setButtonLightMap[btn])
	}
}

func MotorUP(){
	motorChan <- UP
}

func MotorDOWN(){
	motorChan <- DOWN
}

func MotorIDLE(){
	motorChan <- NONE
}


func listenButton(){
	var buttonMap = map[int] button {
	BUTTON_UP1 : {UP, 1},
	BUTTON_UP2 : {UP, 2},
	BUTTON_UP3 : {UP, 3},

	BUTTON_DOWN2 : {DOWN, 2},
	BUTTON_DOWN3 : {DOWN, 3},
	BUTTON_DOWN4 : {DOWN, 4},

	BUTTON_COMMAND1 : {NONE, 1},
	BUTTON_COMMAND2 : {NONE, 2},
	BUTTON_COMMAND3 : {NONE, 3},
	BUTTON_COMMAND4 : {NONE, 4},
}
	buttonsPressed := make(map[int]bool)
	for key, _ := range buttonMap{
		buttonsPressed[key] = (io_read_bit(key)!=0)

	}

	for{
	for key, val := range buttonMap{
		if (io_read_bit(key) != 0 && !buttonsPressed[key]){
			newButton := val
			buttonsPressed[key] = true
			setButtonLamp(val, 1)
			buttonChan <- newButton			
		}

	}
	}
}
func GetButton() (Direction, Floor){
	newButton := <- buttonChan
	return newButton.dir, newButton.floor
}

func GetFloor() (Floor, Floor, bool){
	newFloor := <- floorChan
	return newFloor.currentFloor, newFloor.prevFloor, newFloor.atFloor
}

func listenFloor(){
	var floorMap = map[int] Floor {
	SENSOR_FLOOR1 : 1,
	SENSOR_FLOOR2 : 2,
	SENSOR_FLOOR3 : 3,
	SENSOR_FLOOR4 : 4,
} 
	var currentFloor Floor = 0
	var prevFloor Floor = -1

	atFloor:= make(map[int]bool)
	for key, _ := range floorMap{
		atFloor[key] = (io_read_bit(key) !=0)
		println(atFloor[key])
	}

	for{
		for key, val := range floorMap{
			if (io_read_bit(key) != 0 && !atFloor[key]){
				prevFloor = currentFloor
				atFloor[key] = true
				currentFloor = val
				newFloorStatus := floorStatus{currentFloor, prevFloor, atFloor[key]}
				floorChan <- newFloorStatus
			}
			if (io_read_bit(key) == 0 && atFloor[key]){
				atFloor[key] = false
				prevFloor = currentFloor
				newFloorStatus := floorStatus{currentFloor, prevFloor, atFloor[key]}
				floorChan <- newFloorStatus
				
			}
	
			
		}

	}

}