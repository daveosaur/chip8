package main

import (
	"fmt"
)

func (s *State) decodeInstruction(inst uint16) {
	a := uint8(inst & 0xF000 >> 12)
	b := uint8(inst & 0x0F00 >> 8)
	c := uint8(inst & 0x00F0 >> 4)
	d := uint8(inst & 0x000F)

	_ = b //TODO: remove

	switch a {
	case 0x0:
		switch d {
		case 0x0:
			fmt.Println("clear")
		case 0xE:
			fmt.Println("return from subroutine")
		}
	case 1:
		fmt.Println("jump nnn")
	case 2:
		fmt.Println("call nnn")
	case 3:
		fmt.Println("skip if vx equals kk")
	case 4:
		fmt.Println("skip if vx not equals kk")
	case 5:
		fmt.Println("skip if vx == vy")
	case 6:
		fmt.Println("set vx to kk")
	case 7:
		fmt.Println("add kk to vx")
	case 8:
		switch d {
		case 0:
		//vx = vy
		case 1:
		//vx = vx OR vy
		case 2:
		//vx = vx AND vy
		case 3:
		//vx = vx XOR vy
		case 4:
		//vx = vx + vy
		//if result > 255, vf = 1, else 0 (carry bit)
		case 5:
		//vx = vx - vy
		//if vx > vy, vf = 1. else 0
		case 6:
		//vx = vx >> 1
		//if least-significant bit of vx is 1, then vf = 1. else 0.
		case 7:
		//vx = vy - vx
		//if vy > vx, vf = 1. else 0.
		case 0xE:
			//vx = vx << 1
			//if most significant bit of fx is 1, then fv = 1. else 0.
		}
	case 9:
		//skip next instruction if vx != vy
	case 0xA:
		//set I = nnn
	case 0xB:
		//jump to nnn + v0
	case 0xC:
		//set vx random byte & kk
	case 0xD:
		//collision thing
	case 0xE:
		switch [2]byte{c, d} {
		case [2]byte{0x9, 0xE}:
			//skip VX
			//skip instruction if key of value vx is pressed
		case [2]byte{0xA, 0x1}:
			//skip !VX
			//skip instruction if key of value vx is not pressed
		}
	case 0xF:
		switch [2]byte{c, d} {
		case [2]byte{0x0, 0x7}:
			//vx = DT
		case [2]byte{0x0, 0xA}:
			//wait for key press, store in VX
			//full breakpoint
		case [2]byte{0x1, 0x5}:
			//DT = vx
		case [2]byte{0x1, 0x8}:
			//ST = vx
		case [2]byte{0x1, 0xE}:
			//set I = I + vx
		case [2]byte{0x2, 0x9}:
			//set I = location of sprite for digit vx
		case [2]byte{0x3, 0x3}:
			//store vx BCD in I, I+1, I+2
		case [2]byte{0x5, 0x5}:
			//copy values v0 - vx to memory location I
		case [2]byte{0x6, 0x5}:
			//read values from I into registers v0 - vx

		}
	default:
		fmt.Println("bad instruction")

	}

}
