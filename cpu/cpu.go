package cpu

import (
  "math/rand"
  "log"
  "os"
  "encoding/binary"
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
}

func LoadRom(path string) *Cpu {
  file, err := os.Open(path)

  if err != nil {
    log.Fatal(err)
  }

  cpu := Cpu{
    Memory: make([]byte, 4096),
    Stack: make([]uint16, 16),
    Vx: make([]uint16, 16),
    Vy: make([]uint16, 16),
    PC: 0x200, //Program Counter
    SP: 0,  //Stack Pointer
    I: 0,
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

func (cpu *Cpu) Tick() {
  inst := binary.BigEndian.Uint16(cpu.Memory[cpu.PC : cpu.PC + 2])
  op_1 := (inst & 0xF000) >> 12
  op_2 := (inst & 0x0F00) >> 8
  op_3 := (inst & 0x00F0) >> 4
  op_4 := (inst & 0x000F)

  nnn  := (inst & 0x0FFF)

  switch op_1 {
  case 0xA:
    cpu.Annn(nnn)
  case 0xC:
    cpu.Cxkk(inst, op_1, op_2, op_3, op_4)
  }

  log.Printf("%X %X", inst, op_1)
  cpu.PC += 2
}

func (cpu *Cpu) Debug() {
  for _, b := range cpu.Memory {
    log.Printf("%X", b)
  }
}

func (cpu *Cpu) Annn(nnn uint16) {
  cpu.I = nnn
}
