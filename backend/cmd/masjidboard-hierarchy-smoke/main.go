package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/provider/masjidboardlive"
)

type problem struct {
	Country string
	Region  string
	Detail  string
}

type warning struct {
	Country string
	Region  string
	Detail  string
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	client := masjidboardlive.DiscoveryClient{}

	countriesRaw, err := client.Countries(ctx)
	if err != nil {
		fatalf("fetch countries: %v", err)
	}
	countries := merge(countriesRaw, false)

	fmt.Printf("Countries detected: %d\n", len(countries))
	fmt.Println()

	var globalExpected int
	var globalResolved int
	var globalUnresolved int
	var regionCount int
	var cityCount int
	var problems []problem
	var warnings []warning

	for _, country := range countries {
		globalExpected += country.Count

		regionsRaw, err := client.Regions(ctx, country.Name)
		if err != nil {
			fmt.Printf("COUNTRY %s: expected boards=%d\n", country.Name, country.Count)
			fmt.Printf("  ERROR fetch regions: %v\n\n", err)
			problems = append(problems, problem{Country: country.Name, Detail: fmt.Sprintf("fetch regions: %v", err)})
			continue
		}
		regions := merge(regionsRaw, true)
		regionCount += len(regions)

		regionTotal := 0
		countryResolved := 0
		countryUnresolved := 0

		fmt.Printf("COUNTRY %s: expected boards=%d, regions=%d\n", country.Name, country.Count, len(regions))

		for _, region := range regions {
			regionTotal += region.Count
			regionLabel := region.Name
			if regionLabel == "" {
				regionLabel = "<blank>"
			}

			// A single blank region means the country has no province layer and is
			// still queried normally. A blank bucket alongside named regions is an
			// upstream orphan bucket; preserve its count as unresolved rather than
			// calling the endpoint that has been observed to return HTTP 500.
			if region.Name == "" && len(regions) > 1 {
				countryUnresolved += region.Count
				warnings = append(warnings, warning{Country: country.Name, Region: region.Name, Detail: fmt.Sprintf("%d board(s) advertised in mixed blank region with no resolvable city bucket", region.Count)})
				fmt.Printf("  REGION %-30s expected=%-4d cities=%-3d resolved=%-4d unresolved=%-4d OK*\n", regionLabel, region.Count, 0, 0, region.Count)
				continue
			}

			citiesRaw, err := client.Cities(ctx, country.Name, region.Name)
			if err != nil {
				fmt.Printf("  REGION %-30s expected=%-4d ERROR: %v\n", regionLabel, region.Count, err)
				problems = append(problems, problem{Country: country.Name, Region: region.Name, Detail: fmt.Sprintf("fetch cities: %v", err)})
				continue
			}
			cities := merge(citiesRaw, false)
			cityCount += len(cities)

			cityTotal := 0
			for _, city := range cities {
				cityTotal += city.Count
			}
			if cityTotal > region.Count {
				problems = append(problems, problem{Country: country.Name, Region: region.Name, Detail: fmt.Sprintf("city counts total %d, exceed expected %d", cityTotal, region.Count)})
				fmt.Printf("  REGION %-30s expected=%-4d cities=%-3d resolved=%-4d unresolved=%-4d FAIL\n", regionLabel, region.Count, len(cities), cityTotal, 0)
				continue
			}

			unresolved := region.Count - cityTotal
			countryResolved += cityTotal
			countryUnresolved += unresolved
			status := "OK"
			if unresolved > 0 {
				status = "OK*"
				warnings = append(warnings, warning{Country: country.Name, Region: region.Name, Detail: fmt.Sprintf("%d board(s) counted by region but not represented by returned city rows", unresolved)})
			}
			fmt.Printf("  REGION %-30s expected=%-4d cities=%-3d resolved=%-4d unresolved=%-4d %s\n", regionLabel, region.Count, len(cities), cityTotal, unresolved, status)
		}

		countryStatus := "OK"
		if regionTotal != country.Count {
			countryStatus = "FAIL"
			problems = append(problems, problem{Country: country.Name, Detail: fmt.Sprintf("region counts total %d, expected %d", regionTotal, country.Count)})
		}
		if countryResolved+countryUnresolved != country.Count {
			countryStatus = "FAIL"
			problems = append(problems, problem{Country: country.Name, Detail: fmt.Sprintf("resolved %d plus unresolved %d = %d, expected %d", countryResolved, countryUnresolved, countryResolved+countryUnresolved, country.Count)})
		}
		if countryStatus == "OK" && countryUnresolved > 0 {
			countryStatus = "OK*"
		}
		fmt.Printf("  COUNTRY TOTAL expected=%d regions=%d resolved=%d unresolved=%d %s\n\n", country.Count, regionTotal, countryResolved, countryUnresolved, countryStatus)

		globalResolved += countryResolved
		globalUnresolved += countryUnresolved
	}

	fmt.Println("=== GLOBAL SUMMARY ===")
	fmt.Printf("Countries:         %d\n", len(countries))
	fmt.Printf("Regions:           %d\n", regionCount)
	fmt.Printf("Cities/towns:      %d\n", cityCount)
	fmt.Printf("Expected boards:   %d\n", globalExpected)
	fmt.Printf("Resolved boards:   %d\n", globalResolved)
	fmt.Printf("Unresolved boards: %d\n", globalUnresolved)
	fmt.Printf("Verified total:    %d\n", globalResolved+globalUnresolved)
	fmt.Printf("Warnings:          %d\n", len(warnings))
	fmt.Printf("Problems:          %d\n", len(problems))

	if len(warnings) > 0 {
		fmt.Println("\n=== UPSTREAM COVERAGE WARNINGS ===")
		for _, w := range warnings {
			where := w.Country
			if w.Region != "" {
				where += " / " + w.Region
			} else {
				where += " / <blank>"
			}
			fmt.Printf("- %s: %s\n", where, w.Detail)
		}
	}

	if len(problems) > 0 {
		fmt.Println("\n=== PROBLEMS ===")
		for _, p := range problems {
			where := p.Country
			if p.Region != "" {
				where += " / " + p.Region
			}
			fmt.Printf("- %s: %s\n", where, p.Detail)
		}
		fmt.Println("Result:            FAIL")
		os.Exit(1)
	}

	if globalResolved+globalUnresolved != globalExpected {
		fatalf("global resolved %d plus unresolved %d = %d, expected %d", globalResolved, globalUnresolved, globalResolved+globalUnresolved, globalExpected)
	}
	if globalUnresolved > 0 {
		fmt.Println("Result:            PASS (with upstream coverage warnings)")
		return
	}
	fmt.Println("Result:            PASS")
}

func merge(entries []masjidboardlive.HierarchyEntry, allowBlank bool) []masjidboardlive.HierarchyEntry {
	counts := make(map[string]int)
	for _, entry := range entries {
		name := strings.TrimSpace(entry.Name)
		if name == "" && !allowBlank {
			continue
		}
		counts[name] += entry.Count
	}

	out := make([]masjidboardlive.HierarchyEntry, 0, len(counts))
	for name, count := range counts {
		out = append(out, masjidboardlive.HierarchyEntry{Name: name, Count: count})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "FAIL: "+format+"\n", args...)
	os.Exit(1)
}
