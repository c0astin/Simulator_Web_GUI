package main

import (
	"extruder_web_gui/config"
	"extruder_web_gui/controls"
	"extruder_web_gui/data"
	"extruder_web_gui/pipes"
	"extruder_web_gui/tcp"
	"log"
	"net/http"
)

func main() {
	config.LoadConfig("ExtruderConfig.json")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	switch config.Cfg.Mode {
	case "SimMode":
		log.Println("Starte im Modus: SimLog-Pipe")
		go pipes.FromSimPipeHandler()
	case "PipeMode":
		log.Println("Starte im Modus: Msg-Pipe")
		go pipes.FromSimPipeHandler()
		go pipes.FromUserPipeHandler()
	case "TCPMode":
		manager := tcp.GetConnectionManager(config.Cfg.TCPAddress)
		defer manager.CloseConnection()
		log.Println("Starte im Modus: TCP/IP-Socket")
		go tcp.TCPDataHandler(config.Cfg.TCPAddress)
	default:
		log.Println("Modus nicht bekannt:", config.Cfg.Mode)
	}

	http.HandleFunc("/", data.MainViewHandler)
	http.HandleFunc("/data", data.DataHandler)
	http.HandleFunc("/control/start", controls.ButtonStartHandler)
	http.HandleFunc("/control/stop", controls.ButtonEmergencyStopHandler)
	http.HandleFunc("/control/screw-rpm", controls.ScrewRpmHandler)
	http.HandleFunc("/control/spooler-rpm", controls.SpoolerRpmHandler)
	http.HandleFunc("/control/heater-pwm", controls.HeaterPwmHandler)
	http.HandleFunc("/control/mode", controls.ModeSwitchHandler)

	http.HandleFunc("/messages", data.MessagesHandler)

	log.Fatal(http.ListenAndServe(":"+config.Cfg.HttpPort, nil))
}
