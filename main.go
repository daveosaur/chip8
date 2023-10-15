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
}

// funcs
func initState() *State {
	return &State{
		PC: 512, //set program counter to start of program memory
	}
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
		//get x/y of the byte of pixels
		x := i % 8
		y := i / 8
		//loop through the 8 bits in each byte that represent every pixel
		for j := 0; j < 8; j++ {
			//mask off each byte
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

	s := initState()
	s.loadChars()
	rom, err := os.ReadFile("roms/test_.ch8")
	if err != nil {
		panic(err)
	}
	//copy the loaded rom into the emulator
	romLength := copy(s.Mem[512:], rom)
	_ = romLength

	for !rl.WindowShouldClose() {
		err := s.update()
		if err != nil {
			panic(err)
		}
		// s.randomizeVmem()
		err = s.draw()
		if err != nil {
			panic(err)
		}
	}
}
