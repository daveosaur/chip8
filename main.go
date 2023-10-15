package main

//build without console:
//go build -ldflags -H=windowsgui

import (
	"fmt"

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
	Dreg  uint8      //delay register
	Sreg  uint8      //sound register
	PC    uint16     //program counter
	SP    uint8      //stack pointer
	Stack [16]uint16 //the stack!
	Mem   [4096]byte //memory!
}

// funcs
func initState() *State {
	return &State{
		PC: 512, //set program counter to start of program memory
	}
}

func (s *State) update() error {
	f, err := s.decodeInstruction()
	if err != nil {
		return err
	}
	f()
	s.PC += 2

	return nil
}

func main() {
	rl.SetTargetFPS(60)
	rl.InitWindow(winX, winY, "chip8")

	s := initState()
	s.loadChars()
	rom, err := os.ReadFile("roms/particle.ch8")
	if err != nil {
		panic(err)
	}
	//copy the loaded rom into the emulator
	romLength := copy(s.Mem[512:], rom)
	i := 0

	for !rl.WindowShouldClose() {
		err := s.update()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("did instruction", i)
		i += 2
		if i > romLength {
			os.Exit(0)
		}
	}

}
