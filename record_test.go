package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

// TestRecordString tests the toString() function
func TestRecordString(t *testing.T) {
	tmptime := time.Now()
	r := Record{ID: 1, Created: tmptime, Value: 12.3}

	got := r.String()
	expected := fmt.Sprintf("%d - %.2f (%s)", r.ID, r.Value, r.Created.Format(time.UnixDate))

	if got != expected {
		t.Fatalf("Record.String(): got %s, expected %s\n", got, expected)
	}
}

func TestParseRecord(t *testing.T) {
	tests := []struct {
		in  string
		out Record
		err error
	}{
		{"2", Record{Value: 2.0}, nil},
		{"1.2", Record{Value: 1.2}, nil},
		{"2,3", Record{Value: 2.3}, nil},
		{" -1", Record{Value: -1.0}, nil},
		{"abc", Record{}, strconv.ErrSyntax},
		{"", Record{}, ErrEmptyValue},
	}

	for _, tt := range tests {
		got, err := ParseRecord(tt.in)

		if (tt.err != nil) && err.Error() != tt.err.Error() {
			t.Errorf("ParseRecord(%s): want %v, got %v\n", tt.in, tt.err, err)
		}
		if got.Value == tt.out.Value && err == tt.err {
			// success
			continue
		}
	}
}

func TestRecord(t *testing.T) {
	_id := 1337
	_created := time.Now()
	_value := float32(-0.1)

	r := Record{ID: _id, Created: _created, Value: _value}

	if r.ID != _id || r.Created != _created || r.Value != _value {
		t.Errorf("Record{} does not do good stuff")
	}
}
