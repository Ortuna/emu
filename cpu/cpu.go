package cpu

type Cpu struct {
  rom [512]uint8
}

func LoadRom(path string) *Cpu {
  file, err := os.Open("./rom.ch8")

	if err != nil {
	  log.Fatal(err)
  }

  cpu := Cpu{}


  return cpu
}


