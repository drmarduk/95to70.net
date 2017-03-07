package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Record holds a current weight
type Record struct {
	ID      int
	Created time.Time
	Value   float32
}

func (r *Record) String() string {
	return fmt.Sprintf("%d - %.2f (%s)", r.ID, r.Value, r.Created.Format(time.UnixDate))
}

// ParseRecord returns a fully parsed record based on the input string
// and adds the created member
func ParseRecord(value string) (r Record, err error) {
	if value == "" {
		return r, errors.New("empty value")
	}
	value = strings.Trim(value, " ")
	value = strings.Replace(value, ",", ".", -1)
	x, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return r, err
	}

	r.Created = time.Now()
	r.Value = float32(x)
	return
}
