package collector

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// When unmarshalling JSON, the "xstring" type can be one of the following
// - nil
// - string
// - float64
// - integer
// - [{"Member": "VALUE"}]
type xstring string

func (w *xstring) UnmarshalJSON(data []byte) error {
	var x any
	var s string

	defer func() {
		*w = xstring(s)
	}()

	err := json.Unmarshal(data, &x)
	if err != nil {
		return err
	}

	// case: nil
	if x == nil {
		return nil
	}

	// case: string
	if v, ok := x.(string); ok {
		s = v
		return nil
	}

	// case: int
	if v, ok := x.(int); ok {
		s = fmt.Sprintf("%v", v)
		return nil
	}

	// case: float64
	if v, ok := x.(float64); ok {
		s = fmt.Sprintf("%v", v)
		return nil
	}

	// strange iDRAC case
	list, ok := x.([]any)
	if !ok {
		return nil
	}

	dict, ok := list[0].(map[string]any)
	if !ok {
		return nil
	}

	if v, ok := dict["Member"].(string); ok {
		s = v
		return nil
	}

	return nil
}

func (w *xstring) String() string {
	return string(*w)
}

// When unmarshalling JSON, the "xint" type can be a number or a string that
// holds one, as some BMCs send sizes such as StripSizeBytes as strings. Any
// other value reads as 0 instead of failing the whole response.
type xint int64

func (w *xint) UnmarshalJSON(data []byte) error {
	var x any

	*w = 0
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	switch v := x.(type) {
	case float64:
		*w = xint(v)
	case string:
		if n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil {
			*w = xint(n)
		}
	}

	return nil
}

func asPtr[T any](v T) *T {
	return &v
}

// asFloat64 coerces a value decoded into an "any" field (such as ReadingVolts,
// which some BMCs serialize as a JSON string) into a float64.
func asFloat64(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(x), 64)
		if err != nil {
			return 0, false
		}
		return f, true
	}
	return 0, false
}
