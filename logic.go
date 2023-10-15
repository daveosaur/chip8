package main

import (
	"errors"
	"fmt"
	"math/rand"
	// rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	MAX_STACK   = 16
	PROG_OFFSET = 0x200
)

// global things. errors?
var (
	tooStacked     = errors.New("overflow stack!")
	noStack        = errors.New("underflow stack!!!")
	badInstruction = errors.New("invalid instruction??")
)

func (s *State) push(ins uint16) error {
	if s.SP >= MAX_STACK {
		return tooStacked
	}
	if s.Stack[s.SP] != 0 {
		s.SP++
	}
	s.Stack[s.SP] = ins

	return nil
}

func (s *State) pop() (uint16, error) {
	if s.SP < 0 {
		return 0, noStack
	}
	popped := s.Stack[s.SP]
	s.Stack[s.SP] = 0
	if s.SP > 0 {
		s.SP--
	}
	return popped, nil
}

func (s *State) increment() error {
	s.PC += 2
	return nil
}

func (s *State) decodeInstruction() error {
	//instruction is the 2 bytes after program counter
	//decode the 4 bits into a,b,c,d
	a := (s.Mem[s.PC] & 0xF0 >> 4)
	b := (s.Mem[s.PC] & 0x0F)
	c := (s.Mem[s.PC+1] & 0xF0 >> 4)
	d := (s.Mem[s.PC+1] & 0x0F)

	// debug instruction prints
	// fmt.Println(a, b, c, d)
	// fmt.Printf("%x\n", uint16(uint16(a)<<12|uint16(b)<<8|uint16(c)<<4|uint16(d)))

	// keys:
	// a, b, c, d bytes 0-3 (left to right)
	// x = byte 1
	// y = byte 2
	// kk = bytes 2,3
	// nnn = bytes 1,2,3
	switch a {
	case 0:
		switch d {
		case 0x0:
			// fmt.Println("clear")
			for i := range s.Vmem {
				s.Vmem[i] = 0
			}
		case 0xE:
			// fmt.Println("return from subroutine")
			// probably considered a jump for return reasons
			ret, err := s.pop()
			if err != nil {
				return err
			}
			s.PC = ret
			return nil
		}
	case 1:
		// fmt.Println("jump nnn")
		addr := uint16(uint16(b)<<8 | uint16(c)<<4 | uint16(d))
		s.PC = addr
		return nil
	case 2:
		// fmt.Println("call nnn")
		// pushes current PC to stack, then jumps to nnn
		err := s.push(s.PC)
		if err != nil {
			return err
		}
		addr := uint16(uint16(b)<<8 | uint16(c)<<4 | uint16(d))
		s.PC = addr
		return nil
	case 3:
		// fmt.Println("skip if vx equals kk")
		if s.Reg[b] == uint8(uint8(c)<<4|uint8(d)) {
			s.increment()
		}
	case 4:
		// fmt.Println("skip if vx not equals kk")
		if s.Reg[b] != uint8(uint8(c)<<4|uint8(d)) {
			s.increment()
		}
	case 5:
		//5xy0
		if d == 0 {
			// fmt.Println("skip if vx == vy")
			if s.Reg[b] == s.Reg[c] {
				s.increment()
			}
		}
	case 6:
		//6xkk
		// fmt.Println("set vx to kk")
		val := uint8(uint8(c)<<4 | uint8(d))
		s.Reg[b] = val
	case 7:
		// fmt.Println("add kk to vx")
		val := uint8(uint8(c)<<4 | uint8(d))
		s.Reg[b] += val
	case 8:
		switch d {
		case 0:
			//vx = vy
			s.Reg[b] = s.Reg[c]
		case 1:
			//vx = vx OR vy
			s.Reg[b] = s.Reg[b] | s.Reg[c]
		case 2:
			//vx = vx AND vy
			s.Reg[b] = s.Reg[b] & s.Reg[c]
		case 3:
			//vx = vx XOR vy
			s.Reg[b] = s.Reg[b] ^ s.Reg[c]
		case 4:
			//vx = vx + vy
			//if result > 255, vf = 1, else 0 (carry bit)
			val := int(s.Reg[b]) + int(s.Reg[c])
			if val > 255 {
				s.Reg[15] = 1
			} else {
				s.Reg[15] = 0
			}
			s.Reg[b] = uint8(val)
		case 5:
			//vx = vx - vy
			//if vx > vy, vf = 1. else 0
			if s.Reg[b] > s.Reg[c] {
				s.Reg[15] = 1
			} else {
				s.Reg[15] = 0
			}
			s.Reg[b] -= s.Reg[c]
		case 6:
			//vx = vx >> 1
			//if least-significant bit of vx is 1, then vf = 1. else 0.
			check := s.Reg[b] & 1 //pull the least significant bit
			s.Reg[15] = check
			s.Reg[b] = s.Reg[b] >> 1
		case 7:
			//vx = vy - vx
			//if vy > vx, vf = 1. else 0.
			if s.Reg[c] > s.Reg[b] {
				s.Reg[15] = 1
			} else {
				s.Reg[15] = 0
			}
			s.Reg[b] = s.Reg[c] - s.Reg[b]
		case 0xE:
			//vx = vx << 1
			//if most significant bit of fx is 1, then fv = 1. else 0.
			check := (s.Reg[b] & 128) >> 7 //pull most significant bit
			s.Reg[15] = check
			s.Reg[b] = s.Reg[b] << 1
		}
	case 9:
		//skip next instruction if vx != vy
		if s.Reg[b] != s.Reg[c] {
			s.increment()
		}
	case 0xA:
		//set I = nnn
		nnn := uint16(uint16(b)<<8 | uint16(c)<<4 | uint16(d))
		s.Ireg = nnn
	case 0xB:
		//jump to nnn + v0
		nnn := uint16(uint16(b)<<8 | uint16(c)<<4 | uint16(d))
		s.PC = nnn + uint16(s.Reg[0])
		return nil
	case 0xC:
		//set vx random byte & kk
		num := uint8(rand.Int() % 256)
		kk := (c << 4) | d
		s.Reg[b] = num & kk
	case 0xD:
		//draw/collision thing
		//TODO:
		xPos := s.Reg[b]
		yPos := s.Reg[c]
		for i := 0; i < int(d); i++ {
			curByte := s.Mem[s.Ireg+uint16(i)]
			rem := xPos % 8
			pos := (yPos+byte(i))*8 + rem
			if rem == 0 { //does it line up on the byte and make things easy??
				s.Vmem[pos] = s.Vmem[pos] ^ curByte
			} else { //nope
				nextByte := curByte << (8 - rem)
				curByte = curByte >> rem
				s.Vmem[pos] = s.Vmem[pos] ^ curByte
				s.Vmem[pos+1] = s.Vmem[pos+1] ^ nextByte
			}
		}

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
			s.Reg[b] = s.Dreg
		case [2]byte{0x0, 0xA}:
			//wait for key press, store in VX
			//full breakpoint
			//TODO
		case [2]byte{0x1, 0x5}:
			//DT = vx
			s.Dreg = s.Reg[b]
		case [2]byte{0x1, 0x8}:
			//ST = vx
			s.Sreg = s.Reg[b]
		case [2]byte{0x1, 0xE}:
			//set I = I + vx
			s.Ireg = s.Ireg + uint16(s.Reg[b])
		case [2]byte{0x2, 0x9}:
			//set I = location of sprite for digit vx
			//assuming sprites are stored 0x0 and 5 bytes each
			s.Ireg = uint16(5 * s.Reg[b])
		case [2]byte{0x3, 0x3}:
			//store vx BCD in I, I+1, I+2
			//extract digits
			num := s.Reg[b]
			hund := num / 100 //hundreds digit
			num -= hund * 100
			ten := num / 10 //tens digit
			num -= ten * 10 //ones
			s.Mem[s.Ireg] = hund
			s.Mem[s.Ireg+1] = ten
			s.Mem[s.Ireg+2] = num
		case [2]byte{0x5, 0x5}:
			//copy values v0 - vx to memory location I
			for i := 0; i < int(b); i++ {
				s.Mem[int(s.Ireg)+i] = s.Reg[i]
			}
		case [2]byte{0x6, 0x5}:
			//read values from I into registers v0 - vx
			for i := 0; i < int(b); i++ {
				s.Reg[i] = s.Mem[int(s.Ireg)+i]
			}
		}
	default:
		fmt.Println("bad instruction")
		return badInstruction
	}
	//increment PC
	//use return in switch for jumps to avoid increment
	s.increment()

	return nil
}

func (s *State) randomizeVmem() {
	for i := range s.Vmem {
		s.Vmem[i] = uint8(rand.Int() % 256)
	}
}
