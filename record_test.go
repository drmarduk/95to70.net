package main

import (
	"errors"
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
		{"", Record{}, errors.New("empty value")},
	}

	for _, tt := range tests {
		got, err := ParseRecord(tt.in)

		if tt.in == "" && err != errors.New("empty value") {
			t.Errorf("ParseRecord(%s) = '%s', want error '%v'", tt.in, err, errors.New("empty value"))
		}
		if got == tt.out && err == tt.err {
			continue
		}
		if got == tt.out && err != tt.err {
			t.Errorf("ParseRecord(%s): '%v', want error '%v'", tt.in, err, tt.err)
		}
		if got.Value != tt.out.Value {
			t.Errorf("ParseRecord(%s) = %v, want = %v\n%v\n", tt.in, got.Value, tt.out.Value, err)
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
