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
	var globalActual int
	var regionCount int
	var cityCount int
	var problems []problem

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
		countryCityTotal := 0
		countryComplete := true

		fmt.Printf("COUNTRY %s: expected boards=%d, regions=%d\n", country.Name, country.Count, len(regions))

		for _, region := range regions {
			regionTotal += region.Count
			regionLabel := region.Name
			if regionLabel == "" {
				regionLabel = "<blank>"
			}

			citiesRaw, err := client.Cities(ctx, country.Name, region.Name)
			if err != nil {
				countryComplete = false
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
			countryCityTotal += cityTotal

			status := "OK"
			if cityTotal != region.Count {
				status = "FAIL"
				countryComplete = false
				problems = append(problems, problem{
					Country: country.Name,
					Region:  region.Name,
					Detail:  fmt.Sprintf("city counts total %d, expected %d", cityTotal, region.Count),
				})
			}
			fmt.Printf("  REGION %-30s expected=%-4d cities=%-3d actual=%-4d %s\n", regionLabel, region.Count, len(cities), cityTotal, status)
		}

		countryStatus := "OK"
		if regionTotal != country.Count {
			countryStatus = "FAIL"
			countryComplete = false
			problems = append(problems, problem{Country: country.Name, Detail: fmt.Sprintf("region counts total %d, expected %d", regionTotal, country.Count)})
		}
		if countryComplete && countryCityTotal != country.Count {
			countryStatus = "FAIL"
			problems = append(problems, problem{Country: country.Name, Detail: fmt.Sprintf("city counts total %d, expected %d", countryCityTotal, country.Count)})
		}

		if countryComplete {
			fmt.Printf("  COUNTRY TOTAL expected=%d regions=%d cities=%d %s\n\n", country.Count, regionTotal, countryCityTotal, countryStatus)
			globalActual += countryCityTotal
		} else {
			fmt.Printf("  COUNTRY TOTAL expected=%d regions=%d cities=INCOMPLETE %s\n\n", country.Count, regionTotal, countryStatus)
		}
	}

	fmt.Println("=== GLOBAL SUMMARY ===")
	fmt.Printf("Countries:       %d\n", len(countries))
	fmt.Printf("Regions:         %d\n", regionCount)
	fmt.Printf("Cities/towns:    %d\n", cityCount)
	fmt.Printf("Expected boards: %d\n", globalExpected)
	fmt.Printf("Counted boards:  %d (complete countries only)\n", globalActual)
	fmt.Printf("Problems:        %d\n", len(problems))

	if len(problems) > 0 {
		fmt.Println("\n=== PROBLEMS ===")
		for _, p := range problems {
			where := p.Country
			if p.Region != "" {
				where += " / " + p.Region
			} else if strings.Contains(p.Detail, "fetch cities") {
				where += " / <blank>"
			}
			fmt.Printf("- %s: %s\n", where, p.Detail)
		}
		fmt.Println("Result:          FAIL")
		os.Exit(1)
	}

	if globalActual != globalExpected {
		fatalf("global city counts total %d, expected %d", globalActual, globalExpected)
	}
	fmt.Println("Result:          PASS")
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
