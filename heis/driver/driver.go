package driver

import "fmt"
import "time"

const MOTOR_SPEED = 2800
const MOTOR_STOP = 0

func Init(buttonChan chan Button, floorChan chan FloorStatus, motorDir *Direction) {
	init_status := io_init()
	if init_status == 0 {
		fmt.Printf("Unable to initialize elevator hardware! \n")
	}

	motorChan = make(chan Direction) //Hva er greia med make? og deklarere utenfor funksjonene?

	for _, val := range ButtonToLightMap {
		io_clear_bit(val)
	}
	for _, val := range floorIndMap {
		io_clear_bit(val)
	}

	SetDoorLamp(0)

	go setMotorDirection(motorDir)
	go ListenButton(buttonChan)
	go ListenFloor(floorChan)

	//Moving elevator to closest floor
	for currentFloor == 0 {
		MotorDOWN()
	}

	MotorIDLE()
	for {
		time.Sleep(100 * time.Second)
	}

	//io_write_analog(MOTOR, int(NONE))

}

type Direction int
type Floor int

var motorChan chan Direction
var currentFloor Floor = 0
var currentFloorStatus FloorStatus
var buttonsPressed map[int]bool

const (
	NONE Direction = iota
	UP
	DOWN
)

type Button struct {
	Dir   Direction
	Floor Floor
}

type FloorStatus struct {
	CurrentFloor Floor
	PrevFloor    Floor
	AtFloor      bool
}

var LightToButtonMap = map[int]Button{
	LIGHT_UP1: {UP, 1},
	LIGHT_UP2: {UP, 2},
	LIGHT_UP3: {UP, 3},

	LIGHT_DOWN2: {DOWN, 2},
	LIGHT_DOWN3: {DOWN, 3},
	LIGHT_DOWN4: {DOWN, 4},

	LIGHT_COMMAND1: {NONE, 1},
	LIGHT_COMMAND2: {NONE, 2},
	LIGHT_COMMAND3: {NONE, 3},
	LIGHT_COMMAND4: {NONE, 4},
}

var ButtonToLightMap = map[Button]int{
	{UP, 1}: LIGHT_UP1,
	{UP, 2}: LIGHT_UP2,
	{UP, 3}: LIGHT_UP3,

	{DOWN, 2}: LIGHT_DOWN2,
	{DOWN, 3}: LIGHT_DOWN3,
	{DOWN, 4}: LIGHT_DOWN4,

	{NONE, 1}: LIGHT_COMMAND1,
	{NONE, 2}: LIGHT_COMMAND2,
	{NONE, 3}: LIGHT_COMMAND3,
	{NONE, 4}: LIGHT_COMMAND4,
}

var ButtonToIntMap = map[Button]int{
	{UP, 1}: BUTTON_UP1,
	{UP, 2}: BUTTON_UP2,
	{UP, 3}: BUTTON_UP3,

	{DOWN, 2}: BUTTON_DOWN2,
	{DOWN, 3}: BUTTON_DOWN3,
	{DOWN, 4}: BUTTON_DOWN4,

	{NONE, 1}: BUTTON_COMMAND1,
	{NONE, 2}: BUTTON_COMMAND2,
	{NONE, 3}: BUTTON_COMMAND3,
	{NONE, 4}: BUTTON_COMMAND4,
}

var floorIndMap = map[Floor]int{
	1: LIGHT_FLOOR_IND2,
	2: LIGHT_FLOOR_IND1,
	//	LIGHT_FLOOR_IND3: 3,
	//	LIGHT_FLOOR_IND4: 4, // DETTE FUNGERER IKKE...
}

func setMotorDirection(motorDir *Direction) {

	//hva gjør clear/setbit av MOTORDIR?
	for {
		*motorDir = <-motorChan
		if *motorDir == NONE {
			io_write_analog(MOTOR, int(NONE))
		} else if *motorDir == UP {
			io_clear_bit(MOTORDIR)
			io_write_analog(MOTOR, int(MOTOR_SPEED))
		} else if *motorDir == DOWN {
			io_set_bit(MOTORDIR)
			io_write_analog(MOTOR, int(MOTOR_SPEED))

		}
	}
}

func SetButtonLamp(btn Button, value int) {
	if value != 0 {
		io_set_bit(ButtonToLightMap[btn])
	} else {
		io_clear_bit(ButtonToLightMap[btn])
		buttonsPressed[ButtonToIntMap[btn]] = false
	}

}

func SetDoorLamp(val int) {
	if val != 0 {
		io_set_bit(LIGHT_DOOR_OPEN)
	} else {
		io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func setFloorLamp(floor int) {
	switch floor{
		case 1:
	io_clear_bit(LIGHT_FLOOR_IND1)
	io_clear_bit(LIGHT_FLOOR_IND2)
case 2:
	io_set_bit(LIGHT_FLOOR_IND2)
	io_clear_bit(LIGHT_FLOOR_IND1)
case 3:
	io_set_bit(LIGHT_FLOOR_IND1)
	io_clear_bit(LIGHT_FLOOR_IND2)
case 4:
	io_set_bit(LIGHT_FLOOR_IND1)
	io_set_bit(LIGHT_FLOOR_IND2)
	}


}

func MotorUP() {
	motorChan <- UP
}

func MotorDOWN() {
	motorChan <- DOWN
}

func MotorIDLE() {
	motorChan <- NONE
}

func GetFloor(floorChan chan FloorStatus) FloorStatus {
	select {
	case currentFloorStatus = <-floorChan:

	case <-time.After(10 * time.Millisecond):
	}
	return currentFloorStatus
}

func ListenButton(buttonChan chan Button) {
	var intToButtonMap = map[int]Button{
		BUTTON_UP1: {UP, 1},
		BUTTON_UP2: {UP, 2},
		BUTTON_UP3: {UP, 3},

		BUTTON_DOWN2: {DOWN, 2},
		BUTTON_DOWN3: {DOWN, 3},
		BUTTON_DOWN4: {DOWN, 4},

		BUTTON_COMMAND1: {NONE, 1},
		BUTTON_COMMAND2: {NONE, 2},
		BUTTON_COMMAND3: {NONE, 3},
		BUTTON_COMMAND4: {NONE, 4},
	}

	buttonsPressed = make(map[int]bool)
	for key, _ := range intToButtonMap {
		buttonsPressed[key] = (io_read_bit(key) != 0)

	}

	for {
		for key, val := range intToButtonMap {
			if io_read_bit(key) != 0 && !buttonsPressed[key] {
				newButton := val
				buttonsPressed[key] = true
				SetButtonLamp(val, 1)
				buttonChan <- newButton
			}

		}
	}
}

func ListenFloor(floorChan chan FloorStatus) {

	var floorMap = map[int]Floor{
		SENSOR_FLOOR1: 1,
		SENSOR_FLOOR2: 2,
		SENSOR_FLOOR3: 3,
		SENSOR_FLOOR4: 4,
	}

	var prevFloor Floor = -1
	var AtFloor = make(map[int]bool)
	for key, _ := range floorMap {
		//AtFloor[key] = (io_read_bit(key) != 0)
		AtFloor[key] = false
	}

	for {
		for key, val := range floorMap {
			if io_read_bit(key) != 0 && !AtFloor[key] {
				prevFloor = currentFloor
				AtFloor[key] = true
				currentFloor = val
				setFloorLamp(int(currentFloor))
				newFloorStatus := FloorStatus{currentFloor, prevFloor, AtFloor[key]}
				floorChan <- newFloorStatus
			}
			if io_read_bit(key) == 0 && AtFloor[key] {
				AtFloor[key] = false
				prevFloor = currentFloor
				newFloorStatus := FloorStatus{currentFloor, prevFloor, AtFloor[key]}
				floorChan <- newFloorStatus

			}

		}
	}

}
