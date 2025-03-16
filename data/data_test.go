package data

import (
	"bytes"
	"encoding/hex"
	"log"
	"testing"
)

// Mock message channel
var mockMessageChannel = make(chan string, 1)

// Override the messageChannel for testing
func init() {
	messageChannel = mockMessageChannel
}

// Test for GetValuesFromRow
func TestGetValuesFromRow(t *testing.T) {
	// Setup mock data
	colMap = map[int]*Datapoint{
		1: &Datapoint{},
		2: &Datapoint{},
	}
	//timestampLayout := "2006-01-02"

	line := "2024-12-08 12:00:00,25.5,13.7"
	separator := ","

	// Execute the function
	GetValuesFromRow(line, separator)

	// Validate the results
	if colMap[1].Value != 25.5 {
		t.Errorf("Expected value 25.5, got %v", colMap[1].Value)
	}

	if colMap[2].Value != 13.7 {
		t.Errorf("Expected value 13.7, got %v", colMap[2].Value)
	}

	if colMap[1].Timestamp.IsZero() {
		t.Errorf("Expected timestamp to be parsed, but got zero value")
	}
}

// Test for GetValuesFromRow with invalid input
func TestGetValuesFromRow_InvalidInput(t *testing.T) {
	// Setup mock data
	colMap = map[int]*Datapoint{
		1: &Datapoint{},
		2: &Datapoint{},
	}
	timestampLayout := "2006-01-02"

	line := "Invalid data"
	separator := ","

	// Execute the function
	GetValuesFromRow(line, separator)

	// Validate that no data points were updated
	if colMap[1].Value != 0 {
		t.Errorf("Expected value 0, got %v", colMap[1].Value)
	}

	if colMap[2].Value != 0 {
		t.Errorf("Expected value 0, got %v", colMap[2].Value)
	}

	if !colMap[1].Timestamp.IsZero() {
		t.Errorf("Expected timestamp to remain zero, but got %v", colMap[1].Timestamp)
	}
}

// Test for GetValueFromMsg with valid input
func TestGetValueFromMsg_ValidInput(t *testing.T) {
	// Setup mock data
	msgInMap = map[byte]*Datapoint{
		1: &Datapoint{},
	}

	// Correct Hexadecimal string with 8 bytes
	line := "01020304AABBCCDD" // 8 bytes
	decodedMsg, _ := hex.DecodeString(line)
	expectedValue := uint32(decodedMsg[4])<<24 | uint32(decodedMsg[5])<<16 | uint32(decodedMsg[6])<<8 | uint32(decodedMsg[7])

	// Execute the function
	GetValueFromMsg(line)

	// Validate the results
	if datapoint, exists := msgInMap[1]; exists {
		if datapoint.Value != float32(expectedValue) {
			t.Errorf("Expected value %v, got %v", float32(expectedValue), datapoint.Value)
		}

		if datapoint.Timestamp.IsZero() {
			t.Errorf("Expected timestamp to be set, but got zero value")
		}
	} else {
		t.Errorf("Datapoint with ID 1 not found in msgInMap")
	}
}

// Test for GetValueFromMsg with short input
func TestGetValueFromMsg_ShortInput(t *testing.T) {
	// Capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(nil) // Reset log output after the test

	// Short Hexadecimal string
	line := "01020304" // Only 4 bytes

	// Execute the function
	GetValueFromMsg(line)

	// Validate the log output and ensure no panic occurred
	expectedMessage := "Expected 8 bytes, got 4 bytes"
	if !bytes.Contains(logBuffer.Bytes(), []byte(expectedMessage)) {
		t.Errorf("Expected log message '%s', got '%s'", expectedMessage, logBuffer.String())
	}
}
