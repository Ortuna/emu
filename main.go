package main

import (
  _ "log"
  "image/color"
	"github.com/hajimehoshi/ebiten"
  "github.com/ortuna/emu/cpu"
)

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() { return nil }

  c := color.RGBA{ uint8(rand.Intn(30) + 225), uint8(rand.Intn(50) + 100), 0, 255 }

  screen.Set(32, 16, c)
  return nil
}

func main() {
//	if err := ebiten.Run(update, 64, 32, 6, "Hello, World!"); err != nil {
//		log.Fatal(err)
//}

  cpu := cpu.LoadRom("rom.ch8")

  for i := 0; i < cpu.RomSize / 2; i++ {
    cpu.Tick()
  }
//  cpu.Debug()
}
