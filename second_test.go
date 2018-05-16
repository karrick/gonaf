package gonaf_test

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"strconv"
	"testing"
)

// The gist is to compose a JSON decoder with something that interprets null and infinite tokens as
// their corresponding floating point equivalents, and to compose a JSON encoder with something that
// represetnts null and infinite values as their corresponding tokens.

func TestSecondEncoder(t *testing.T) {
	// takes normal JSON values, and encodes into JSON
}

type Encoder interface {
	Encode(interface{}) error
}

type encoder struct {
	iow io.Writer
	err error
}

func NewEncoder(w io.Writer) Encoder {
	return &encoder{iow: w}
}

func (e encoder) Encode(value interface{}) error {
	if e.err != nil {
		return e.err
	}

	var buf []byte

	if v, ok := value.(float64); ok {
		if math.IsNaN(v) {
			buf = []byte("null")
		} else if math.IsInf(v, 1) {
			buf = []byte("1e999")
		} else if math.IsInf(v, -1) {
			buf = []byte("-1e999")
		} else {
			// NOTE: magic parameters in call to strconv.FormatFloat
			buf = []byte(strconv.FormatFloat(v, 'g', -1, 64))
		}
		_, e.err = e.iow.Write(buf)
		return e.err
	}

	buf, e.err = json.Marshal(value)
	if e.err == nil {
		_, e.err = e.iow.Write(buf)
	}
	return e.err
}

func TestRoundTripADTSecond(t *testing.T) {
	source := struct {
		Name   string
		Values []float64
	}{
		Name: "foo",
		Values: []float64{
			math.Pi,
			math.NaN(),
			math.Inf(1),
			math.Inf(-1),
		},
	}
	bb := new(bytes.Buffer)
	encoder := NewEncoder(bb)

	if err := encoder.Encode(source); err != nil {
		t.Fatalf("Actual: %#v; Expected: %#v", err, nil)
	}

	t.Logf(bb.String())

	decoder := json.NewDecoder(bb)
	var destination testType

	if err := decoder.Decode(&destination); err != nil {
		t.Fatalf("Actual: %#v; Expected: %#v", err, nil)
	}
	if actual, expected := len(destination.Values), len(source.Values); actual != expected {
		t.Fatalf("Actual: %#v; Expected: %#v", actual, expected)
	}
	if actual, expected := destination.Values[0], source.Values[0]; float64(actual) != expected {
		t.Errorf("Actual: %#v; Expected: %#v", actual, expected)
	}
	if actual, expected := float64(destination.Values[1]), float64(source.Values[1]); !math.IsNaN(actual) {
		t.Errorf("Actual: %#v; Expected: %#v", actual, expected)
	}
	if actual, expected := float64(destination.Values[2]), float64(source.Values[2]); !math.IsInf(actual, 1) {
		t.Errorf("Actual: %#v; Expected: %#v", actual, expected)
	}
	if actual, expected := float64(destination.Values[3]), float64(source.Values[3]); !math.IsInf(actual, -1) {
		t.Errorf("Actual: %#v; Expected: %#v", actual, expected)
	}
}
