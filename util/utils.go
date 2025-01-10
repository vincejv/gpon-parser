package util

import (
	"errors"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
)

func NormalizeString(s string) string {
	return strings.TrimSpace(strings.Join(strings.Fields(s), " "))
}

func RemoveLastNChars(s string, lengthNChars int) string {
	return s[:len(s)-lengthNChars]
}

func ParseDuration(timeString string) int64 {
	var days, hours, minutes, seconds int64

	// Split the input string into components (separated by whitespace)
	components := strings.Fields(timeString)

	// Iterate through each component
	for _, component := range components {
		var unit string
		var valueStr string

		// Handle cases with space between number and unit (like "1 d", "3 h", "32 m", "5 s")
		if len(component) > 1 && (component[len(component)-1] == 'd' || component[len(component)-1] == 'h' || component[len(component)-1] == 'm' || component[len(component)-1] == 's') {
			// This is the case where the numeric value and the unit are separated by a space
			unit = string(component[len(component)-1]) // 'd', 'h', 'm', or 's'
			valueStr = component[:len(component)-1]
		} else {
			// Handle full units (like "days", "hours", "minutes", "seconds")
			if strings.HasSuffix(component, "days") {
				unit = "d"
				valueStr = component[:len(component)-4]
			} else if strings.HasSuffix(component, "hours") {
				unit = "h"
				valueStr = component[:len(component)-5]
			} else if strings.HasSuffix(component, "minutes") {
				unit = "m"
				valueStr = component[:len(component)-7]
			} else if strings.HasSuffix(component, "seconds") {
				unit = "s"
				valueStr = component[:len(component)-7]
			}
		}

		// Parse the numeric part into an integer
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			continue // Skip invalid parts
		}

		// Update the respective duration component based on the unit
		switch unit {
		case "d":
			days = value
		case "h":
			hours = value
		case "m":
			minutes = value
		case "s":
			seconds = value
		}
	}

	// Convert everything to seconds and return the total
	return days*86400 + hours*3600 + minutes*60 + seconds
}

func RandInt(min int, max int) int {
	return rand.Intn(max-min) + min
}

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func ParseHtmlPage(elemList *[]string, page *DocPage, xpath string) error {
	htmlNode := htmlquery.FindOne(page.GetPage(), xpath)
	if htmlNode != nil {
		*elemList = append(*elemList, NormalizeString(htmlquery.InnerText(htmlNode)))
		return nil
	}
	return errors.New("unable to find xpath: " + xpath)
}

func ParseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func ParseInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func ParseFloat(s string) float64 {
	if s == "" {
		return 0 // Return a default value
	}
	// Remove non-numeric characters (e.g., dBm, V, mA)
	s = strings.TrimSpace(s)
	for _, ch := range s {
		if (ch < '0' || ch > '9') && ch != '.' {
			s = strings.TrimSuffix(s, string(ch)) // Remove characters after the number
		}
	}

	// Parse the cleaned string
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0 // Return 0 if parsing fails
	}
	return f
}

// ConvertPowerToDBm converts raw power value to dBm
func ConvertPowerToDBm(power string) (float64, error) {
	// Convert string to number
	powerVal, err := strconv.ParseFloat(power, 64)
	if err != nil {
		return 0, err // Return 0 and the error if conversion fails
	}

	// Perform the conversion: log10(power / 10000)
	convertedPower := math.Log10(powerVal / 1e4)

	// Approximate the result
	convertedPower = math.Round(convertedPower*100000) / 10000

	// Return the raw float value instead of a formatted string
	return convertedPower, nil
}

// convertWorkTemperature converts the WorkTemperature value to a float.
func ConvertWorkTemperature(workTemperature string) (float64, error) {
	// Convert the input string to a number
	r, err := strconv.ParseFloat(workTemperature, 64)
	if err != nil {
		return 0, err // Return 0 and the error if conversion fails
	}

	// Calculate the temperature based on the conditions
	var temperature float64
	if r >= math.Pow(2, 15) {
		temperature = -((math.Pow(2, 16) - r) / 256) // No rounding to keep decimal precision
	} else {
		temperature = r / 256 // No rounding to keep decimal precision
	}

	return temperature, nil
}
