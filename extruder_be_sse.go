package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Base int

const (
	SimMode Base = iota
	CanMode
	UserLogMode
)

var currentMode Base = SimMode

type Datapoint struct {
	Timestamp time.Time
	Value     float32
}

type Dataset struct {
	Temperature Datapoint
	Diameter    Datapoint
	SpoolerRpm  Datapoint
	ScrewRpm    Datapoint
	HeaterPwm   Datapoint
}

var currentData Dataset

const timestampCol int = 0
const timestampLayout string = "15:04:05.000"

var colMap = map[int]*Datapoint{
	7: &currentData.Diameter,
	8: &currentData.Temperature,
	4: &currentData.SpoolerRpm,
	2: &currentData.ScrewRpm,
	6: &currentData.HeaterPwm,
}

var msgMap = map[byte]*Datapoint{
	0x01: &currentData.Diameter,
	0x02: &currentData.Temperature,
	0x03: &currentData.SpoolerRpm,
	0x04: &currentData.ScrewRpm,
	0x05: &currentData.HeaterPwm,
}

func main() {
	//c := make(chan string, 10)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	go dataInPipeHandler("/tmp/simulator")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, currentData)
	})

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"diameter": {"timestamp": "%s", "value": %.2f},
			"temperature": {"timestamp": "%s", "value": %.2f},
			"spoolerRpm": {"timestamp": "%s", "value": %.2f},
			"screwRpm": {"timestamp": "%s", "value": %.2f},
			"heaterPwm": {"timestamp": "%s", "value": %.2f}
		}`,
			currentData.Diameter.Timestamp.Format(time.RFC3339), currentData.Diameter.Value,
			currentData.Temperature.Timestamp.Format(time.RFC3339), currentData.Temperature.Value,
			currentData.SpoolerRpm.Timestamp.Format(time.RFC3339), currentData.SpoolerRpm.Value,
			currentData.ScrewRpm.Timestamp.Format(time.RFC3339), currentData.ScrewRpm.Value,
			currentData.HeaterPwm.Timestamp.Format(time.RFC3339), currentData.HeaterPwm.Value,
		)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func dataInPipeHandler(path string) {
	tty, err := os.Open("/tmp/simulator")
	if err != nil {
		log.Fatal("Make named pipe file error:", err)
	}
	defer tty.Close()
	reader := bufio.NewReader(tty)

	for {
		line, err := reader.ReadBytes('\n')
		if len(line) == 0 {
			continue
		}
		if err == nil {
			switch currentMode {
			case SimMode:
				getValuesFromRow(string(line), "|")
			}
			fmt.Println(currentData.Diameter.Value)
			fmt.Println(currentData.Diameter.Timestamp)
			//c <- string(line)

		} else {
			_ = line
			fmt.Printf("READ err: %v\n", err)
		}
	}
	tty.Close()
}

func getValuesFromRow(line string, separator string) {

	strArr := strings.Split(line, separator)
	fmt.Println(line)
	if len(strArr) > 2 {
		for idx, str := range strArr {
			strArr[idx] = strings.TrimSpace(str)
		}
		timeStr := strings.Split(strArr[0], " ")
		fmt.Println(timeStr[0])
		for key, val := range colMap {
			val64, _ := strconv.ParseFloat(strArr[key], 32)
			val.Value = float32(val64)
			val.Timestamp, _ = time.Parse(timestampLayout, timeStr[0])

		}

	}

	fmt.Println("-----")
}
