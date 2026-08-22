package model

import (
	"testing"
	"time"
)

func TestCalculateIslamicDateMatchesMasjidBoardLiveReference(t *testing.T) {
	loc := time.FixedZone("GMT+02", 2*60*60)
	sunset := &ClockTime{Hour: 18, Minute: 0}
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, loc)

	got := CalculateIslamicDate(now, sunset, 0, false)
	want := IslamicDate{Day: 1, Month: 2, Year: 1448}
	if got != want {
		t.Fatalf("CalculateIslamicDate() = %+v, want %+v", got, want)
	}
	if got.String() != "1 Safar 1448" {
		t.Fatalf("IslamicDate.String() = %q", got.String())
	}
}

func TestCalculateIslamicDateRollsOverAfterMasjidBoardLiveSunsetDelay(t *testing.T) {
	loc := time.FixedZone("GMT+02", 2*60*60)
	sunset := &ClockTime{Hour: 18, Minute: 0}

	before := CalculateIslamicDate(time.Date(2026, 7, 16, 18, 3, 5, 0, loc), sunset, 0, false)
	if before != (IslamicDate{Day: 1, Month: 2, Year: 1448}) {
		t.Fatalf("at rollover boundary = %+v, want 1 Safar 1448", before)
	}

	after := CalculateIslamicDate(time.Date(2026, 7, 16, 18, 3, 6, 0, loc), sunset, 0, false)
	if after != (IslamicDate{Day: 2, Month: 2, Year: 1448}) {
		t.Fatalf("after rollover = %+v, want 2 Safar 1448", after)
	}
}

func TestCalculateIslamicDateAppliesMasjidAdjustment(t *testing.T) {
	loc := time.FixedZone("GMT+02", 2*60*60)
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, loc)

	plusOne := CalculateIslamicDate(now, nil, 1, false)
	if plusOne != (IslamicDate{Day: 2, Month: 2, Year: 1448}) {
		t.Fatalf("+1 adjustment = %+v, want 2 Safar 1448", plusOne)
	}

	minusOne := CalculateIslamicDate(now, nil, -1, false)
	if minusOne != (IslamicDate{Day: 30, Month: 1, Year: 1448}) {
		t.Fatalf("-1 adjustment = %+v, want 30 Muharram 1448", minusOne)
	}
}

func TestCalculateIslamicDateCanForcePreviousMonthDay30(t *testing.T) {
	loc := time.FixedZone("GMT+02", 2*60*60)
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, loc)

	got := CalculateIslamicDate(now, nil, 0, true)
	want := IslamicDate{Day: 30, Month: 1, Year: 1448}
	if got != want {
		t.Fatalf("forced date = %+v, want %+v", got, want)
	}
}

func TestCalculateIslamicDateForceDay30WrapsMuharramToPreviousYear(t *testing.T) {
	loc := time.FixedZone("GMT+02", 2*60*60)
	now := time.Date(2026, 6, 16, 12, 0, 0, 0, loc)

	got := CalculateIslamicDate(now, nil, 0, true)
	want := IslamicDate{Day: 30, Month: 12, Year: 1447}
	if got != want {
		t.Fatalf("forced Muharram boundary = %+v, want %+v", got, want)
	}
}
