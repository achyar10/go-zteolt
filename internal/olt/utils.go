package olt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ConvertStringToUint16 Convert String to Uint16
func ConvertStringToUint16(str string) uint16 {
	// Convert string to uint16
	value, err := strconv.ParseUint(str, 10, 16)
	if err != nil {
		return 0
	}

	return uint16(value)
}

// ConvertStringToInteger Convert String to Integer
func ConvertStringToInteger(str string) int {
	// Convert string to int
	value, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return value
}

// ConvertDurationToString Convert duration to human-readable format
func ConvertDurationToString(duration time.Duration) string {
	days := int(duration / (24 * time.Hour))
	duration = duration % (24 * time.Hour)
	hours := int(duration / time.Hour)
	duration = duration % time.Hour
	minutes := int(duration / time.Minute)
	duration = duration % time.Minute
	seconds := int(duration / time.Second)

	return strconv.Itoa(days) + " days " + strconv.Itoa(hours) + " hours " + strconv.Itoa(minutes) + " minutes " + strconv.Itoa(seconds) + " seconds"
}

// ConvertByteArrayToDateTime Convert byte array to human-readable date time
func ConvertByteArrayToDateTime(byteArray []byte) (string, error) {

	// Check if byteArray length is exactly 8
	if len(byteArray) != 8 {
		return "", errors.New("invalid byte array length: expected 8 bytes")
	}

	// Extract the year from the first two bytes
	year := int(binary.BigEndian.Uint16(byteArray[0:2]))

	// Extract other components
	month := time.Month(byteArray[2]) // Month
	day := int(byteArray[3])          // Day
	hour := int(byteArray[4])         // Hour
	minute := int(byteArray[5])       // Minute
	second := int(byteArray[6])       // Second

	// Validate extracted values
	if month < 1 || month > 12 {
		return "", fmt.Errorf("invalid month: %d", month)
	}
	if day < 1 || day > 31 {
		return "", fmt.Errorf("invalid day: %d", day)
	}
	if hour < 0 || hour > 23 {
		return "", fmt.Errorf("invalid hour: %d", hour)
	}
	if minute < 0 || minute > 59 {
		return "", fmt.Errorf("invalid minute: %d", minute)
	}
	if second < 0 || second > 59 {
		return "", fmt.Errorf("invalid second: %d", second)
	}

	// Create a time.Time object UTC
	datetime := time.Date(year, month, day, hour, minute, second, 0, time.UTC)

	// Convert to Unix epoch time (seconds since Jan 1, 1970)
	return datetime.Format("2006-01-02 15:04:05"), nil
}

func ExtractONUID(oid string) string {
	// Split the OID name and take the last component
	parts := strings.Split(oid, ".")
	if len(parts) > 0 {
		// Check if the last component is a valid number
		lastComponent := parts[len(parts)-1]
		if _, err := strconv.Atoi(lastComponent); err == nil {
			return lastComponent
		}
	}
	return "" // Return an empty string if the OID is invalid or empty (default value)
}

func ExtractIDOnuID(oid interface{}) int {
	if oid == nil {
		return 0
	}

	switch v := oid.(type) {
	case string:
		parts := strings.Split(v, ".")
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]
			id, err := strconv.Atoi(lastPart)
			if err == nil {
				return id
			}
		}
		return 0
	default:
		return 0
	}
}

func ExtractName(oidValue interface{}) string {
	switch v := oidValue.(type) {
	case string:
		// Data is string, return it
		return v
	case []byte:
		// Data is byte slice, convert to string
		return string(v)
	default:
		// Data type is not recognized, you can handle this case according to your needs.
		return "Unknown" // Return "Unknown" if the OID is invalid or empty
	}
}

// ExtractSerialNumber function is used to extract serial number from OID value
func ExtractSerialNumber(oidValue interface{}) string {
	switch v := oidValue.(type) {
	case string:
		// If the string starts with "1,", remove it from the string
		if strings.HasPrefix(v, "1,") {
			return v[2:]
		}
		return v
	case []byte:
		// Convert byte slice to string
		strValue := string(v)
		if strings.HasPrefix(strValue, "1,") {
			return strValue[2:]
		}
		return strValue // Data is byte slice, convert to string
	default:
		// Data type is not recognized, you can handle this case according to your needs.
		return "" // Return 0 if the OID is invalid or empty (default value)
	}
}

func ConvertAndMultiply(pduValue interface{}) (string, error) {
	// Type assert pduValue to an integer type
	intValue, ok := pduValue.(int)
	if !ok {
		return "", fmt.Errorf("value is not an integer")
	}

	// Multiply the integer by 0.002
	result := float64(intValue) * 0.002

	// Subtract 30
	result -= 30.0

	// Convert the result to a string with two decimal places
	resultStr := strconv.FormatFloat(result, 'f', 2, 64)

	return resultStr, nil
}

func ExtractAndGetStatus(oidValue interface{}) string {
	// Check if oidValue is not an integer
	intValue, ok := oidValue.(int)
	if !ok {
		return "Unknown"
	}

	switch intValue {
	case 1:
		return "Logging"
	case 2:
		return "LOS"
	case 3:
		return "Synchronization"
	case 4:
		return "Online"
	case 5:
		return "Dying Gasp"
	case 6:
		return "Auth Failed"
	case 7:
		return "Offline"
	default:
		return "Unknown"
	}
}

func ExtractLastOfflineReason(oidValue interface{}) string {
	// Check if oidValue is not an integer
	intValue, ok := oidValue.(int)
	if !ok {
		return "Unknown"
	}

	switch intValue {
	case 1:
		return "Unknown"
	case 2:
		return "LOS"
	case 3:
		return "LOSi"
	case 4:
		return "LOFi"
	case 5:
		return "sfi"
	case 6:
		return "loai"
	case 7:
		return "loami"
	case 8:
		return "AuthFail"
	case 9:
		return "PowerOff"
	case 10:
		return "deactiveSucc"
	case 11:
		return "deactiveFail"
	case 12:
		return "Reboot"
	case 13:
		return "Shutdown"
	default:
		return "Unknown"
	}
}

func ExtractGponOpticalDistance(oidValue interface{}) string {
	// Check if oidValue is not an integer
	intValue, ok := oidValue.(int)
	if !ok {
		return "Unknown"
	}

	return strconv.Itoa(intValue)
}