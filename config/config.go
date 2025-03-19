package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

// Json config structure
type Config struct {
	Mode           string `json:"mode"` //Options "PipeMode", "TCPMode", "SimMode"
	HttpPort       string `json:"httpPort"`
	TCPAddress     string `json:"tcpAddress"`
	SimModePipe    string `json:"simModePipe"`    //Path to SimMode Pipe
	MsgFromSimPipe string `json:"msgFromSimPipe"` //Path to IN Pipe
	MsgToSimPipe   string `json:"msgToSimPipe"`   //Path to OUT Pipe
}

var Cfg *Config

// Load config from json file
func LoadConfig(filePath string) error {
	file, err := os.Open(filePath)
	//Check if file exists
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("Warning: File not found (%s). Loading default values.", filePath)
		// Set default values
		Cfg = &Config{
			Mode:           "PipeMode",
			TCPAddress:     "localhost:8081",
			SimModePipe:    "/tmp/simulator",
			MsgFromSimPipe: "/tmp/msgFromSim",
			MsgToSimPipe:   "/tmp/msgToSim",
			HttpPort:       "8080",
		}
		return nil

	} else if err != nil {
		//Check for other errors
		return fmt.Errorf("Error opening file: %w", err)
	}
	defer file.Close()
	//Decode json --> config struct
	decoder := json.NewDecoder(file)
	var config Config
	if err := decoder.Decode(&config); err != nil {
		return fmt.Errorf("Error decoding .json-File: %w", err)
	}
	Cfg = &config
	log.Println("Config loaded.")
	return nil
}
