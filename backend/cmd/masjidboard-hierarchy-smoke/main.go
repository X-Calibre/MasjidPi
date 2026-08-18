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

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
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

	for _, country := range countries {
		globalExpected += country.Count

		regionsRaw, err := client.Regions(ctx, country.Name)
		if err != nil {
			fatalf("fetch regions for %q: %v", country.Name, err)
		}
		regions := merge(regionsRaw, true)
		regionCount += len(regions)

		regionTotal := 0
		countryCityTotal := 0

		fmt.Printf("COUNTRY %s: expected boards=%d, regions=%d\n", country.Name, country.Count, len(regions))

		for _, region := range regions {
			regionTotal += region.Count

			citiesRaw, err := client.Cities(ctx, country.Name, region.Name)
			if err != nil {
				fatalf("fetch cities for %q / %q: %v", country.Name, region.Name, err)
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
			}
			regionLabel := region.Name
			if regionLabel == "" {
				regionLabel = "<blank>"
			}
			fmt.Printf("  REGION %-30s expected=%-4d cities=%-3d actual=%-4d %s\n", regionLabel, region.Count, len(cities), cityTotal, status)

			if cityTotal != region.Count {
				fatalf("city counts for %q / %q total %d, expected %d", country.Name, region.Name, cityTotal, region.Count)
			}
		}

		countryStatus := "OK"
		if regionTotal != country.Count || countryCityTotal != country.Count {
			countryStatus = "FAIL"
		}
		fmt.Printf("  COUNTRY TOTAL expected=%d regions=%d cities=%d %s\n\n", country.Count, regionTotal, countryCityTotal, countryStatus)

		if regionTotal != country.Count {
			fatalf("region counts for %q total %d, expected %d", country.Name, regionTotal, country.Count)
		}
		if countryCityTotal != country.Count {
			fatalf("city counts for %q total %d, expected %d", country.Name, countryCityTotal, country.Count)
		}

		globalActual += countryCityTotal
	}

	fmt.Println("=== GLOBAL SUMMARY ===")
	fmt.Printf("Countries:       %d\n", len(countries))
	fmt.Printf("Regions:         %d\n", regionCount)
	fmt.Printf("Cities/towns:    %d\n", cityCount)
	fmt.Printf("Expected boards: %d\n", globalExpected)
	fmt.Printf("Counted boards:  %d\n", globalActual)

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
