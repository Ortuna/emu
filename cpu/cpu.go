package cpu

import (
  "math/rand"
  "fmt"
  "log"
  "os"
  "time"
  "encoding/binary"
  "image/color"
	"github.com/hajimehoshi/ebiten"
	tm "github.com/buger/goterm"
)

type Cpu struct {
  Memory []byte
  Stack  []uint16
  RomSize   int
  PC     uint16
  SP     uint8
  I      uint16
  Vx     []uint16
  Vy     []uint16
  DrawMatrix [][]int
  Screen *ebiten.Image
}

func NewCpu(path string, screen *ebiten.Image) *Cpu {
  file, err := os.Open(path)

  if err != nil {
    log.Fatal(err)
  }

  rand.Seed(time.Now().UnixNano())

  cpu := Cpu{
    Memory: make([]byte, 4096),
    Stack: make([]uint16, 16),
    Vx: make([]uint16, 16),
    PC: 0x200, //Program Counter
    SP: 0,  //Stack Pointer
    I: 0,
    Screen: screen,
    DrawMatrix: make([][]int, 64),
  }

  fontSet := [80]byte {
    0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
    0x20, 0x60, 0x20, 0x20, 0x70, // 1
    0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
    0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
    0x90, 0x90, 0xF0, 0x10, 0x10, // 4
    0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
    0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
    0xF0, 0x10, 0x20, 0x40, 0x40, // 7
    0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
    0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
    0xF0, 0x90, 0xF0, 0x90, 0x90, // A
    0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
    0xF0, 0x80, 0x80, 0x80, 0xF0, // C
    0xE0, 0x90, 0x90, 0x90, 0xE0, // D
    0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
    0xF0, 0x80, 0xF0, 0x80, 0x80,  // F
  }

  for i, b := range fontSet {
    cpu.Memory[i] = b
  }

  for i := range cpu.DrawMatrix {
    cpu.DrawMatrix[i] = make([]int, 32)
  }

  pointer := 0
  b := make([]byte, 1)

  for {
    bytes_read, _ := file.Read(b)

    if bytes_read == 0 {
      break
    }

    cpu.Memory[0x200 + pointer] = b[0]
    pointer++
  }

  cpu.RomSize = pointer

  log.Printf("Rom Size %d", pointer)

  return &cpu
}

func (cpu *Cpu) DecomposeInstAt(PC uint16) (uint16, uint16, uint16, uint16, uint16, uint16) {
  inst := binary.BigEndian.Uint16(cpu.Memory[PC : PC + 2])
  op_1 := (inst & 0xF000) >> 12
  op_2 := (inst & 0x0F00) >> 8
  op_3 := (inst & 0x00F0) >> 4
  op_4 := (inst & 0x000F)

  kk  := (inst & 0x00FF)
  nnn := (inst & 0x0FFF)

  return op_1, op_2, op_3, op_4, kk, nnn
}

func (cpu *Cpu) DrawMatrixToScreen() {
  for x := range cpu.DrawMatrix {
    for y := range cpu.DrawMatrix[x] {
      if cpu.DrawMatrix[x][y] == 1 {
        c := color.RGBA{ 255, 255, 255, 1 }
        cpu.Screen.Set(x, y, c)
      }
    }
  }
}

func (cpu *Cpu) Tick() {
//  op_1, op_2, op_3, op_4, kk, nnn := cpu.DecomposeInstAt(cpu.PC)

  cpu.RunTick(cpu.DecomposeInstAt(cpu.PC))
  cpu.PC += 2

  cpu.DrawMatrixToScreen()

  return
}

func (cpu *Cpu) RunTick(op_1, op_2, op_3, op_4, kk, nnn uint16) {
  switch op_1 {
  case 0x1:
    cpu.i1nnn(nnn)
  case 0xA:
    cpu.Annn(nnn)
  case 0xC:
    cpu.Cxkk(op_2, kk)
  case 0x3:
    cpu.i3xkk(op_2, kk)
  case 0x6:
    cpu.i6xkk(op_2, kk)
  case 0x7:
    cpu.i7xkk(op_2, kk)
  case 0x8:
    cpu.i8xyn(op_1, op_2, op_3, op_4)
  case 0xD:
    cpu.Dxyn(op_2, op_3, op_4)
  default:
    log.Printf("Unknown instruction: %X", op_1)
  }
}

func (cpu *Cpu) friendlyOpCode(currentPos uint16) string {
  op_1, op_2, op_3, op_4, kk, nnn := cpu.DecomposeInstAt(currentPos)

  switch op_1 {
  case 0x1:
    return fmt.Sprintf("JP  %X",  nnn)
  case 0xA:
    return fmt.Sprintf("LD  I,  %X", nnn)
  case 0xC:
    return fmt.Sprintf("RND V%d, %X",  op_2, kk)
  case 0x3:
    return fmt.Sprintf("SE  V%d, %X",  op_2, kk)
  case 0x6:
    return fmt.Sprintf("LD  V%d, %X",  op_2, kk)
  case 0x7:
    return fmt.Sprintf("ADD V%d, %X",  op_2, kk)
  case 0x8:
    switch op_4 {
    case 0:
      return fmt.Sprintf("LD  V%d, V%d",  op_2, op_3)
    }
  case 0xD:
    return fmt.Sprintf("DRW V%d, V%d, %X",  op_2, op_3, op_4)
  default:
    return fmt.Sprintf("UNK")
  }

  return ""
}

func (cpu *Cpu) DebugTick() {
  tm.Clear()
  tm.MoveCursor(1, 1)

  for i := 0 ; i < cpu.RomSize; i += 2 {
    currentPos := 0x200 + uint16(i)
    op_1, op_2, op_3, op_4, _, _ := cpu.DecomposeInstAt(currentPos)
    inst := cpu.friendlyOpCode(currentPos)

    if currentPos == cpu.PC {
      s := fmt.Sprintf("%d: %X%X%X%X %s\n", currentPos, op_1, op_2, op_3, op_4, inst)
      tm.Printf(tm.Background(s, tm.RED))
      tm.Printf("\n")
    } else {
      tm.Printf("%d: %X%X%X%X %s\n", currentPos, op_1, op_2, op_3, op_4, inst)
    }
  }

  tm.MoveCursor(40, 1)

  for i, v := range cpu.Vx {
    tm.Printf("[V%d: %X] ", i, v)
  }

  tm.MoveCursor(40, 10)
  tm.Printf("I:%d PC:%d SC:%d", cpu.I, cpu.PC, cpu. SP)
  tm.Flush()

  cpu.RunTick(cpu.DecomposeInstAt(cpu.PC))
  cpu.PC += 2

  cpu.DrawMatrixToScreen()
//  cpu.DecomposeInstAt(0x200 + uint16(i))
}

func (cpu *Cpu) Annn(nnn uint16) {
  cpu.I = nnn
}

func (cpu *Cpu) Cxkk(x, kk uint16) {
  cpu.Vx[x] =  uint16(rand.Intn(255)) & kk
  //cpu.Vx[x] =  uint16(rand.Intn(255))
}

func (cpu *Cpu) i3xkk(Vx, kk uint16) {
  if cpu.Vx[Vx] == kk {
    cpu.PC += 2
  }
}

func (cpu *Cpu) i6xkk(Vx, kk uint16) {
  cpu.Vx[Vx] = kk
}

func (cpu *Cpu) i7xkk(Vx, kk uint16) {
  cpu.Vx[Vx] += kk
}

func (cpu *Cpu) i8xyn(op_1, Vx, Vy, op_4 uint16) {
  switch op_4 {
  case 0:
    cpu.i8xy0(Vx, Vy)
  default:
    log.Printf("Unknown i8xy*: %X", op_4)
  }
}

func (cpu *Cpu) i8xy0(Vx, Vy uint16) {
  cpu.Vx[Vx] = cpu.Vx[Vy]
}

func (cpu *Cpu) i1nnn(nnn uint16) {
  cpu.PC = (nnn - 0x02)
}

func (cpu *Cpu) Dxyn(Vx, Vy, n uint16) {
  sprite := cpu.Memory[cpu.I : cpu.I + n]

  for y, value := range sprite {
    for x := 7; x > 0; x-- {
      if (value >> x) & 0x01 == 0 {
        continue
      }
      cpu.DrawMatrix[int(cpu.Vx[Vx]) + (7 - x)][int(cpu.Vx[Vy]) + y] ^= 0x01
    }
  }
}
