package carbon

import (
	"bytes"
	"fmt"
	"strconv"
)

// DateTime defines a DateTime struct.
// 定义 DateTime 结构体
type DateTime struct {
	Carbon
}

// DateTimeMilli defines a DateTimeMilli struct.
// 定义 DateTimeMilli 结构体
type DateTimeMilli struct {
	Carbon
}

// DateTimeMicro defines a DateTimeMicro struct.
// 定义 DateTimeMicro 结构体
type DateTimeMicro struct {
	Carbon
}

// DateTimeNano defines a DateTimeNano struct.
// 定义 DateTimeNano 结构体
type DateTimeNano struct {
	Carbon
}

// Date defines a Date struct.
// 定义 Date 结构体
type Date struct {
	Carbon
}

// DateMilli defines a DateMilli struct.
// 定义 DateMilli 结构体
type DateMilli struct {
	Carbon
}

// DateMicro defines a DateMicro struct.
// 定义 DateMicro 结构体
type DateMicro struct {
	Carbon
}

// DateNano defines a DateNano struct.
// 定义 DateNano 结构体
type DateNano struct {
	Carbon
}

// Timestamp defines a Timestamp struct.
// 定义 Timestamp 结构体
type Timestamp struct {
	Carbon
}

// TimestampMilli defines a TimestampMilli struct.
// 定义 TimestampMilli 结构体
type TimestampMilli struct {
	Carbon
}

// TimestampMicro defines a TimestampMicro struct.
// 定义 TimestampMicro 结构体
type TimestampMicro struct {
	Carbon
}

// TimestampNano defines a TimestampNano struct.
// 定义 TimestampNano 结构体
type TimestampNano struct {
	Carbon
}

// MarshalJSON implements the interface json.Marshal for DateTime struct.
// 实现 MarshalJSON 接口
func (t DateTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToDateTimeString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for DateTime struct.
// 实现 UnmarshalJSON 接口
func (t *DateTime) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), DateTimeLayout)
	if c.Error == nil {
		*t = DateTime{c}
	}
	return c.Error
}

// String implements the interface Stringer for DateTime struct.
// 实现 Stringer 接口
func (t DateTime) String() string {
	return t.ToDateTimeString()
}

// MarshalJSON implements the interface json.Marshal for DateTimeMilli struct.
// 实现 MarshalJSON 接口
func (t DateTimeMilli) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToDateTimeMilliString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for DateTimeMilli struct.
// 实现 UnmarshalJSON 接口
func (t *DateTimeMilli) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), DateTimeMilliLayout)
	if c.Error == nil {
		*t = DateTimeMilli{c}
	}
	return c.Error
}

// String implements the interface Stringer for DateTimeMilli struct.
// 实现 Stringer 接口
func (t DateTimeMilli) String() string {
	return t.ToDateTimeMilliString()
}

// MarshalJSON implements the interface json.Marshal for DateTimeMicro struct.
// 实现 MarshalJSON 接口
func (t DateTimeMicro) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToDateTimeMicroString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for DateTimeMicro struct.
// 实现 UnmarshalJSON 接口
func (t *DateTimeMicro) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), DateTimeMicroLayout)
	if c.Error == nil {
		*t = DateTimeMicro{c}
	}
	return c.Error
}

// String implements the interface Stringer for DateTimeMicro struct.
// 实现 Stringer 接口
func (t DateTimeMicro) String() string {
	return t.ToDateTimeMicroString()
}

// MarshalJSON implements the interface json.Marshal for DateTimeNano struct.
// 实现 MarshalJSON 接口
func (t DateTimeNano) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToDateTimeNanoString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for DateTimeNano struct.
// 实现 UnmarshalJSON 接口
func (t *DateTimeNano) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), DateTimeNanoLayout)
	if c.Error == nil {
		*t = DateTimeNano{c}
	}
	return c.Error
}

// String implements the interface Stringer for DateTimeNano struct.
// 实现 Stringer 接口
func (t DateTimeNano) String() string {
	return t.ToDateTimeNanoString()
}

// MarshalJSON implements the interface json.Marshal for Date struct.
// 实现 MarshalJSON 接口
func (t Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToDateString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for Date struct.
// 实现 UnmarshalJSON 接口
func (t *Date) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), DateLayout)
	if c.Error == nil {
		*t = Date{c}
	}
	return c.Error
}

// String implements the interface Stringer for Date struct.
// 实现 Stringer 接口
func (t Date) String() string {
	return t.ToDateString()
}

// MarshalJSON implements the interface json.Marshal for DateMilli struct.
// 实现 MarshalJSON 接口
func (t DateMilli) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToDateMilliString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for DateMilli struct.
// 实现 UnmarshalJSON 接口
func (t *DateMilli) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), DateMilliLayout)
	if c.Error == nil {
		*t = DateMilli{c}
	}
	return c.Error
}

// String implements the interface Stringer for DateMilli struct.
// 实现 Stringer 接口
func (t DateMilli) String() string {
	return t.ToDateMilliString()
}

// MarshalJSON implements the interface json.Marshal for DateMicro struct.
// 实现 MarshalJSON 接口
func (t DateMicro) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToDateMicroString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for DateMicro struct.
// 实现 UnmarshalJSON 接口
func (t *DateMicro) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), DateMicroLayout)
	if c.Error == nil {
		*t = DateMicro{c}
	}
	return c.Error
}

// String implements the interface Stringer for DateMicro struct.
// 实现 Stringer 接口
func (t DateMicro) String() string {
	return t.ToDateMicroString()
}

// MarshalJSON implements the interface json.Marshal for DateNano struct.
// 实现 MarshalJSON 接口
func (t DateNano) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToDateNanoString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for DateNano struct.
// 实现 UnmarshalJSON 接口
func (t *DateNano) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), DateNanoLayout)
	if c.Error == nil {
		*t = DateNano{c}
	}
	return c.Error
}

// String implements the interface Stringer for DateNano struct.
// 实现 Stringer 接口
func (t DateNano) String() string {
	return t.ToDateNanoString()
}

// MarshalJSON implements the interface json.Marshal for Timestamp struct.
// 实现 MarshalJSON 接口
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`%d`, t.Timestamp())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for Timestamp struct.
// 实现 UnmarshalJSON 接口
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, _ := strconv.ParseInt(string(b), 10, 64)
	c := CreateFromTimestamp(ts)
	if c.Error == nil {
		*t = Timestamp{c}
	}
	return c.Error
}

// String implements the interface Stringer for Timestamp struct.
// 实现 Stringer 接口
func (t Timestamp) String() string {
	return t.ToDateTimeString()
}

// MarshalJSON implements the interface json.Marshal for TimestampMilli struct.
// 实现 MarshalJSON 接口
func (t TimestampMilli) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`%d`, t.TimestampMilli())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for TimestampMilli struct.
// 实现 UnmarshalJSON 接口
func (t *TimestampMilli) UnmarshalJSON(b []byte) error {
	ts, _ := strconv.ParseInt(string(b), 10, 64)
	c := CreateFromTimestampMilli(ts)
	if c.Error == nil {
		*t = TimestampMilli{c}
	}
	return c.Error
}

// String implements the interface Stringer for TimestampMilli struct.
// 实现 Stringer 接口
func (t TimestampMilli) String() string {
	return t.ToDateTimeMilliString()
}

// MarshalJSON implements the interface MarshalJSON for TimestampMicro struct.
// 实现 MarshalJSON 接口
func (t TimestampMicro) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`%d`, t.TimestampMicro())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for TimestampMicro struct.
// 实现 UnmarshalJSON 接口
func (t *TimestampMicro) UnmarshalJSON(b []byte) error {
	ts, _ := strconv.ParseInt(string(b), 10, 64)
	c := CreateFromTimestampMicro(ts)
	if c.Error == nil {
		*t = TimestampMicro{c}
	}
	return c.Error
}

// String implements the interface Stringer for TimestampMicro struct.
// 实现 Stringer 接口
func (t TimestampMicro) String() string {
	return t.ToDateTimeMicroString()
}

// MarshalJSON implements the interface json.Marshal for TimestampNano struct.
// 实现 MarshalJSON 接口
func (t TimestampNano) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`%d`, t.TimestampNano())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for TimestampNano struct.
// 实现 UnmarshalJSON 接口
func (t *TimestampNano) UnmarshalJSON(b []byte) error {
	ts, _ := strconv.ParseInt(string(b), 10, 64)
	c := CreateFromTimestampNano(ts)
	if c.Error == nil {
		*t = TimestampNano{c}
	}
	return c.Error
}

// String implements the interface Stringer for TimestampNano struct.
// 实现 Stringer 接口
func (t TimestampNano) String() string {
	return t.ToDateTimeNanoString()
}
