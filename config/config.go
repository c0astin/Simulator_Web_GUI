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
		log.Printf("Warnung: Datei nicht gefunden (%s). Default-Werte werden geladen.", filePath)
		// Set default values
		Cfg = &Config{
			Mode:           "PipeMode",
			TCPAddress:     "localhost:8081",
			SimModePipe:    "/tmp/simulator",
			MsgFromSimPipe: "/tmp/msgFromSim",
			MsgToSimPipe:   "/tmp/msgToSim",
		}
		return nil

	} else if err != nil {
		//Check for other errors
		return fmt.Errorf("Fehler beim Öffnen der Datei: %w", err)
	}
	defer file.Close()
	//Decode json --> config struct
	decoder := json.NewDecoder(file)
	var config Config
	if err := decoder.Decode(&config); err != nil {
		return fmt.Errorf("Fehler beim Decodieren der json-Datei: %w", err)
	}
	Cfg = &config
	log.Println("Konfiguration erfolgreich geladen.")
	return nil
}
