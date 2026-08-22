package model

import (
	"fmt"
	"math"
	"time"
)

// IslamicDate is a calculated Hijri date.
type IslamicDate struct {
	Day   int
	Month int // 1-12
	Year  int
}

var islamicMonthNames = [...]string{
	"Muharram",
	"Safar",
	"Rabi al-Awwal",
	"Rabi al-Thani",
	"Jumada al-Awwal",
	"Jumada al-Thani",
	"Rajab",
	"Sha'ban",
	"Ramadan",
	"Shawwal",
	"Dhul Qa'dah",
	"Dhul Hijjah",
}

// String returns the English date form used by MasjidPi displays.
func (d IslamicDate) String() string {
	if d.Day < 1 || d.Month < 1 || d.Month > len(islamicMonthNames) || d.Year < 1 {
		return ""
	}
	return fmt.Sprintf("%d %s %d", d.Day, islamicMonthNames[d.Month-1], d.Year)
}

// CalculateIslamicDate reproduces the civil/Kuwaiti Hijri calculation used by
// MasjidBoard Live. now must already be in the board's local timezone.
//
// MasjidBoard Live advances the Islamic date shortly after the board's sunset
// (185 seconds after the published sunset time), then applies the masjid's
// manual day adjustment. When forcePreviousMonthDay30 is true, the date is
// held on day 30 of the previous Hijri month for local moon-sighting use.
func CalculateIslamicDate(now time.Time, sunset *ClockTime, adjustment int, forcePreviousMonthDay30 bool) IslamicDate {
	local := now
	if adjustment != 0 {
		local = local.AddDate(0, 0, adjustment)
	}

	if sunset != nil {
		rollover := time.Date(
			now.Year(), now.Month(), now.Day(),
			sunset.Hour, sunset.Minute, 0, 0, now.Location(),
		).Add(185 * time.Second)
		if now.After(rollover) {
			local = local.AddDate(0, 0, 1)
		}
	}

	date, previousMonth := kuwaitiCalendar(local)
	if !forcePreviousMonthDay30 {
		return date
	}

	date.Day = 30
	date.Month = previousMonth + 1
	if date.Month < 1 {
		date.Month = 12
		date.Year--
	}
	return date
}

// kuwaitiCalendar is a direct Go port of the calendar arithmetic used by
// MasjidBoard Live's kuwaiticalendar() JavaScript function. previousMonth is
// zero-based and may be -1 at the Muharram boundary, matching the upstream
// intermediate value.
func kuwaitiCalendar(t time.Time) (date IslamicDate, previousMonth int) {
	day := t.Day()
	month := int(t.Month()) - 1
	year := t.Year()

	m := month + 1
	y := year
	if m < 3 {
		y--
		m += 12
	}

	a := math.Floor(float64(y) / 100)
	b := 2 - a + math.Floor(a/4)
	if y < 1583 {
		b = 0
	}
	if y == 1582 {
		if m > 10 {
			b = -10
		}
		if m == 10 {
			b = 0
			if day > 4 {
				b = -10
			}
		}
	}

	jd := math.Floor(365.25*float64(y+4716)) +
		math.Floor(30.6001*float64(m+1)) +
		float64(day) + b - 1524

	b = 0
	if jd > 2299160 {
		a = math.Floor((jd - 1867216.25) / 36524.25)
		b = 1 + a - math.Floor(a/4)
	}

	bb := jd + b + 1524
	cc := math.Floor((bb - 122.1) / 365.25)
	dd := math.Floor(365.25 * cc)
	ee := math.Floor((bb - dd) / 30.6001)

	const (
		islamicCycleDays = 10631.0
		islamicYearDays  = islamicCycleDays / 30.0
		epochAstronomical = 1948084.0
	)
	shift := 8.01 / 60.0
	z := jd - epochAstronomical
	cycle := math.Floor(z / islamicCycleDays)
	z -= islamicCycleDays * cycle
	j := math.Floor((z - shift) / islamicYearDays)
	hijriYear := int(30*cycle + j)
	z -= math.Floor(j*islamicYearDays + shift)
	hijriMonth := int(math.Floor((z + 28.5001) / 29.5))
	if hijriMonth == 13 {
		hijriMonth = 12
	}
	hijriDay := int(z - math.Floor(29.5001*float64(hijriMonth)-29))

	// Retain these Gregorian intermediate calculations because they are part of
	// the upstream algorithm and determine the Julian day above.
	_ = cc
	_ = dd
	_ = ee

	return IslamicDate{
		Day:   hijriDay,
		Month: hijriMonth,
		Year:  hijriYear,
	}, hijriMonth - 2
}
