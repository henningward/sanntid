package driver

const MOTOR_SPEED 2800
const MOTOR_STOP

func Init() {
	init_status := io_init()
	if !init_status {
		fmt.Printf("Unable to initialize elevator hardware! \n")
	}

	motorChan := make(chan Direction)
	go setMotorDirection()
	// Alt av buttons og lys må nullstilles her
}


type Direction int 
/*
var motorChan chan Direction
var motorDir Direction
*/

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
	LIGHT_UP1 : {UP, 1}
	LIGHT_UP2 : {UP, 2}
	LIGHT_UP3 : {UP, 3}

	LIGHT_DOWN2 : {DOWN, 2}
	LIGHT_DOWN3 : {DOWN, 3}
	LIGHT_DOWN4 : {DOWN, 4}

	LIGHT_COMMAND1 : {NONE, 1}
	LIGHT_COMMAND2 : {NONE, 2}
	LIGHT_COMMAND3 : {NONE, 3}
	LIGHT_COMMAND4 : {NONE, 4}
}




func setMotorDirection(){
	//hva gjør clear/setbit av MOTORDIR?
	
	for{
		motorDir <- motorChan
		if (motorDir == NONE){
			io_write_analog(MOTOR, NONE)

	} else if (motorDir == UP){
		io_clear_bit(MOTORDIR)
		io_write_analog(MOTOR, UP)

	} else if (motorDir == DOWN){
		io_set_bit(MOTORDIR)
		io_write_analog(MOTOR, DOWN)

	}
	}
}

func setButtonLamp(btn button, value int){

	if (btn.Direction == 0){
		//floor og value må settes
	}
}