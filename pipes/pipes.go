package pipes

import (
	"bufio"
	"encoding/binary"
	"extruder_web_gui/config"
	"extruder_web_gui/data"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// Handler for (User Program -> Web UI) Pipe
func FromUserPipeHandler() {
	var delayReconnect time.Duration = 2 * time.Second
	for {
		file, err := os.OpenFile(config.Cfg.MsgFromSimPipe, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			log.Printf("Failed to open pipe (User->UI) file (%s): %v", config.Cfg.MsgFromSimPipe, err)
			time.Sleep(delayReconnect)
			continue
		}
		log.Println("Pipe (User->UI) connected successfully")

		reader := bufio.NewReader(file)
		buffer := make([]byte, 8)
		for {
			_, err := io.ReadFull(reader, buffer)
			if err != nil {
				if err == io.EOF {
					log.Println("Pipe (User->UI) closed by writer. Trying to reestablish connection...")
					break
				}
				log.Fatal("Error reading from Pipe (User -> UI):", err)
			}
			buffer = reverseBytes(buffer)
			data.GetValueFromMsg(buffer)
		}
		//only reached when writer closes connection
		log.Println("Pipe (User->UI) disconnected. Reconnecting...")
		file.Close()
	}
}

// Handler for (Simulator --> Web-UI) Pipe
func FromSimPipeHandler() {
	var delayReconnect time.Duration = 2 * time.Second
	for {
		//Open pipe with Read Only permissions
		file, err := os.OpenFile(config.Cfg.SimModePipe, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			log.Printf("Failed to open pipe (Sim->UI) file (%s): %v", config.Cfg.SimModePipe, err)
			time.Sleep(delayReconnect)
			continue
		}
		log.Println("Pipe (Sim->UI) connected successfully")

		reader := bufio.NewReader(file)
		for {
			// Read pipe by line
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					log.Println("Pipe (Sim->UI) closed by writer. Trying to reestablish connection...")
					break
				}
			}

			data.GetStatsFromRow(line, "|")
			if config.Cfg.Mode == "SimMode" {
				data.GetValuesFromRow(line, "|")
			}
		}
		//only reached when writer closes connection
		log.Println("Pipe (Sim->UI) disconnected, reconnecting...")
		file.Close()
	}
}

// Function to write Msg to Pipe
func ToUserPipe(pipePath string, id byte, val uint32) {
	//Open pipe with Write Only permissions
	pipe, err := os.OpenFile(pipePath, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatalf("Failed to open file.")
	}
	defer pipe.Close()

	//Compose 8 byte message from id and value
	messageBytes := make([]byte, 8)
	messageBytes[7] = id
	binary.LittleEndian.PutUint32(messageBytes[0:], val)

	//Write 8 byte message to pipe
	_, err = pipe.Write(messageBytes)
	if err != nil {
		log.Fatalf("Failed to write to pipe")
	}

	fmt.Printf("Message sent to pipe: %x\n", reverseBytes(messageBytes))
}

// Helper function to reverse byte order
func reverseBytes(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
