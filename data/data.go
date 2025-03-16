package data

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// Data Structure for incoming metrics
type Datapoint struct {
	Timestamp time.Time
	Value     float32
}
type Dataset struct {
	Temperature   Datapoint
	Diameter      Datapoint
	SpoolerRpm    Datapoint
	ScrewRpm      Datapoint
	HeaterPwm     Datapoint
	ContactSwitch Datapoint
}

type SpoolStats struct {
	WindingDiameter Datapoint
	AvgFilDiameter  Datapoint
	NbrOfWindings   Datapoint
	FilamentMass    Datapoint
}

var curSpoolStats SpoolStats
var prevSpoolStats SpoolStats
var currentData Dataset

// SimMode: Assign columns to parameter
var colMap = map[int]*Datapoint{
	5: &currentData.Diameter,
	6: &currentData.Temperature,
	3: &currentData.SpoolerRpm,
	2: &currentData.ScrewRpm,
	4: &currentData.HeaterPwm,
	7: &currentData.ContactSwitch,
}

var statsColMap = map[int]*Datapoint{
	8:  &curSpoolStats.WindingDiameter,
	9:  &curSpoolStats.AvgFilDiameter,
	10: &curSpoolStats.NbrOfWindings,
	11: &curSpoolStats.FilamentMass,
}

// MsgMode: Assign IDs of incoming messages to parameter
var msgInMap = map[byte]*Datapoint{
	0x01: &currentData.Diameter,
	0x02: &currentData.Temperature,
	0x03: &currentData.SpoolerRpm,
	0x04: &currentData.ScrewRpm,
	0x05: &currentData.HeaterPwm,
	0x06: &currentData.ContactSwitch,
}

var messageChannel = make(chan string)

var timestampLayout string = "15:04:05.000"

func DataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
		"diameter": {"timestamp": "%s", "value": %.2f},
		"temperature": {"timestamp": "%s", "value": %.2f},
		"spoolerRpm": {"timestamp": "%s", "value": %.2f},
		"screwRpm": {"timestamp": "%s", "value": %.2f},
		"heaterPwm": {"timestamp": "%s", "value": %.2f},
		"contactSwitch": {"timestamp": "%s", "value": %.2f},
		"windingDiameter": {"timestamp": "%s", "value": %.2f},
		"avgFilDiameter": {"timestamp": "%s", "value": %.2f},
		"nbrOfWindings": {"timestamp": "%s", "value": %.2f},
		"filamentMass": {"timestamp": "%s", "value": %.2f},
		"prevWindingDiameter": {"timestamp": "%s", "value": %.2f},
		"prevAvgFilDiameter": {"timestamp": "%s", "value": %.2f},
		"prevNbrOfWindings": {"timestamp": "%s", "value": %.2f},
		"prevFilamentMass": {"timestamp": "%s", "value": %.2f}
	}`,
		currentData.Diameter.Timestamp.Format(time.RFC3339), currentData.Diameter.Value,
		currentData.Temperature.Timestamp.Format(time.RFC3339), currentData.Temperature.Value,
		currentData.SpoolerRpm.Timestamp.Format(time.RFC3339), currentData.SpoolerRpm.Value,
		currentData.ScrewRpm.Timestamp.Format(time.RFC3339), currentData.ScrewRpm.Value,
		currentData.HeaterPwm.Timestamp.Format(time.RFC3339), currentData.HeaterPwm.Value,
		currentData.ContactSwitch.Timestamp.Format(time.RFC3339), currentData.ContactSwitch.Value,
		curSpoolStats.WindingDiameter.Timestamp.Format(time.RFC3339), curSpoolStats.WindingDiameter.Value,
		curSpoolStats.AvgFilDiameter.Timestamp.Format(time.RFC3339), curSpoolStats.AvgFilDiameter.Value,
		curSpoolStats.NbrOfWindings.Timestamp.Format(time.RFC3339), curSpoolStats.NbrOfWindings.Value,
		curSpoolStats.FilamentMass.Timestamp.Format(time.RFC3339), curSpoolStats.FilamentMass.Value,
		prevSpoolStats.WindingDiameter.Timestamp.Format(time.RFC3339), prevSpoolStats.WindingDiameter.Value,
		prevSpoolStats.AvgFilDiameter.Timestamp.Format(time.RFC3339), prevSpoolStats.AvgFilDiameter.Value,
		prevSpoolStats.NbrOfWindings.Timestamp.Format(time.RFC3339), prevSpoolStats.NbrOfWindings.Value,
		prevSpoolStats.FilamentMass.Timestamp.Format(time.RFC3339), prevSpoolStats.FilamentMass.Value,
	)
}

func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Flush writer to send initial headers
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	for message := range messageChannel {
		fmt.Fprintf(w, "data: %s\n\n", message)
		flusher.Flush()
	}
}

// Handler - Update main view schematics
func MainViewHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, currentData)
}

// SimMode: Function to get Values from Simulator Pipe row
func GetValuesFromRow(line string, separator string) {
	strArr := strings.Split(line, separator)
	if len(strArr) > 2 {
		for idx, str := range strArr {
			strArr[idx] = strings.TrimSpace(str)
		}
		timeStr := strings.Split(strArr[0], " ")[0]
		parsedTime, err := time.Parse(timestampLayout, timeStr)
		if err != nil {
			fmt.Println("Error parsing timestamp:", err)
			return
		}
		for key, val := range colMap {
			if key < len(strArr) { // Ensure key is within bounds of strArr
				val64, _ := strconv.ParseFloat(strArr[key], 32)
				val.Value = float32(val64)
				val.Timestamp = parsedTime
			}
		}
	}
	messageChannel <- line
}

// Function to get Winding Stats from Simulator Pipe row
func GetStatsFromRow(line string, separator string) {
	strArr := strings.Split(line, separator)
	if len(strArr) > 2 {
		for idx, str := range strArr {
			strArr[idx] = strings.TrimSpace(str)
		}
		timeStr := strings.Split(strArr[0], " ")[0]
		parsedTime, err := time.Parse(timestampLayout, timeStr)
		if err != nil {
			fmt.Println("Error parsing timestamp:", err)
			return
		}
		var tempStats SpoolStats = curSpoolStats
		for key, val := range statsColMap {
			if key < len(strArr) { // Ensure key is within bounds of strArr
				val64, _ := strconv.ParseFloat(strArr[key], 32)
				val.Value = float32(val64)
				val.Timestamp = parsedTime
			}
		}
		if tempStats.FilamentMass.Value > curSpoolStats.FilamentMass.Value && tempStats.WindingDiameter.Value > 12 {
			prevSpoolStats = tempStats
		}

	}
	messageChannel <- line
}

func GetValueFromMsg(msg []byte) {
	//Check if msg length is valid
	if len(msg) != 8 {
		log.Printf("Expected 8 bytes, got %d bytes", len(msg))
		return // Skip this message
	}

	//fmt.Printf("Received message: %x\n", msg)

	//Decode id from first byte
	id := msg[0]
	//Decode value from the last 4 bytes
	value := (uint32(msg[4]) << 24) | (uint32(msg[5]) << 16) | (uint32(msg[6]) << 8) | uint32(msg[7])

	//Check if id matches metric and assign value
	if datapoint, ok := msgInMap[id]; ok {
		datapoint.Value = float32(value)
		datapoint.Timestamp = time.Now() // Set the current time as the timestamp
	} else {
		//fmt.Printf("No matching ID %d for incoming message\n", id)
		return // Skip this message
	}

	message := fmt.Sprintf("%s: %x", time.Now().Format("15:04:05"), msg)
	messageChannel <- message
}
