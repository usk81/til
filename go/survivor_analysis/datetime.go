package main

import "time"

type DateTime struct {
	time.Time
}

// MarshalCSV Converts the internal date as CSV string
func (date *DateTime) MarshalCSV() (string, error) {
	return date.Time.Format(time.RFC3339), nil
}

// You could also use the standard Stringer interface
func (date *DateTime) String() string {
	return date.String() // Redundant, just for example
}

// UnmarshalCSV Converts the CSV string as internal date
func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse(time.RFC3339, csv)
	return err
}
