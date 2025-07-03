// копия файла из https://github.com/jackc/pgtype/timestamp.go
// чтоб не выдавала ошибку на null
// чтобы дата NULL = time.Time{}
package postgres_pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"strings"
	"time"
	//"github.com/jackc/pgx/v5/internal/pgio"
)

const pgTimestampFormat = "2006-01-02 15:04:05.999999999"

//type TimestampScanner interface {
//	ScanTimestamp(v Timestamp) error
//}
//
//type TimestampValuer interface {
//	TimestampValue() (Timestamp, error)
//}

// Timestamp represents the PostgreSQL timestamp type.
type Timestamp struct {
	Time             time.Time // Time zone will be ignored when encoding to PostgreSQL.
	InfinityModifier pgtype.InfinityModifier
	Valid            bool
}

func (ts *Timestamp) ScanTimestamp(v Timestamp) error {
	*ts = v
	return nil
}

func (ts Timestamp) TimestampValue() (Timestamp, error) {
	return ts, nil
}

// Scan implements the database/sql Scanner interface.
func (ts *Timestamp) Scan(src any) error {
	if src == nil {
		*ts = Timestamp{Valid: true} //sanek
		//*ts = Timestamp{}
		return nil
	}

	switch src := src.(type) {
	case string:
		return (&scanPlanTextTimestampToTimestampScanner{}).Scan([]byte(src), ts)
	case time.Time:
		*ts = Timestamp{Time: src, Valid: true}
		return nil
	}

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (ts Timestamp) Value() (driver.Value, error) {
	if !ts.Valid {
		return nil, nil
	}

	if ts.InfinityModifier != pgtype.Finite {
		return ts.InfinityModifier.String(), nil
	}
	return ts.Time, nil
}

func (ts Timestamp) MarshalJSON() ([]byte, error) {
	if !ts.Valid {
		return []byte("null"), nil
	}

	var s string

	switch ts.InfinityModifier {
	case pgtype.Finite:
		s = ts.Time.Format(time.RFC3339Nano)
	case pgtype.Infinity:
		s = "infinity"
	case pgtype.NegativeInfinity:
		s = "-infinity"
	}

	return json.Marshal(s)
}

func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	var s *string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	if s == nil {
		*ts = Timestamp{}
		return nil
	}

	switch *s {
	case "infinity":
		*ts = Timestamp{Valid: true, InfinityModifier: pgtype.Infinity}
	case "-infinity":
		*ts = Timestamp{Valid: true, InfinityModifier: -pgtype.Infinity}
	default:
		// PostgreSQL uses ISO 8601 wihout timezone for to_json function and casting from a string to timestampt
		tim, err := time.Parse(time.RFC3339Nano, *s+"Z")
		if err != nil {
			return err
		}

		*ts = Timestamp{Time: tim, Valid: true}
	}

	return nil
}

type TimestampCodec struct {
	// ScanLocation is the location that the time is assumed to be in for scanning. This is different from
	// TimestamptzCodec.ScanLocation in that this setting does change the instant in time that the timestamp represents.
	ScanLocation *time.Location
}

func (*TimestampCodec) FormatSupported(format int16) bool {
	return format == pgtype.TextFormatCode || format == pgtype.BinaryFormatCode
}

func (*TimestampCodec) PreferredFormat() int16 {
	return pgtype.BinaryFormatCode
}

func (*TimestampCodec) PlanEncode(m *pgtype.Map, oid uint32, format int16, value any) pgtype.EncodePlan {
	if _, ok := value.(pgtype.TimestampValuer); !ok {
		return nil
	}

	switch format {
	case pgtype.BinaryFormatCode:
		return encodePlanTimestampCodecBinary{}
	case pgtype.TextFormatCode:
		return encodePlanTimestampCodecText{}
	}

	return nil
}

type encodePlanTimestampCodecBinary struct{}

func (encodePlanTimestampCodecBinary) Encode(value any, buf []byte) (newBuf []byte, err error) {
	ts, err := value.(pgtype.TimestampValuer).TimestampValue()
	if err != nil {
		return nil, err
	}

	if !ts.Valid {
		return nil, nil
	}

	var microsecSinceY2K int64
	switch ts.InfinityModifier {
	case pgtype.Finite:
		t := discardTimeZone(ts.Time)
		microsecSinceUnixEpoch := t.Unix()*1000000 + int64(t.Nanosecond())/1000
		microsecSinceY2K = microsecSinceUnixEpoch - microsecFromUnixEpochToY2K
	case pgtype.Infinity:
		microsecSinceY2K = infinityMicrosecondOffset
	case pgtype.NegativeInfinity:
		microsecSinceY2K = negativeInfinityMicrosecondOffset
	}

	buf = AppendInt64(buf, microsecSinceY2K)

	return buf, nil
}

type encodePlanTimestampCodecText struct{}

func (encodePlanTimestampCodecText) Encode(value any, buf []byte) (newBuf []byte, err error) {
	ts, err := value.(pgtype.TimestampValuer).TimestampValue()
	if err != nil {
		return nil, err
	}

	if !ts.Valid {
		return nil, nil
	}

	var s string

	switch ts.InfinityModifier {
	case pgtype.Finite:
		t := discardTimeZone(ts.Time)

		// Year 0000 is 1 BC
		bc := false
		if year := t.Year(); year <= 0 {
			year = -year + 1
			t = time.Date(year, t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
			bc = true
		}

		s = t.Truncate(time.Microsecond).Format(pgTimestampFormat)

		if bc {
			s = s + " BC"
		}
	case pgtype.Infinity:
		s = "infinity"
	case pgtype.NegativeInfinity:
		s = "-infinity"
	}

	buf = append(buf, s...)

	return buf, nil
}

func discardTimeZone(t time.Time) time.Time {
	if t.Location() != time.UTC {
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
	}

	return t
}

func (c *TimestampCodec) PlanScan(m *pgtype.Map, oid uint32, format int16, target any) pgtype.ScanPlan {
	switch format {
	case pgtype.BinaryFormatCode:
		name := getInterfaceName(target) //sanek
		switch name {
		case "*pgtype.timeWrapper":
			return &scanPlanBinaryTimestampToTimestampScanner{location: c.ScanLocation}
		}
	case pgtype.TextFormatCode:
		name := getInterfaceName(target) //sanek
		switch name {
		case "*pgtype.timeWrapper":
			return &scanPlanTextTimestampToTimestampScanner{location: c.ScanLocation}
		}
	}

	return nil
}

type scanPlanBinaryTimestampToTimestampScanner struct{ location *time.Location }

func (plan *scanPlanBinaryTimestampToTimestampScanner) Scan(src []byte, dst any) error {
	scanner := (dst).(pgtype.TimestampScanner)

	if src == nil {
		return scanner.ScanTimestamp(pgtype.Timestamp{Valid: true})
		//return scanner.ScanTimestamp(pgtype.Timestamp{})
	}

	if len(src) != 8 {
		return fmt.Errorf("invalid length for timestamp: %v", len(src))
	}

	var ts pgtype.Timestamp
	microsecSinceY2K := int64(binary.BigEndian.Uint64(src))

	switch microsecSinceY2K {
	case infinityMicrosecondOffset:
		ts = pgtype.Timestamp{Valid: true, InfinityModifier: pgtype.Infinity}
	case negativeInfinityMicrosecondOffset:
		ts = pgtype.Timestamp{Valid: true, InfinityModifier: -pgtype.Infinity}
	default:
		tim := time.Unix(
			microsecFromUnixEpochToY2K/1000000+microsecSinceY2K/1000000,
			(microsecFromUnixEpochToY2K%1000000*1000)+(microsecSinceY2K%1000000*1000),
		).UTC()
		if plan.location != nil {
			tim = time.Date(tim.Year(), tim.Month(), tim.Day(), tim.Hour(), tim.Minute(), tim.Second(), tim.Nanosecond(), plan.location)
		}
		ts = pgtype.Timestamp{Time: tim, Valid: true}
	}

	return scanner.ScanTimestamp(ts)
}

type scanPlanTextTimestampToTimestampScanner struct{ location *time.Location }

func (plan *scanPlanTextTimestampToTimestampScanner) Scan(src []byte, dst any) error {
	scanner := (dst).(pgtype.TimestampScanner)

	if src == nil {
		return scanner.ScanTimestamp(pgtype.Timestamp{})
	}

	var ts pgtype.Timestamp
	sbuf := string(src)
	switch sbuf {
	case "infinity":
		ts = pgtype.Timestamp{Valid: true, InfinityModifier: pgtype.Infinity}
	case "-infinity":
		ts = pgtype.Timestamp{Valid: true, InfinityModifier: -pgtype.Infinity}
	default:
		bc := false
		if strings.HasSuffix(sbuf, " BC") {
			sbuf = sbuf[:len(sbuf)-3]
			bc = true
		}
		tim, err := time.Parse(pgTimestampFormat, sbuf)
		if err != nil {
			return err
		}

		if bc {
			year := -tim.Year() + 1
			tim = time.Date(year, tim.Month(), tim.Day(), tim.Hour(), tim.Minute(), tim.Second(), tim.Nanosecond(), tim.Location())
		}

		if plan.location != nil {
			tim = time.Date(tim.Year(), tim.Month(), tim.Day(), tim.Hour(), tim.Minute(), tim.Second(), tim.Nanosecond(), plan.location)
		}

		ts = pgtype.Timestamp{Time: tim, Valid: true}
	}

	return scanner.ScanTimestamp(ts)
}

func (c *TimestampCodec) DecodeDatabaseSQLValue(m *pgtype.Map, oid uint32, format int16, src []byte) (driver.Value, error) {
	if src == nil {
		return nil, nil
	}

	var ts Timestamp
	err := codecScan(c, m, oid, format, src, &ts)
	if err != nil {
		return nil, err
	}

	if ts.InfinityModifier != pgtype.Finite {
		return ts.InfinityModifier.String(), nil
	}

	return ts.Time, nil
}

func (c *TimestampCodec) DecodeValue(m *pgtype.Map, oid uint32, format int16, src []byte) (any, error) {
	if src == nil {
		return nil, nil
	}

	var ts Timestamp
	err := codecScan(c, m, oid, format, src, &ts)
	if err != nil {
		return nil, err
	}

	if ts.InfinityModifier != pgtype.Finite {
		return ts.InfinityModifier, nil
	}

	return ts.Time, nil
}
