package controls

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleControlRequest(t *testing.T) {
	testCases := []struct {
		name           string
		id             byte
		value          uint32
		expectedStatus int
	}{
		{
			name:           "Valid Screw RPM Request",
			id:             screw_rpm_id,
			value:          1000,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid Spooler RPM Request",
			id:             spooler_rpm_id,
			value:          500,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid Heater PWM Request",
			id:             heater_pwm_id,
			value:          75,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Prepare request body
			controlData := ControlData{Value: tc.value}
			jsonBody, _ := json.Marshal(controlData)

			// Create request
			req, err := http.NewRequest("POST", "/control", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			// Create ResponseRecorder
			w := httptest.NewRecorder()

			// Call handler based on ID
			switch tc.id {
			case screw_rpm_id:
				ScrewRpmHandler(w, req)
			case spooler_rpm_id:
				SpoolerRpmHandler(w, req)
			case heater_pwm_id:
				HeaterPwmHandler(w, req)
			default:
				t.Fatalf("Unexpected ID: %x", tc.id)
			}

			// Check response status
			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			// Check response body
			response := w.Body.String()
			if !strings.Contains(response, `{"status": "success"}`) {
				t.Errorf("Unexpected response body: %s", response)
			}
		})
	}
}

func TestButtonHandlers(t *testing.T) {
	testCases := []struct {
		name           string
		handler        func(http.ResponseWriter, *http.Request)
		expectedStatus int
	}{
		{
			name:           "Start Button",
			handler:        ButtonStartHandler,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Emergency Stop Button",
			handler:        ButtonEmergencyStopHandler,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("POST", "/control", nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			// Create ResponseRecorder
			w := httptest.NewRecorder()

			// Call handler
			tc.handler(w, req)

			// Check response status
			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}

	// Test invalid method
	t.Run("Invalid Method", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/control", nil)
		w := httptest.NewRecorder()
		ButtonStartHandler(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d for invalid method, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}
