package controls

import (
	"encoding/json"
	"extruder_web_gui/config"
	"extruder_web_gui/pipes"
	"extruder_web_gui/tcp"
	"net/http"
)

// Assign IDs to outgoing msgs
const man_auto_switch_id byte = 0x01
const auto_start_id byte = 0x02
const spooler_rpm_id byte = 0x03
const screw_rpm_id byte = 0x04
const heater_pwm_id byte = 0x05
const emergency_stop_id byte = 0x06

type ControlData struct {
	Value uint32 `json:"value"`
}

// Process json input and call function to send data via Pipe or TCP/IP Socket
func handleControlRequest(w http.ResponseWriter, r *http.Request, id byte) {
	var data ControlData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if config.Cfg.Mode == "TCPMode" {
		tcp.SendTCPData(id, data.Value)
	} else {
		pipes.ToUserPipe(config.Cfg.MsgToSimPipe, id, data.Value)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))
}

// Handler for Screw RPM Input
func ScrewRpmHandler(w http.ResponseWriter, r *http.Request) {
	handleControlRequest(w, r, screw_rpm_id)
}

// Handler for Spooler RPM Input
func SpoolerRpmHandler(w http.ResponseWriter, r *http.Request) {
	handleControlRequest(w, r, spooler_rpm_id)
}

// Handler for Heater PWM Input
func HeaterPwmHandler(w http.ResponseWriter, r *http.Request) {
	handleControlRequest(w, r, heater_pwm_id)
}

// Handler for Automatic Mode: Start Button
func ButtonStartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if config.Cfg.Mode == "TCPMode" {
			tcp.SendTCPData(auto_start_id, 1)
		} else {
			pipes.ToUserPipe(config.Cfg.MsgToSimPipe, auto_start_id, 1)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Aktion erfolgreich ausgeführt"))
	} else {
		http.Error(w, "Nur POST-Anfragen erlaubt", http.StatusMethodNotAllowed)
	}
}

// Handler for Emergency Stop Button
func ButtonEmergencyStopHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if config.Cfg.Mode == "TCPMode" {
			tcp.SendTCPData(emergency_stop_id, 1)
		} else {
			pipes.ToUserPipe(config.Cfg.MsgToSimPipe, emergency_stop_id, 1)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Aktion erfolgreich ausgeführt"))
	} else {
		http.Error(w, "Nur POST-Anfragen erlaubt", http.StatusMethodNotAllowed)
	}
}

// Handler for Mode Switch
func ModeSwitchHandler(w http.ResponseWriter, r *http.Request) {
	handleControlRequest(w, r, man_auto_switch_id)
}
