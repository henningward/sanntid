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
	
	for _, val := range setButtonLightMap{
		io_clear_bit(val)

	}

	go setMotorDirection()

	go listenButton()
	//testing motor direction function
	go func () {
	for {
		//motorChan <- UP
	}
	}()

	

	time.Sleep(5*time.Second)
	io_write_analog(MOTOR, int(NONE))

}

type Direction int 

var motorChan chan Direction
var motorDir Direction



const (
	NONE Direction = iota
	UP
	DOWN
)
type button struct{
	dir Direction
	floor int
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




func listenButton(){
	for{
	for key, val := range buttonMap{
		if (io_read_bit(key) != 0){
			setButtonLamp(val, 1)
		}

	}
	}
}
	