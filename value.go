package tree

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	BoolType = 0 + iota
	FloatType
	MapType
)

type Value struct {
	t int
	b *bool
	f *float64
	m map[string]*Value
}

func NewValue(v interface{}) (*Value, error) {
	value := &Value{}
	switch t := v.(type) {
	case bool:
		value.t = BoolType
		value.b = &t
	case float64:
		value.t = FloatType
		value.f = &t
	case map[string]interface{}:
		value.t = MapType
		value.m = map[string]*Value{}
		var err error
		for k, v := range t {
			value.m[k], err = NewValue(v)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("unsupported Value type")
	}

	return value, nil
}

func (v *Value) Bool() (bool, error) {
	if v.t != BoolType {
		return false, fmt.Errorf("value is not bool type")
	}

	return *v.b, nil
}

func (v *Value) Float() (float64, error) {
	if v.t != FloatType {
		return 0, fmt.Errorf("value is not float type")
	}

	return *v.f, nil
}

func (v *Value) Map() (map[string]*Value, error) {
	if v.t != MapType {
		return nil, fmt.Errorf("value is not map type")
	}

	return v.m, nil
}

func (v *Value) MarshalJSON() ([]byte, error) {
	b := []byte{}
	switch v.t {
	case BoolType:
		return strconv.AppendBool(b, *v.b), nil
	case FloatType:
		return strconv.AppendFloat(b, *v.f, 'f', -1, 64), nil
	case MapType:
		return json.Marshal(v.m)
	default:
		return nil, fmt.Errorf("unsupported Value type")
	}
}

func (v *Value) UnmarshalJSON(b []byte) error {
	s := string(b)
	var err error
	if *v.b, err = parseBool(b); err == nil {
		v.t = BoolType
		return nil
	}

	if *v.f, err = strconv.ParseFloat(s, 64); err == nil {
		v.t = FloatType
		return nil
	}

	m := map[string]interface{}{}
	if err = json.Unmarshal(b, &m); err == nil {
		val, err := NewValue(m)
		if err != nil {
			return err
		}
		v.t = MapType
		v.m = map[string]*Value{}
		for key, value := range val.m {
			v.m[key] = value
		}
		return nil
	}

	return fmt.Errorf("unable to unmarshal")
}

func parseBool(b []byte) (bool, error) {

	switch string(b) {
	case "t", "T", "true", "TRUE", "True":
		return true, nil
	case "f", "F", "false", "FALSE", "False":
		return false, nil
	}

	return false, fmt.Errorf("syntax error")
}
