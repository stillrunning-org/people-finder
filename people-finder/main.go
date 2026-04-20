package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	w "stillrunning.org/people-finder/people-finder/wikidata"
)

const (
	defaultStartYear = -5000
	defaultEndYear   = 2026
)

func formatWikidataYear(year int) string {
	if year < 0 {
		return fmt.Sprintf("-%04d", -year)
	}
	return fmt.Sprintf("%04d", year)
}

func yearStartDate(year int) string {
	return fmt.Sprintf("%s-01-01T00:00:00Z", formatWikidataYear(year))
}

func fetchWikidataDeathsWithRetry(date1, date2 string, maxAttempts int, retryDelay time.Duration) ([]w.Person, error) {
	if maxAttempts < 1 {
		maxAttempts = 1
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		people, err := w.FetchWikidataDeaths(date1, date2)
		if err == nil {
			return people, nil
		}
		lastErr = err
		if attempt < maxAttempts {
			time.Sleep(retryDelay)
		}
	}
	return nil, lastErr
}

func dedupePeopleByID(people []w.Person) []w.Person {
	if len(people) <= 1 {
		return people
	}

	seen := make(map[string]struct{}, len(people))
	unique := make([]w.Person, 0, len(people))
	for _, p := range people {
		if p.Id == "" {
			continue
		}
		if _, ok := seen[p.Id]; ok {
			continue
		}
		seen[p.Id] = struct{}{}
		unique = append(unique, p)
	}
	return unique
}

func main() {
	startYear := flag.Int("start-year", defaultStartYear, "starting year, inclusive")
	endYear := flag.Int("end-year", defaultEndYear, "ending year, inclusive")
	retries := flag.Int("retries", 3, "max retries per year query")
	retryDelayMs := flag.Int("retry-delay-ms", 1500, "delay between retries in milliseconds")
	requestDelayMs := flag.Int("request-delay-ms", 150, "delay between successful year queries in milliseconds")
	flag.Parse()

	if *startYear > *endYear {
		log.Fatalf("invalid range: start-year (%d) must be <= end-year (%d)", *startYear, *endYear)
	}

	db := initDatabase()
	defer db.Close()

	totalFetched := 0
	totalSaved := 0
	retryDelay := time.Duration(*retryDelayMs) * time.Millisecond
	requestDelay := time.Duration(*requestDelayMs) * time.Millisecond

	for year := *startYear; year <= *endYear; year++ {
		date1 := yearStartDate(year)
		date2 := yearStartDate(year + 1)

		people, err := fetchWikidataDeathsWithRetry(date1, date2, *retries, retryDelay)
		if err != nil {
			log.Printf("year %d: fetch failed after %d attempts: %v", year, *retries, err)
			continue
		}

		for i := range people {
			people[i].Age = calculateAgeAtDeath(people[i].BirthDate, people[i].DeathDate)
			if people[i].Age < 0 {
				people[i].Age = 0
			}
		}

		uniquePeople := dedupePeopleByID(people)
		if err := upsertPeople(db, uniquePeople); err != nil {
			log.Printf("year %d: db upsert failed: %v", year, err)
			continue
		}

		totalFetched += len(people)
		totalSaved += len(uniquePeople)
		fmt.Printf("year %d: fetched=%d saved=%d total_saved=%d\n", year, len(people), len(uniquePeople), totalSaved)

		if requestDelay > 0 {
			time.Sleep(requestDelay)
		}
	}

	fmt.Printf("completed: start_year=%d end_year=%d total_fetched=%d total_saved=%d\n", *startYear, *endYear, totalFetched, totalSaved)
}
