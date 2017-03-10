package gonaf

import (
	"bytes"
	"math"
	"strconv"
)

// JsonFloat is a float64 with JSON encoding and decoding support for NaN, +Infinity, and -Infinity.
//
// 	type MetricSeries struct {
// 	    Times  []time.Time
// 	    Values []gonaf.JsonFloat
// 	}
type JsonFloat float64

// MarshalJSON implements the json.Marshaler interface for the JsonFloat type, encoding NaN,
// +Infinity, and -Infinity using conventions found in many other JSON serialization libraries.
//
// Normally one doesn't directly invoke a type's MarshalJSON method, but it's called by the standard
// library's json.Marshal method when encoding a value that is of the specified type.  Rather than
// demonstrate directly invoking the MarshalJSON method, a more typical use is demonstrated below.
//
// 	baseTime := time.Now()
// 	source := MetricSeries{
// 		Times:  []time.Time{
// 			baseTime.Add(time.Minute),
// 			baseTime.Add(2 * time.Minute),
// 			baseTime.Add(3 * time.Minute),
// 			baseTime.Add(4 * time.Minute),
// 		},
// 		Values: []gonaf.JsonFloat{
// 			gonaf.JsonFloat(math.Pi),
// 			gonaf.JsonFloat(math.NaN()),
// 			gonaf.JsonFloat(math.Inf(1)),
// 			gonaf.JsonFloat(math.Inf(-1)),
// 		},
// 	}
// 	bb := new(bytes.Buffer)
// 	encoder := json.NewEncoder(bb)
//
// 	if err := encoder.Encode(source); err != nil {
// 		return fmt.Errorf("cannot encode MetricSeries: %s", err)
// 	}
func (jf JsonFloat) MarshalJSON() ([]byte, error) {
	if math.IsNaN(float64(jf)) {
		return []byte("null"), nil
	} else if math.IsInf(float64(jf), 1) {
		return []byte("1e999"), nil
	} else if math.IsInf(float64(jf), -1) {
		return []byte("-1e999"), nil
	}
	// NOTE: magic parameters in call to strconv.FormatFloat
	return []byte(strconv.FormatFloat(float64(jf), 'g', -1, 64)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for the JsonFloat type, decoding NaN,
// +Infinity, and -Infinity using conventions found in many other JSON serialization libraries.
//
// Normally one doesn't directly invoke a type's UnmarshalJSON method, but it's called by the
// standard library's json.Unmarshal method when decoding a value that is of the specified type.
// Rather than demonstrate directly invoking the UnmarshalJSON method, a more typical use is
// demonstrated below.
//
// 	decoder := json.NewDecoder(someReader)
// 	var destination MetricSeries
//
// 	if err := decoder.Decode(&destination); err != nil {
// 		return fmt.Errorf("cannot decode MetricSeries: %s", err)
// 	}
func (jf *JsonFloat) UnmarshalJSON(blob []byte) error {
	l := len(blob)
	if l >= 4 {
		if bytes.Equal(blob[:4], []byte("null")) {
			*jf = JsonFloat(math.NaN())
			return nil
		}
		if l >= 5 {
			if bytes.Equal(blob[:5], []byte("1e999")) {
				*jf = JsonFloat(math.Inf(1))
				return nil
			}
			if l >= 6 && bytes.Equal(blob[:6], []byte("-1e999")) {
				*jf = JsonFloat(math.Inf(-1))
				return nil
			}
		}
	}
	val, err := strconv.ParseFloat(string(blob), 64)
	// do not want to alter this variable when error
	if err != nil {
		return err
	}
	*jf = JsonFloat(val)
	return nil
}
