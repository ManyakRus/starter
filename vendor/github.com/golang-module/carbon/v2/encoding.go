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

// Time defines a Time struct.
// 定义 Time 结构体
type Time struct {
	Carbon
}

// TimeMilli defines a TimeMilli struct.
// 定义 TimeMilli 结构体
type TimeMilli struct {
	Carbon
}

// TimeMicro defines a TimeMicro struct.
// 定义 TimeMicro 结构体
type TimeMicro struct {
	Carbon
}

// TimeNano defines a TimeNano struct.
// 定义 TimeNano 结构体
type TimeNano struct {
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
	c := ParseByLayout(string(bytes.Trim(b, `"`)), DateTimeLayout, t.Location())
	if c.Error == nil {
		*t = DateTime{Carbon: c}
	}
	return c.Error
}

// MarshalJSON implements the interface json.Marshal for DateTimeMilli struct.
// 实现 MarshalJSON 接口
func (t DateTimeMilli) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToDateTimeMilliString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for DateTimeMilli struct.
// 实现 UnmarshalJSON 接口
func (t *DateTimeMilli) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), DateTimeMilliLayout, t.Location())
	if c.Error == nil {
		*t = DateTimeMilli{Carbon: c}
	}
	return c.Error
}

// MarshalJSON implements the interface json.Marshal for DateTimeMicro struct.
// 实现 MarshalJSON 接口
func (t DateTimeMicro) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToDateTimeMicroString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for DateTimeMicro struct.
// 实现 UnmarshalJSON 接口
func (t *DateTimeMicro) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), DateTimeMicroLayout, t.Location())
	if c.Error == nil {
		*t = DateTimeMicro{Carbon: c}
	}
	return c.Error
}

// MarshalJSON implements the interface json.Marshal for DateTimeNano struct.
// 实现 MarshalJSON 接口
func (t DateTimeNano) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToDateTimeNanoString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for DateTimeNano struct.
// 实现 UnmarshalJSON 接口
func (t *DateTimeNano) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), DateTimeNanoLayout, t.Location())
	if c.Error == nil {
		*t = DateTimeNano{Carbon: c}
	}
	return c.Error
}

// MarshalJSON implements the interface json.Marshal for Date struct.
// 实现 MarshalJSON 接口
func (t Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToDateString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for Date struct.
// 实现 UnmarshalJSON 接口
func (t *Date) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), DateLayout, t.Location())
	if c.Error == nil {
		*t = Date{Carbon: c}
	}
	return c.Error
}

// MarshalJSON implements the interface json.Marshal for DateMilli struct.
// 实现 MarshalJSON 接口
func (t DateMilli) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToDateMilliString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for DateMilli struct.
// 实现 UnmarshalJSON 接口
func (t *DateMilli) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), DateMilliLayout, t.Location())
	if c.Error == nil {
		*t = DateMilli{Carbon: c}
	}
	return c.Error
}

// MarshalJSON implements the interface json.Marshal for DateMicro struct.
// 实现 MarshalJSON 接口
func (t DateMicro) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToDateMicroString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for DateMicro struct.
// 实现 UnmarshalJSON 接口
func (t *DateMicro) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), DateMicroLayout, t.Location())
	if c.Error == nil {
		*t = DateMicro{Carbon: c}
	}
	return c.Error
}

// MarshalJSON implements the interface json.Marshal for DateNano struct.
// 实现 MarshalJSON 接口
func (t DateNano) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToDateNanoString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for DateNano struct.
// 实现 UnmarshalJSON 接口
func (t *DateNano) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), DateNanoLayout, t.Location())
	if c.Error == nil {
		*t = DateNano{Carbon: c}
	}
	return c.Error
}

// MarshalJSON implements the interface json.Marshal for Time struct.
// 实现 MarshalJSON 接口
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToTimeString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for Time struct.
// 实现 UnmarshalJSON 接口
func (t *Time) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), TimeLayout, t.Location())
	if c.Error == nil {
		*t = Time{Carbon: c}
	}
	return c.Error
}

// MarshalJSON implements the interface json.Marshal for TimeMilli struct.
// 实现 MarshalJSON 接口
func (t TimeMilli) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToTimeMilliString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for TimeMilli struct.
// 实现 UnmarshalJSON 接口
func (t *TimeMilli) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), TimeMilliLayout, t.Location())
	if c.Error == nil {
		*t = TimeMilli{Carbon: c}
	}
	return c.Error
}

// MarshalJSON implements the interface json.Marshal for TimeMicro struct.
// 实现 MarshalJSON 接口
func (t TimeMicro) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToTimeMicroString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for TimeMicro struct.
// 实现 UnmarshalJSON 接口
func (t *TimeMicro) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), TimeMicroLayout, t.Location())
	if c.Error == nil {
		*t = TimeMicro{Carbon: c}
	}
	return c.Error
}

// MarshalJSON implements the interface json.Marshal for TimeNano struct.
// 实现 MarshalJSON 接口
func (t TimeNano) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.ToTimeNanoString())), nil
}

// UnmarshalJSON implements the interface json.Unmarshal for TimeNano struct.
// 实现 UnmarshalJSON 接口
func (t *TimeNano) UnmarshalJSON(b []byte) error {
	c := ParseByLayout(string(bytes.Trim(b, `"`)), TimeNanoLayout, t.Location())
	if c.Error == nil {
		*t = TimeNano{Carbon: c}
	}
	return c.Error
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
		*t = Timestamp{Carbon: c}
	}
	return c.Error
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
		*t = TimestampMilli{Carbon: c}
	}
	return c.Error
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
		*t = TimestampMicro{Carbon: c}
	}
	return c.Error
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
		*t = TimestampNano{Carbon: c}
	}
	return c.Error
}

// MarshalJSON implements the json.Marshaler interface.
// 实现 json.Marshaler 接口
func (c Carbon) MarshalJSON() ([]byte, error) {
	if c.Error != nil {
		return nil, c.Error
	}
	key, value, tz := c.parseTag()
	data := ""
	if key == "layout" {
		data = fmt.Sprintf(`"%s"`, c.Layout(value, tz))
	}
	if key == "format" {
		// timestamp without double quotes in json
		if value == "U" || value == "V" || value == "X" || value == "Z" {
			data = fmt.Sprintf(`%s`, c.Format(value, tz))
		} else {
			data = fmt.Sprintf(`"%s"`, c.Format(value, tz))
		}
	}
	return []byte(data), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// 实现 json.Unmarshaler 接口
func (c *Carbon) UnmarshalJSON(b []byte) error {
	if c.Error != nil {
		return c.Error
	}
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	key, value, tz := c.parseTag()
	data := fmt.Sprintf("%s", bytes.Trim(b, `"`))
	if key == "layout" {
		*c = ParseByLayout(data, value, tz)
	}
	if key == "format" {
		*c = ParseByFormat(data, value, tz)
	}
	c.tag = &tag{
		carbon: fmt.Sprintf("%s:%s", key, value),
		tz:     tz,
	}
	return c.Error
}

// Int64 outputs timestamp with second.
// 输出秒级时间戳
func (t Timestamp) Int64() int64 {
	return t.Timestamp()
}

// Int64 outputs timestamp with millisecond.
// 输出豪秒级时间戳
func (t TimestampMilli) Int64() int64 {
	return t.TimestampMilli()
}

// Int64 outputs timestamp with microsecond.
// 输出微秒级时间戳
func (t TimestampMicro) Int64() int64 {
	return t.TimestampMicro()
}

// Int64 outputs timestamp with nanosecond.
// 输出纳秒级时间戳
func (t TimestampNano) Int64() int64 {
	return t.TimestampNano()
}

// String implements the interface Stringer for DateTime struct.
// 实现 Stringer 接口
func (t DateTime) String() string {
	return t.ToDateTimeString()
}

// String implements the interface Stringer for DateTimeMilli struct.
// 实现 Stringer 接口
func (t DateTimeMilli) String() string {
	return t.ToDateTimeMilliString()
}

// String implements the interface Stringer for DateTimeMicro struct.
// 实现 Stringer 接口
func (t DateTimeMicro) String() string {
	return t.ToDateTimeMicroString()
}

// String implements the interface Stringer for DateTimeNano struct.
// 实现 Stringer 接口
func (t DateTimeNano) String() string {
	return t.ToDateTimeNanoString()
}

// String implements the interface Stringer for Date struct.
// 实现 Stringer 接口
func (t Date) String() string {
	return t.ToDateString()
}

// String implements the interface Stringer for DateMilli struct.
// 实现 Stringer 接口
func (t DateMilli) String() string {
	return t.ToDateMilliString()
}

// String implements the interface Stringer for DateMicro struct.
// 实现 Stringer 接口
func (t DateMicro) String() string {
	return t.ToDateMicroString()
}

// String implements the interface Stringer for DateNano struct.
// 实现 Stringer 接口
func (t DateNano) String() string {
	return t.ToDateNanoString()
}

// String implements the interface Stringer for Time struct.
// 实现 Stringer 接口
func (t Time) String() string {
	return t.ToTimeString()
}

// String implements the interface Stringer for TimeMilli struct.
// 实现 Stringer 接口
func (t TimeMilli) String() string {
	return t.ToTimeMilliString()
}

// String implements the interface Stringer for TimeMicro struct.
// 实现 Stringer 接口
func (t TimeMicro) String() string {
	return t.ToTimeMicroString()
}

// String implements the interface Stringer for TimeNano struct.
// 实现 Stringer 接口
func (t TimeNano) String() string {
	return t.ToTimeNanoString()
}

// String implements the interface Stringer for Timestamp struct.
// 实现 Stringer 接口
func (t Timestamp) String() string {
	return strconv.FormatInt(t.Timestamp(), 10)
}

// String implements the interface Stringer for TimestampMilli struct.
// 实现 Stringer 接口
func (t TimestampMilli) String() string {
	return strconv.FormatInt(t.TimestampMilli(), 10)
}

// String implements the interface Stringer for TimestampMicro struct.
// 实现 Stringer 接口
func (t TimestampMicro) String() string {
	return strconv.FormatInt(t.TimestampMicro(), 10)
}

// String implements the interface Stringer for TimestampNano struct.
// 实现 Stringer 接口
func (t TimestampNano) String() string {
	return strconv.FormatInt(t.TimestampNano(), 10)
}

// GormDataType implements the interface GormDataTypeInterface for DateTime struct.
// 实现 GormDataTypeInterface 接口
func (t DateTime) GormDataType() string {
	return "time"
}

// GormDataType implements the interface GormDataTypeInterface for DateTimeMilli struct.
// 实现 GormDataTypeInterface 接口
func (t DateTimeMilli) GormDataType() string {
	return "time"
}

// GormDataType implements the interface GormDataTypeInterface for DateTimeMicro struct.
// 实现 GormDataTypeInterface 接口
func (t DateTimeMicro) GormDataType() string {
	return "time"
}

// GormDataType implements the interface GormDataTypeInterface for DateTimeNano struct.
// 实现 GormDataTypeInterface 接口
func (t DateTimeNano) GormDataType() string {
	return "time"
}

// GormDataType implements the interface GormDataTypeInterface for Date struct.
// 实现 GormDataTypeInterface 接口
func (t Date) GormDataType() string {
	return "time"
}

// GormDataType implements the interface GormDataTypeInterface for DateMilli struct.
// 实现 GormDataTypeInterface 接口
func (t DateMilli) GormDataType() string {
	return "time"
}

// GormDataType implements the interface GormDataTypeInterface for DateMicro struct.
// 实现 GormDataTypeInterface 接口
func (t DateMicro) GormDataType() string {
	return "time"
}

// GormDataType implements the interface GormDataTypeInterface for DateNano struct.
// 实现 GormDataTypeInterface 接口
func (t DateNano) GormDataType() string {
	return "time"
}

// GormDataType implements the interface GormDataTypeInterface for Time struct.
// 实现 GormDataTypeInterface 接口
func (t Time) GormDataType() string {
	return "time"
}

// GormDataType implements the interface GormDataTypeInterface for TimeMilli struct.
// 实现 GormDataTypeInterface 接口
func (t TimeMilli) GormDataType() string {
	return "time"
}

// GormDataType implements the interface GormDataTypeInterface for TimeMicro struct.
// 实现 GormDataTypeInterface 接口
func (t TimeMicro) GormDataType() string {
	return "time"
}

// GormDataType implements the interface GormDataTypeInterface for TimeNano struct.
// 实现 GormDataTypeInterface 接口
func (t TimeNano) GormDataType() string {
	return "time"
}

// GormDataType implements the interface GormDataTypeInterface for Timestamp struct.
// 实现 GormDataTypeInterface 接口
func (t Timestamp) GormDataType() string {
	return "int"
}

// GormDataType implements the interface GormDataTypeInterface for TimestampMilli struct.
// 实现 GormDataTypeInterface 接口
func (t TimestampMilli) GormDataType() string {
	return "int"
}

// GormDataType implements the interface GormDataTypeInterface for TimestampMicro struct.
// 实现 GormDataTypeInterface 接口
func (t TimestampMicro) GormDataType() string {
	return "int"
}

// GormDataType implements the interface GormDataTypeInterface for TimestampNano struct.
// 实现 GormDataTypeInterface 接口
func (t TimestampNano) GormDataType() string {
	return "int"
}

// ToDateTimeStruct converts Carbon to DateTime.
// 将 Carbon 结构体转换成 DateTime 结构体
func (c Carbon) ToDateTimeStruct() DateTime {
	return DateTime{Carbon: c}
}

// ToDateTimeMilliStruct converts Carbon to DateTimeMilli.
// 将 Carbon 结构体转换成 DateTimeMilli 结构体
func (c Carbon) ToDateTimeMilliStruct() DateTimeMilli {
	return DateTimeMilli{Carbon: c}
}

// ToDateTimeMicroStruct converts Carbon to DateTimeMicro.
// 将 Carbon 结构体转换成 DateTimeMicro 结构体
func (c Carbon) ToDateTimeMicroStruct() DateTimeMicro {
	return DateTimeMicro{Carbon: c}
}

// ToDateTimeNanoStruct converts Carbon to DateTimeNano.
// 将 Carbon 结构体转换成 DateTimeNano 结构体
func (c Carbon) ToDateTimeNanoStruct() DateTimeNano {
	return DateTimeNano{Carbon: c}
}

// ToDateStruct converts Carbon to Date.
// 将 Carbon 结构体转换成 Date 结构体
func (c Carbon) ToDateStruct() Date {
	return Date{Carbon: c}
}

// ToDateMilliStruct converts Carbon to DateMilli.
// 将 Carbon 结构体转换成 DateMilli 结构体
func (c Carbon) ToDateMilliStruct() DateMilli {
	return DateMilli{Carbon: c}
}

// ToDateMicroStruct converts Carbon to DateMicro.
// 将 Carbon 结构体转换成 DateMicro 结构体
func (c Carbon) ToDateMicroStruct() DateMicro {
	return DateMicro{Carbon: c}
}

// ToDateNanoStruct converts Carbon to DateNano.
// 将 Carbon 结构体转换成 DateNano 结构体
func (c Carbon) ToDateNanoStruct() DateNano {
	return DateNano{Carbon: c}
}

// ToTimeStruct converts Carbon to Time.
// 将 Carbon 结构体转换成 Time 结构体
func (c Carbon) ToTimeStruct() Time {
	return Time{Carbon: c}
}

// ToTimeMilliStruct converts Carbon to TimeMilli.
// 将 Carbon 结构体转换成 TimeMilli 结构体
func (c Carbon) ToTimeMilliStruct() TimeMilli {
	return TimeMilli{Carbon: c}
}

// ToTimeMicroStruct converts Carbon to TimeMicro.
// 将 Carbon 结构体转换成 TimeMicro 结构体
func (c Carbon) ToTimeMicroStruct() TimeMicro {
	return TimeMicro{Carbon: c}
}

// ToTimeNanoStruct converts Carbon to TimeNano.
// 将 Carbon 结构体转换成 TimeNano 结构体
func (c Carbon) ToTimeNanoStruct() TimeNano {
	return TimeNano{Carbon: c}
}

// ToTimestampStruct converts Carbon to Timestamp.
// 将 Carbon 结构体转换成 Timestamp 结构体
func (c Carbon) ToTimestampStruct() Timestamp {
	return Timestamp{Carbon: c}
}

// ToTimestampMilliStruct converts Carbon to TimestampMilli.
// 将 Carbon 结构体转换成 TimestampMilli 结构体
func (c Carbon) ToTimestampMilliStruct() TimestampMilli {
	return TimestampMilli{Carbon: c}
}

// ToTimestampMicroStruct converts Carbon to TimestampMicro.
// 将 Carbon 结构体转换成 TimestampMicro 结构体
func (c Carbon) ToTimestampMicroStruct() TimestampMicro {
	return TimestampMicro{Carbon: c}
}

// ToTimestampNanoStruct converts Carbon to TimestampNano.
// 将 Carbon 结构体转换成 TimestampNano 结构体
func (c Carbon) ToTimestampNanoStruct() TimestampNano {
	return TimestampNano{Carbon: c}
}
