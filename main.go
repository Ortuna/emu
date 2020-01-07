package main

import (
  "log"
  "os"
  "image/color"
  "math/rand"
	"github.com/hajimehoshi/ebiten"
  "github.com/ortuna/emu"
)

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() { return nil }

  c := color.RGBA{ uint8(rand.Intn(30) + 225), uint8(rand.Intn(50) + 100), 0, 255 }

  screen.Set(32, 16, c)
  return nil
}

func main() {
	if err := ebiten.Run(update, 64, 32, 6, "Hello, World!"); err != nil {
		log.Fatal(err)
	}

}
