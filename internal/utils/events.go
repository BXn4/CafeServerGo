/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package utils

import (
	"errors"
	"math"
	"time"
)

// EasterYear returns time.Time for midnight on Orthodox Easter of a given year
// use Meeus Julian algorithm and return time based on Gregorian calendar
//
// @param int year The year as a number greater than 325
// @return time.Time The easter date as a time.Time
// @return errors.Error Error if exists, else nil
// from: https://github.com/vjeantet/eastertime
func EasterYear(year int) (time.Time, error) {

	if year < 326 {
		return time.Now(), errors.New("year have to be greater than 325")
	}

	var a, b, c, d, e int
	var month time.Month
	var day float64

	a = year % 4
	b = year % 7
	c = year % 19
	d = (19*c + 15) % 30
	e = (2*a + 4*b - d + 34) % 7
	month = time.Month((d + e + 114) / 31)
	day = math.Floor(float64((d+e+114)%31 + 1))
	day = day + 13

	return time.Date(year, month, int(day), 0, 0, 0, 0, time.UTC), nil
}

func IsEvent(t time.Time) bool {
	y, _, _ := t.UTC().Date()

	easter, _ := EasterYear(y)
	halloween := time.Date(y, time.October, 31, 0, 0, 0, 0, time.UTC)
	christmas := time.Date(y, time.December, 25, 0, 0, 0, 0, time.UTC)

	holidays := []time.Time{easter, halloween, christmas}

	for _, h := range holidays {
		days := int(t.Sub(h).Hours() / 24)

		switch {
		case h.Month() == time.December && h.Day() == 25:
			// christmas starts 15 days before and goes until Dec 31
			if days >= -15 {
				ty, tm, td := t.Date()
				if tm == time.December && ty == y && td <= 31 {
					return true
				}
			}

		default:
			// other holidays: ±7 days
			if days >= -7 && days <= 7 {
				return true
			}
		}
	}

	return false
}

func GetDaysLeft(t time.Time) int {
	y, _, _ := t.UTC().Date()

	easter, _ := EasterYear(y)
	halloween := time.Date(y, time.October, 31, 0, 0, 0, 0, time.UTC)
	christmas := time.Date(y, time.December, 25, 0, 0, 0, 0, time.UTC)

	holidays := []struct {
		date        time.Time
		startOffset int
		endOffset   int
	}{
		{easter, -7, 7},
		{halloween, -7, 7},
		{christmas, -15, 6}, // starts 15 days earlier, finish always on 31.
	}

	for _, h := range holidays {
		days := int(t.Sub(h.date).Hours() / 24)
		if days >= h.startOffset && days <= h.endOffset {
			return h.endOffset - days
		}
	}

	return 0
}

/*
HALLOWEEN = 2
CHRISTMAS = 3
EASTER = 5
*/
func GetEventType(t time.Time) int {
	y, _, _ := t.UTC().Date()

	easter, _ := EasterYear(y)
	halloween := time.Date(y, time.October, 31, 0, 0, 0, 0, time.UTC)
	christmas := time.Date(y, time.December, 25, 0, 0, 0, 0, time.UTC)

	holidays := []struct {
		id          int
		date        time.Time
		startOffset int
		endOffset   int
	}{
		{5, easter, -7, 7},
		{2, halloween, -7, 7},
		{3, christmas, -15, 6}, // starts 15 days earlier, finish always on 31.
	}

	for _, h := range holidays {
		days := int(t.Sub(h.date).Hours() / 24)
		if days >= h.startOffset && days <= h.endOffset {
			return h.id
		}
	}

	return 0
}
