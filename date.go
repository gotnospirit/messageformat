package messageformat

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

// DateFormatter is an interface for a type that formats
// a date in a variety of formats.
type DateFormatter interface {
	Short(time.Time) string
	Medium(time.Time) string
	Long(time.Time) string
	Full(time.Time) string
}

type AmericanDateFormatter struct {
	Month   map[time.Month]string
	Weekday map[time.Weekday]string
}

func createAmericanDateFormatter() DateFormatter {
	return &AmericanDateFormatter{
		Month: map[time.Month]string{
			time.January:   "January",
			time.February:  "February",
			time.March:     "March",
			time.April:     "April",
			time.May:       "May",
			time.June:      "June",
			time.July:      "July",
			time.August:    "August",
			time.September: "September",
			time.October:   "October",
			time.November:  "November",
			time.December:  "December",
		},
		Weekday: map[time.Weekday]string{
			time.Monday:    "Monday",
			time.Tuesday:   "Tuesday",
			time.Wednesday: "Wednesday",
			time.Thursday:  "Thursday",
			time.Friday:    "Friday",
			time.Saturday:  "Saturday",
			time.Sunday:    "Sunday",
		},
	}
}

// 10/16/1996
func (en *AmericanDateFormatter) Short(t time.Time) string {
	return fmt.Sprintf("%d/%d/%d", t.Month(), t.Day(), t.Year())
}

// October 16, 1996
func (en *AmericanDateFormatter) Medium(t time.Time) string {
	return fmt.Sprintf("%s %d, %d", en.Month[t.Month()], t.Day(), t.Year())
}

// Tuesday October 16, 1996
func (en *AmericanDateFormatter) Long(t time.Time) string {
	return fmt.Sprintf("%s %s %d, %d", en.Weekday[t.Weekday()], en.Month[t.Month()], t.Day(), t.Year())
}

func (en *AmericanDateFormatter) Full(t time.Time) string {
	// TODO: implement format
	return fmt.Sprintf("%s %d. %s %d", en.Weekday[t.Weekday()], t.Day(), en.Month[t.Month()], t.Year())
}

type GermanDateFormatter struct {
	Month   map[time.Month]string
	Weekday map[time.Weekday]string
}

func createGermanDateFormatter() DateFormatter {
	return &GermanDateFormatter{
		Month: map[time.Month]string{
			time.January:   "Januar",
			time.February:  "Februar",
			time.March:     "März",
			time.April:     "April",
			time.May:       "Mai",
			time.June:      "Juni",
			time.July:      "Juli",
			time.August:    "August",
			time.September: "September",
			time.October:   "Oktober",
			time.November:  "November",
			time.December:  "Dezember",
		},
		Weekday: map[time.Weekday]string{
			time.Monday:    "Montag",
			time.Tuesday:   "Dienstag",
			time.Wednesday: "Mittwoch",
			time.Thursday:  "Donnerstag",
			time.Friday:    "Freitag",
			time.Saturday:  "Samstag",
			time.Sunday:    "Sonntag",
		},
	}
}

func (de *GermanDateFormatter) Short(t time.Time) string {
	return fmt.Sprintf("%d.%d.%d", t.Day(), t.Month(), t.Year())
}

func (de *GermanDateFormatter) Medium(t time.Time) string {
	return fmt.Sprintf("%d. %s %d", t.Day(), de.Month[t.Month()], t.Year())
}

func (de *GermanDateFormatter) Long(t time.Time) string {
	return fmt.Sprintf("%s %d. %s %d", de.Weekday[t.Weekday()], t.Day(), de.Month[t.Month()], t.Year())
}

func (de *GermanDateFormatter) Full(t time.Time) string {
	// TODO: implement format
	return fmt.Sprintf("%s %d. %s %d", de.Weekday[t.Weekday()], t.Day(), de.Month[t.Month()], t.Year())
}

type DateExpr struct {
	// Key is the key for the data map when this
	// expression is being formatted.
	Key string `json:"key"`
	// Format represents the DateFormat enum value
	Format DateFormat `json:"format"`
}

type DateFormat = string

const (
	Short  DateFormat = "short"
	Medium DateFormat = "medium"
	Long   DateFormat = "long"
	Full   DateFormat = "full"
	// skeleton represents ICU's datetime skeleton format
	// Skeleton DateFormat = "skeleton"
)

// parseDate attempts to parse the input at the given start position into a DateExpr
func (p *parser) parseDate(varName string, nextChar rune, start, end int, ptr_input *[]rune) (Expression, int, error) {
	var result = DateExpr{
		Key: varName,
	}

	format, _, cursor, err := readVar(start+1, end, ptr_input)
	if err != nil {
		return nil, cursor, errors.New("failed to parse date format")
	}

	switch string(format) {
	case Short:
		result.Format = Short
	case Medium:
		result.Format = Medium
	case Long:
		result.Format = Long
	case Full:
		result.Format = Full
	default:
		return nil, cursor, fmt.Errorf("InvalidDateFormat")
	}

	return &result, cursor, nil
}

func (f *formatter) formatDate(expr Expression, ptrOutput *bytes.Buffer, data map[string]any) error {
	if date, ok := expr.(*DateExpr); ok {
		t, ok := data[date.Key].(time.Time)
		if !ok {
			return fmt.Errorf("InvalidArgType: want time.Time, got %T", t)
		}

		switch date.Format {
		case Short:
			ptrOutput.WriteString(f.date.Short(t))
		case Medium:
			ptrOutput.WriteString(f.date.Medium(t))
		case Long:
			ptrOutput.WriteString(f.date.Long(t))
		case Full:
			ptrOutput.WriteString(f.date.Full(t))
		default:
			return fmt.Errorf("InvalidDateFormat")
		}
	} else {
		return fmt.Errorf("InvalidExprType: want DateExpr, got %T", expr)
	}

	return nil
}

// Symbol represents a format symbol
type Symbol rune

const (
	Era        Symbol = 'G'
	Year       Symbol = 'y'
	ShortMonth Symbol = 'M'
	LongMonth  Symbol = 'L'
	DayOfMonth Symbol = 'd'
	DayOfWeek  Symbol = 'E'
	AmPmMarker Symbol = 'a'
	Hour112    Symbol = 'h'
	Hour023    Symbol = 'H'
	Hour011    Symbol = 'K'
	Hour124    Symbol = 'k'
	Minute     Symbol = 'm'
	Second     Symbol = 's'
	TimeZone   Symbol = 'Z'
)

// DateTimeSkeleton is an ICU datetime skeleton.
type DateTimeSkeleton struct {
}
