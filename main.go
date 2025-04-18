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
	config.LoadConfig("ExtruderUIConfig.json")

	switch config.Cfg.Mode {
	case "SimMode":
		log.Println("Starting in mode: Sim-Pipe")
		go pipes.FromSimPipeHandler()
	case "PipeMode":
		log.Println("Starting in mode: Msg-Pipe")
		go pipes.FromSimPipeHandler()
		go pipes.FromUserPipeHandler()
	case "TCPMode":
		manager := tcp.GetConnectionManager(config.Cfg.TCPAddress)
		defer manager.CloseConnection()
		log.Println("Starting in mode: TCP/IP-Socket")
		go tcp.TCPDataHandler(config.Cfg.TCPAddress)
	default:
		log.Println("Mode unknown:", config.Cfg.Mode)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
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
