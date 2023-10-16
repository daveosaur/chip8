package main

//build without console:
//go build -ldflags -H=windowsgui

import (
	// "fmt"

	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	//64x32 is the chip8 resolution
	winX = 640
	winY = 320
)

type State struct {
	Reg   [16]byte   //general registers
	Ireg  uint16     //16bit register
	Dreg  uint8      //delay timer register
	Sreg  uint8      //sound timer register
	PC    uint16     //program counter
	SP    uint8      //stack pointer
	Stack [16]uint16 //the stack!
	Mem   [4096]byte //memory!
	Vmem  [256]byte  //vram!
	Inp   byte       //current input!
}

// funcs

// create a new emulator state, loading rom into memory from path argument
func initState(path string) (*State, error) {
	var state State = State{
		PC: 512, //set program counter to start of program memory
	}
	loadChars(&state)
	rom, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	//copy the loaded rom into the emulator
	copy(state.Mem[512:], rom)

	return &state, nil
}

func (s *State) update() error {
	err := s.decodeInstruction()
	if err != nil {
		return err
	}
	return nil
}

func (s *State) draw() error {
	rl.BeginDrawing()
	rl.ClearBackground(rl.White)
	for i := range s.Vmem {
		// get x/y of the byte of pixels
		// screen is 8 bytes wide, so 64 total. we're drawing 8 pixels at a time tho, so %8.
		// vram is a flat array, so find the Y level using truncated division
		x := i % 8
		y := i / 8
		//loop through the 8 bits in each byte that represent every pixel
		for j := 0; j < 8; j++ {
			//mask off each byte
			//shift down the current drawing bit to the end, and mask everything off with & 1
			mask := (s.Vmem[i] >> (7 - j)) & 1
			if mask == 1 {
				rl.DrawRectangle(int32((80*x)+(10*j)), int32(y*10), 10, 10, rl.Black)
			}
		}
	}
	rl.EndDrawing()
	return nil
}

func main() {
	rl.SetTargetFPS(60)
	rl.InitWindow(winX, winY, "chip8")

	// s, err := initState("roms/tests/6-keypad.ch8")
	// s, err := initState("roms/tests/3-corax+.ch8")
	s, err := initState("roms/tests/4-flags.ch8")
	// s, err := initState("roms/tests/5-quirks.ch8")
	// s, err := initState("roms/framed2.ch8")
	if err != nil {
		panic(err)
	}

	for !rl.WindowShouldClose() {
		for i := 0; i < 4; i++ {
			s.Inp = getInput()
			if err := s.update(); err != nil {
				panic(err)
			}
		}
		if err := s.draw(); err != nil {
			panic(err)
		}
	}
}
