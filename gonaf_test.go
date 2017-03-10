package gonaf_test

import (
	"bytes"
	"encoding/json"
	"math"
	"testing"

	"github.com/karrick/gonaf"
)

func TestJsonFloatDecodeNaN(t *testing.T) {
	bb := bytes.NewBufferString("null")
	decoder := json.NewDecoder(bb)
	var jf gonaf.JsonFloat

	if err := decoder.Decode(&jf); err != nil {
		t.Fatalf("Actual: %#v; Expected: %#v", err, nil)
	}
	if !math.IsNaN(float64(jf)) {
		t.Fatalf("Actual: %#v; Expected: %#v", jf, math.NaN())
	}
}

func TestJsonFloatDecodePositiveInfinity(t *testing.T) {
	bb := bytes.NewBufferString("1e999")
	decoder := json.NewDecoder(bb)
	var jf gonaf.JsonFloat

	if err := decoder.Decode(&jf); err != nil {
		t.Fatalf("Actual: %#v; Expected: %#v", err, nil)
	}
	if !math.IsInf(float64(jf), 1) {
		t.Fatalf("Actual: %#v; Expected: %#v", jf, math.Inf(1))
	}
}

func TestJsonFloatDecodeNegativeInfinity(t *testing.T) {
	bb := bytes.NewBufferString("-1e999")
	decoder := json.NewDecoder(bb)
	var jf gonaf.JsonFloat

	if err := decoder.Decode(&jf); err != nil {
		t.Fatalf("Actual: %#v; Expected: %#v", err, nil)
	}
	if !math.IsInf(float64(jf), -1) {
		t.Fatalf("Actual: %#v; Expected: %#v", jf, math.Inf(-1))
	}
}

func TestJsonFloatEncodeNaN(t *testing.T) {
	jf := gonaf.JsonFloat(math.NaN())
	bb := new(bytes.Buffer)
	encoder := json.NewEncoder(bb)

	if err := encoder.Encode(jf); err != nil {
		t.Fatalf("Actual: %#v; Expected: %#v", err, nil)
	}
	// NOTE: JSON encoding is free to add leading and trailing whitespace to the encoded value
	if actual, expected := string(bytes.TrimSpace(bb.Bytes())), "null"; actual != expected {
		t.Fatalf("Actual: %#v; Expected: %#v", actual, expected)
	}
}

func TestJsonFloatEncodePositiveInfinity(t *testing.T) {
	jf := gonaf.JsonFloat(math.Inf(1))
	bb := new(bytes.Buffer)
	encoder := json.NewEncoder(bb)

	if err := encoder.Encode(jf); err != nil {
		t.Fatalf("Actual: %#v; Expected: %#v", err, nil)
	}
	// NOTE: JSON encoding is free to add leading and trailing whitespace to the encoded value
	if actual, expected := string(bytes.TrimSpace(bb.Bytes())), "1e999"; actual != expected {
		t.Fatalf("Actual: %#v; Expected: %#v", actual, expected)
	}
}

func TestJsonFloatEncodeNegativeInfinity(t *testing.T) {
	jf := gonaf.JsonFloat(math.Inf(-1))
	bb := new(bytes.Buffer)
	encoder := json.NewEncoder(bb)

	if err := encoder.Encode(jf); err != nil {
		t.Fatalf("Actual: %#v; Expected: %#v", err, nil)
	}
	// NOTE: JSON encoding is free to add leading and trailing whitespace to the encoded value
	if actual, expected := string(bytes.TrimSpace(bb.Bytes())), "-1e999"; actual != expected {
		t.Fatalf("Actual: %#v; Expected: %#v", actual, expected)
	}
}

type testType struct {
	Name   string
	Values []gonaf.JsonFloat
}

func TestRoundTripADT(t *testing.T) {
	source := testType{
		Name: "foo",
		Values: []gonaf.JsonFloat{
			gonaf.JsonFloat(math.Pi),
			gonaf.JsonFloat(math.NaN()),
			gonaf.JsonFloat(math.Inf(1)),
			gonaf.JsonFloat(math.Inf(-1)),
		},
	}
	bb := new(bytes.Buffer)
	encoder := json.NewEncoder(bb)

	if err := encoder.Encode(source); err != nil {
		t.Fatalf("Actual: %#v; Expected: %#v", err, nil)
	}

	decoder := json.NewDecoder(bb)
	var destination testType

	if err := decoder.Decode(&destination); err != nil {
		t.Fatalf("Actual: %#v; Expected: %#v", err, nil)
	}
	if actual, expected := len(destination.Values), len(source.Values); actual != expected {
		t.Fatalf("Actual: %#v; Expected: %#v", actual, expected)
	}
	if actual, expected := destination.Values[0], source.Values[0]; actual != expected {
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
