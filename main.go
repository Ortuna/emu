package main

import (
  "log"
	"github.com/hajimehoshi/ebiten"
  "github.com/ortuna/emu/cpu"
)

var c *cpu.Cpu

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() { return nil }
  if c == nil { 
    c = cpu.NewCpu("maze.rom", screen)
    ebiten.SetRunnableInBackground(true)
//    ebiten.SetMaxTPS(6)
  }

  c.Tick()
  //c.DebugTick()
  //panic("out")
  return nil
}

func main() {
  if err := ebiten.Run(update, 64, 32, 6, "Hello, World!"); err != nil {
    log.Fatal(err)
  }
}
