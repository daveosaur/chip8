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
	I     uint16     //16bit register
	Dreg  uint8      //delay register
	Sreg  uint8      //sound register
	PC    uint16     //program counter
	SP    uint8      //stack pointer
	Stack [16]uint16 //the stack!
	Mem   [4096]byte //memory!
}

// types
// funcs
func initState() *State {
	return &State{}
}

func (s *State) update() {

}

func (s *State) draw() {

}

func main() {
	rl.SetTargetFPS(60)
	rl.InitWindow(winX, winY, "chip8")

	s := initState()
	rom, err := os.ReadFile("roms/logo.ch8")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%b", rom)
	s.decodeInstruction(uint16(rom[0])<<8 | uint16(rom[1]))

	os.Exit(0)

	for !rl.WindowShouldClose() {
		s.update()
		s.draw()
	}

}
