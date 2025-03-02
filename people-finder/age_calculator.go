package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Function to validate and parse the date in the format YYYY-MM-DDTHH:MM:SSZ
func parseDate(date string) (time.Time, int, error) {
	// Regular expression to match date format: YYYY-MM-DDTHH:MM:SSZ (including BC years)
	datePattern := `^(-?[0-9]+)-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})Z$`
	re := regexp.MustCompile(datePattern)
	matches := re.FindStringSubmatch(date)

	if len(matches) != 7 {
		return time.Time{}, 0, fmt.Errorf("invalid date format")
	}

	// Parse the year, month, day, hour, minute, and second
	year, err := strconv.Atoi(matches[1])
	if err != nil {
		return time.Time{}, 0, err
	}

	month, err := strconv.Atoi(matches[2])
	if err != nil {
		return time.Time{}, 0, err
	}

	day, err := strconv.Atoi(matches[3])
	if err != nil {
		return time.Time{}, 0, err
	}

	hour, err := strconv.Atoi(matches[4])
	if err != nil {
		return time.Time{}, 0, err
	}

	minute, err := strconv.Atoi(matches[5])
	if err != nil {
		return time.Time{}, 0, err
	}

	second, err := strconv.Atoi(matches[6])
	if err != nil {
		return time.Time{}, 0, err
	}

	// Return time object for AD dates (positive years)
	if year >= 0 {
		parsedDate := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
		return parsedDate, year, nil
	}

	// For BC dates (negative years), return a year value to indicate it's BC
	// We will handle BC dates separately later
	return time.Time{}, year, nil
}

// Function to calculate the age at death
func calculateAgeAtDeath(birthDate, deathDate string) int {
	// Parse the birth and death dates
	birthTime, birthYear, err := parseDate(birthDate)
	if err != nil {
		return 0
	}

	deathTime, deathYear, err := parseDate(deathDate)
	if err != nil {
		return 0
	}

	// If either birth or death date is BC (negative year), handle it manually
	if birthYear < 0 || deathYear < 0 {
		// Handle BC year calculation separately
		if birthYear < 0 && deathYear < 0 {
			// Both birth and death are BC dates, calculate the difference in years
			age := -birthYear + deathYear
			return age
		}
		// If one date is BC, calculate the total age from BC to AD
		if birthYear < 0 {
			// Birth is BC, death is AD
			age := deathYear + (-birthYear)
			return age
		}
		if deathYear < 0 {
			// Death is BC, birth is AD
			age := -deathYear + birthYear
			return age
		}
	}

	// Calculate the age for AD dates
	if birthTime.Before(deathTime) {
		age := deathTime.Year() - birthTime.Year()
		if deathTime.YearDay() < birthTime.YearDay() {
			age--
		}
		return age
	}

	// Default case for error situations
	return 0
}
