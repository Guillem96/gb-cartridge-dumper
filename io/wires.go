package io

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type yamlGameBoyMappingContainer struct {
	Gameboy GameBoyRaspberryMapping `yaml:"gameboy-pins"`
}

// GameBoyRaspberryMapping structure to map RaspberryPi GPIO pins to GameBoy cartridge pins
type GameBoyRaspberryMapping struct {
	RD int32 `yaml:"RD"`
	WR int32 `yaml:WR`

	A0  int32 `yaml:"A0"`
	A1  int32 `yaml:"A1"`
	A2  int32 `yaml:"A2"`
	A3  int32 `yaml:"A3"`
	A4  int32 `yaml:"A4"`
	A5  int32 `yaml:"A5"`
	A6  int32 `yaml:"A6"`
	A7  int32 `yaml:"A7"`
	A8  int32 `yaml:"A8"`
	A9  int32 `yaml:"A9"`
	A10 int32 `yaml:"A10"`
	A11 int32 `yaml:"A11"`
	A12 int32 `yaml:"A12"`
	A13 int32 `yaml:"A13"`
	A14 int32 `yaml:"A14"`
	A15 int32 `yaml:"A15"`

	D0 int32 `yaml:"D0"`
	D1 int32 `yaml:"D1"`
	D2 int32 `yaml:"D2"`
	D3 int32 `yaml:"D3"`
	D4 int32 `yaml:"D4"`
	D5 int32 `yaml:"D5"`
	D6 int32 `yaml:"D6"`
	D7 int32 `yaml:"D7"`
}

// ParseWireMapping parses the yaml file containing the mapping from GameBoy cartridge pins to Raspberry pins
func ParseWireMapping(path string) *GameBoyRaspberryMapping {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("[ERROR] %v cannot be read. %v\n", path, err)
		os.Exit(1)
	}

	connectionMapping := &yamlGameBoyMappingContainer{}
	if err = yaml.Unmarshal(yamlFile, connectionMapping); err != nil {
		fmt.Printf("[ERROR] Parsing file content")
		os.Exit(1)
	}

	return &connectionMapping.Gameboy
}
