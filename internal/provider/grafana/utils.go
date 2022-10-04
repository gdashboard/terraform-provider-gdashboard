package grafana

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
)

type BoolString struct {
	Flag  bool
	Value string
}

func (s *BoolString) UnmarshalJSON(raw []byte) error {
	if raw == nil || bytes.Equal(raw, []byte(`"null"`)) {
		return nil
	}
	var (
		tmp string
		err error
	)
	if raw[0] != '"' {
		if bytes.Equal(raw, []byte("true")) {
			s.Flag = true
			return nil
		}
		if bytes.Equal(raw, []byte("false")) {
			return nil
		}
		return errors.New("bad boolean value provided")
	}
	if err = json.Unmarshal(raw, &tmp); err != nil {
		return err
	}
	s.Value = tmp
	return nil
}

func (s BoolString) MarshalJSON() ([]byte, error) {
	if s.Value != "" {
		var buf bytes.Buffer
		buf.WriteRune('"')
		buf.WriteString(s.Value)
		buf.WriteRune('"')
		return buf.Bytes(), nil
	}
	return strconv.AppendBool([]byte{}, s.Flag), nil
}

type BoolInt struct {
	Flag  bool
	Value *int64
}

func (s *BoolInt) UnmarshalJSON(raw []byte) error {
	if raw == nil || bytes.Equal(raw, []byte(`"null"`)) {
		return nil
	}
	var (
		tmp int64
		err error
	)
	if tmp, err = strconv.ParseInt(string(raw), 10, 64); err != nil {
		if bytes.Equal(raw, []byte("true")) {
			s.Flag = true
			return nil
		}
		if bytes.Equal(raw, []byte("false")) {
			return nil
		}
		return errors.New("bad value provided")
	}
	s.Value = &tmp
	return nil
}

func (s BoolInt) MarshalJSON() ([]byte, error) {
	if s.Value != nil {
		return strconv.AppendInt([]byte{}, *s.Value, 10), nil
	}
	return strconv.AppendBool([]byte{}, s.Flag), nil
}

// StringSliceString represents special type for json values that could be
// strings or slice of strings: "something" or ["something"].
type StringSliceString struct {
	Value []string
	Valid bool
}

// UnmarshalJSON implements custom unmarshalling for StringSliceString type.
func (v *StringSliceString) UnmarshalJSON(raw []byte) error {
	if raw == nil || bytes.Equal(raw, []byte(`"null"`)) {
		return nil
	}

	// First try with string.
	var str string
	if err := json.Unmarshal(raw, &str); err == nil {
		v.Value = []string{str}
		v.Valid = true
		return nil
	}

	// Lastly try with string slice.
	var strSlice []string
	err := json.Unmarshal(raw, &strSlice)
	if err != nil {
		return err
	}

	v.Value = strSlice
	v.Valid = true
	return nil
}

// MarshalJSON implements custom marshalling for StringSliceString type.
func (v *StringSliceString) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte(`"null"`), nil
	}

	return json.Marshal(v.Value)
}
