package config

import (
	"os"
	"testing"
)

// Test for LoadConfig: Valid config file
func TestLoadConfig_ValidFile(t *testing.T) {
	// Create valid config file
	testConfig := `{
		"mode": "TCPMode",
		"tcpAddress": "192.168.0.10:8081",
		"simModePipe": "/tmp/test_simulator",
		"msgFromSimPipe": "/tmp/test_msgFromSim",
		"msgToSimPipe": "/tmp/test_msgToSim"
	}`
	testFilePath := "test_config.json"
	err := os.WriteFile(testFilePath, []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Error: Can not create config file: %v", err)
	}
	defer os.Remove(testFilePath) // Remove file after test

	// Call LoadConfig with valid file
	// Expected result: Cfg values set to defined values
	err = LoadConfig(testFilePath)
	if err != nil {
		t.Errorf("Error: Can not load test config file: %v", err)
	}

	expectedMode := "TCPMode"
	expectedMsgFromSimPipe := "/tmp/test_msgFromSim"
	// Check if values are correctly loaded
	if Cfg.Mode != expectedMode {
		t.Errorf("Mode wasn't set to expected value. Expected: %v, received: %v", expectedMode, Cfg.Mode)
	}
	if Cfg.MsgFromSimPipe != expectedMsgFromSimPipe {
		t.Errorf("MsgFromSimPipe wasn't set to expected value. Expected: %v, received: %v", expectedMsgFromSimPipe, Cfg.MsgFromSimPipe)
	}
}

// Test for LoadConfig: File does not exist
func TestLoadConfig_MissingFile(t *testing.T) {
	// Call LoadConfig with non existing file name
	// Expected result: Load default values
	LoadConfig("missing_config.json")

	expectedMode := "PipeMode"
	expectedMsgToSimPipe := "/tmp/msgToSim"
	if Cfg.Mode != expectedMode {
		t.Errorf("Mode wasn't set to default value. Expected: %v, received: %v", expectedMode, Cfg.Mode)
	}
	if Cfg.MsgToSimPipe != expectedMsgToSimPipe {
		t.Errorf("MsgToSimPipe wasn't set to default value. Expected: %v, received: %v", expectedMsgToSimPipe, Cfg.MsgToSimPipe)
	}
}

// Test for LoadConfig: Invalid Json
func TestLoadConfig_InvalidJSON(t *testing.T) {
	// Create invalid json for config file
	invalidTestFilePath := "invalid_config.json"
	os.WriteFile(invalidTestFilePath, []byte("INVALID"), 0644)
	defer os.Remove(invalidTestFilePath) // Remove file after test

	// Call LoadConfig with invalid Json File
	// Expected result: Error because of invalid json
	err := LoadConfig(invalidTestFilePath)
	if err == nil {
		t.Errorf("Expected error, but no error thrown.")
	}
}
