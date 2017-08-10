# gonaf

Go Not-A-Float is a library for encoding and decoding IEEE 754 values
to and from JSON.

*DEPRECATED* Please consider using
[https://github.com/karrick/goejs](https://github.com/karrick/goejs),
which serializes both strings and numbers in and out of JSON.

Documentation is available via
[![GoDoc](https://godoc.org/github.com/karrick/gonaf?status.svg)](https://godoc.org/github.com/karrick/gonaf).

## Description

The JSON encoding library in Go's standard library does things
correctly.  While the JSON data interchange format states that all
numbers are floating point values, it unfortunately does not support
three special cases defined by the IEEE 754 floating point standard:
Not-A-Number, also known as NaN, along with positive infinity, and
negative infinity.  These three special floating point values are
defined by the IEEE 754 standard but not defined by the JSON standard.

JSON is a very convenient and commonly used text based serialization
method.  There are JSON serialization libraries for most programming
languages that do support all legal floating point values, including
these special cases.  They do it by representing NaN as `null`,
`1e999` to represent +Infinity, and `-1e999` to represent -Infinity.
This Go library allows encoding and decoding floating point values
using these conventions.

## Usage Example

Normally one doesn't directly invoke a type's MarshalJSON and
UnmarshalJSON methods, but it's called by the standard library's
json.Marshal and json.Unmarshal methods when encoding or decoding a
value that is of the specified type.  Rather than demonstrate directly
invoking the methods required by the json.Marshaler and
json.Unmarshaler methods, a more typical use is demonstrated below.

```Go
    type MetricSeries struct {
        Times  []time.Time
        Values []gonaf.JsonFloat
    }


    func example() {
		baseTime := time.Now()
		source := MetricSeries{
		    Times:  []time.Time{
		        baseTime.Add(time.Minute),
		        baseTime.Add(2 * time.Minute),
		        baseTime.Add(3 * time.Minute),
		        baseTime.Add(4 * time.Minute),
		    },
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
		    return fmt.Errorf("cannot encode MetricSeries: %s", err)
		}
		
		decoder := json.NewDecoder(bb)
		var destination MetricSeries
		
		if err := decoder.Decode(&destination); err != nil {
		    return fmt.Errorf("cannot decode MetricSeries: %s", err)
		}
    }
```
