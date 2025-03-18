package data

import (
	"testing"
)

// Test for GetValuesFromRow: Valid input line
func TestGetValuesFromRow_ValidLine(t *testing.T) {
	// Valid example data line with random values
	validLine := "12:00:00.000 | 1.0 | 2.5 | 50 | 100 | 1.7 | 201.1 | 0 | 50 | 2 | 50 | 100"

	// Call GetValuesFromRow with valid data
	// Expected result: stats in data struct updated to expected values
	GetValuesFromRow(validLine, "|")

	// Expected values
	expectedTemperature := float32(201.1)
	expectedDiameter := float32(1.7)

	// Check if values have changed to target values
	if currentData.Temperature.Value != expectedTemperature {
		t.Errorf("Temperature wasn't updated correctly. Expected: %v, received: %v", expectedTemperature, currentData.Temperature.Value)
	}

	if currentData.Diameter.Value != expectedDiameter {
		t.Errorf("Diameter wasn't updated correctly. Expected: %v, received: %v", expectedDiameter, currentData.Diameter.Value)
	}
}

// Test for GetValuesFromRow: Line with invalid values
func TestGetValuesFromRow_InvalidValues(t *testing.T) {
	// Invalid data line
	invalidLine := "12:00:00.000 | 1.0 | 2.5 | 50 | 100 | 2.0 | X | 0 | 50 | 2 | 50 | 100"

	// Call GetValuesFromRow with right formatting but invalid data
	// Expected result: invalid value in data struct not changed
	GetValuesFromRow(invalidLine, "|")

	// Expected values
	expectedTemperature := float32(0) //Initial value
	expectedDiameter := float32(2)

	// Check if values have changed to target values
	if currentData.Temperature.Value != expectedTemperature {
		t.Errorf("Temperature wasn't updated correctly. Expected: %v, received: %v", expectedTemperature, currentData.Temperature.Value)
	}

	if currentData.Diameter.Value != expectedDiameter {
		t.Errorf("Diameter wasn't updated correctly. Expected: %v, received: %v", expectedDiameter, currentData.Diameter.Value)
	}
}

// Test for GetValuesFromRow: Line with invalid formatting
func TestGetValuesFromRow_InvalidFormat(t *testing.T) {
	// Invalid data line
	invalidLine := "12:00:00.000 | 1.0 | 100 "

	currentData.ScrewRpm.Value = float32(30)
	currentData.Diameter.Value = float32(0)

	// Call GetValuesFromRow with right formatting but invalid data
	// Expected result: No change in data struct if input string has less then 3 "|" separators
	GetValuesFromRow(invalidLine, "|")

	// Expected values
	expectedScrewRpm := float32(30) //Initial value
	expectedDiameter := float32(0)

	// Check if values have changed to target values
	if currentData.ScrewRpm.Value != expectedScrewRpm {
		t.Errorf("ScrewRpm wasn't updated correctly. Expected: %v, received: %v", expectedScrewRpm, currentData.ScrewRpm.Value)
	}

	if currentData.Diameter.Value != expectedDiameter {
		t.Errorf("Diameter wasn't updated correctly. Expected: %v, received: %v", expectedDiameter, currentData.Diameter.Value)
	}
}

// Test for GetStatsFromRow: Valid input line
func TestGetStatsFromRow_ValidLine(t *testing.T) {
	// Valid example data line with random values
	validLine := "12:00:00.000 | 1.0 | 2.5 | 50 | 100 | 1 | 201.1 | 0 | 60.5 | 2.20 | 100 | 200"

	// Call GetStatsFromRow with valid data
	// Expected result: stats in data struct updated to expecte values
	GetStatsFromRow(validLine, "|")

	// Expected values
	expectedWindingDiameter := float32(60.5)
	expectedFilamentMass := float32(200)

	// Check if values have changed to target values
	if curSpoolStats.WindingDiameter.Value != expectedWindingDiameter {
		t.Errorf("WindingDiameter wasn't updated correctly. Expected: %v, received: %v", expectedWindingDiameter, curSpoolStats.WindingDiameter.Value)
	}
	if curSpoolStats.FilamentMass.Value != expectedFilamentMass {
		t.Errorf("FilamentMass wasn't updated correctly. Expected: %v, received: %v", expectedFilamentMass, curSpoolStats.FilamentMass.Value)
	}
}

// Test for GetStatsFromRow: Line with invalid values
func TestGetStatsFromRow_InvalidLine(t *testing.T) {
	// Invalid data line
	invalidLine := "12:00:00.000 | X | 1.2 | X | 100 | 1 | 220.5 | 0 | 70| 32 | X | X"

	// Call GetStatsFromRow with invalid data
	// Expected result: Data unchanged
	GetStatsFromRow(invalidLine, "|")

	// Expected values
	expectedWindingDiameter := float32(70)
	expectedFilamentMass := float32(0) //Initial value

	// Check if values have changed to target values with 2 examples
	if curSpoolStats.WindingDiameter.Value != expectedWindingDiameter {
		t.Errorf("WindingDiameter wasn't updated correctly. Expected: %v, received: %v", expectedWindingDiameter, curSpoolStats.WindingDiameter.Value)
	}
	if curSpoolStats.FilamentMass.Value != expectedFilamentMass {
		t.Errorf("FilamentMass wasn't updated correctly. Expected: %v, received: %v", expectedFilamentMass, curSpoolStats.FilamentMass.Value)
	}
}

// Test for GetValueFromMsg: Valid input msg
func TestGetValueFromMsg_ValidMsg(t *testing.T) {
	validMsg := []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x64} // Valid msg with temperature ID and value 100

	//Call function with valid Msg
	//Expected result: temperature value in data struct changed to 100
	GetValueFromMsg(validMsg)

	// Expected value
	expectedValue := float32(100)

	// Check if temperature value was updated
	if currentData.Temperature.Value != expectedValue {
		t.Errorf("Value wasn't updated correctly. Expected: %v, received: %v", expectedValue, currentData.Temperature.Value)
	}
}

// Test for GetValueFromMsg: Msg not 8bytes long
func TestGetValueFromMsg_ShortMsg(t *testing.T) {
	prevValue := float32(50)
	currentData.Temperature.Value = prevValue
	//Msg with invalid length
	invalidLenMsg := []byte{0x02, 0x00, 0x00}
	// Call function with msg of invalid number of bytes
	//Expected result: temperature value in data struct unchanged
	GetValueFromMsg(invalidLenMsg)
	if currentData.Temperature.Value != prevValue {
		t.Errorf("Invalid message shouldn't change current data value.")
	}
}

// Test for GetValueFromMsg: Unknown msg ID
func TestGetValueFromMsg_UnknownID(t *testing.T) {
	//Msg with undefined msg ID
	invalidIDMsg := []byte{0x99, 0x00, 0x00, 0x00, 0xDC, 0x00, 0x00, 0x00}
	// Save previous data content
	oldData := currentData
	// Call function with invalid ID msg
	//Expected result: temperature value in data struct unchanged
	GetValueFromMsg(invalidIDMsg)
	//Check if values have changed
	if currentData != oldData {
		t.Errorf("Invalid message changed values! Expected: %v, received: %v", oldData, currentData)
	}
}
