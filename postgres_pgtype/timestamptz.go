// копия файла из https://github.com/jackc/pgtype/timestamptz.go
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
)

const pgTimestamptzHourFormat = "2006-01-02 15:04:05.999999999Z07"
const pgTimestamptzMinuteFormat = "2006-01-02 15:04:05.999999999Z07:00"
const pgTimestamptzSecondFormat = "2006-01-02 15:04:05.999999999Z07:00:00"
const microsecFromUnixEpochToY2K = 946684800 * 1000000

const (
	negativeInfinityMicrosecondOffset = -9223372036854775808
	infinityMicrosecondOffset         = 9223372036854775807
)

//type TimestamptzScanner interface {
//	ScanTimestamptz(v Timestamptz) error
//}
//
//type TimestamptzValuer interface {
//	TimestamptzValue() (Timestamptz, error)
//}

// Timestamptz represents the PostgreSQL timestamptz type.
type Timestamptz struct {
	Time             time.Time
	InfinityModifier pgtype.InfinityModifier
	Valid            bool
}

func (tstz *Timestamptz) ScanTimestamptz(v Timestamptz) error {
	*tstz = v
	return nil
}

func (tstz Timestamptz) TimestamptzValue() (Timestamptz, error) {
	return tstz, nil
}

// Scan implements the database/sql Scanner interface.
func (tstz *Timestamptz) Scan(src any) error {
	if src == nil {
		*tstz = Timestamptz{Valid: true} //sanek
		//*tstz = Timestamptz{}
		return nil
	}

	switch src := src.(type) {
	case string:
		return (&scanPlanTextTimestamptzToTimestamptzScanner{}).Scan([]byte(src), tstz)
	case time.Time:
		*tstz = Timestamptz{Time: src, Valid: true}
		return nil
	}

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (tstz Timestamptz) Value() (driver.Value, error) {
	if !tstz.Valid {
		return nil, nil
	}

	if tstz.InfinityModifier != pgtype.Finite {
		return tstz.InfinityModifier.String(), nil
	}
	return tstz.Time, nil
}

func (tstz Timestamptz) MarshalJSON() ([]byte, error) {
	if !tstz.Valid {
		return []byte("null"), nil
	}

	var s string

	switch tstz.InfinityModifier {
	case pgtype.Finite:
		s = tstz.Time.Format(time.RFC3339Nano)
	case pgtype.Infinity:
		s = "infinity"
	case pgtype.NegativeInfinity:
		s = "-infinity"
	}

	return json.Marshal(s)
}

func (tstz *Timestamptz) UnmarshalJSON(b []byte) error {
	var s *string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	if s == nil {
		*tstz = Timestamptz{}
		return nil
	}

	switch *s {
	case "infinity":
		*tstz = Timestamptz{Valid: true, InfinityModifier: pgtype.Infinity}
	case "-infinity":
		*tstz = Timestamptz{Valid: true, InfinityModifier: -pgtype.Infinity}
	default:
		// PostgreSQL uses ISO 8601 for to_json function and casting from a string to timestamptz
		tim, err := time.Parse(time.RFC3339Nano, *s)
		if err != nil {
			return err
		}

		*tstz = Timestamptz{Time: tim, Valid: true}
	}

	return nil
}

type TimestamptzCodec struct {
	// ScanLocation is the location to return scanned timestamptz values in. This does not change the instant in time that
	// the timestamptz represents.
	ScanLocation *time.Location
}

func (*TimestamptzCodec) FormatSupported(format int16) bool {
	return format == pgtype.TextFormatCode || format == pgtype.BinaryFormatCode
}

func (*TimestamptzCodec) PreferredFormat() int16 {
	return pgtype.BinaryFormatCode
}

func (*TimestamptzCodec) PlanEncode(m *pgtype.Map, oid uint32, format int16, value any) pgtype.EncodePlan {
	if _, ok := value.(pgtype.TimestamptzValuer); !ok {
		return nil
	}

	switch format {
	case pgtype.BinaryFormatCode:
		return encodePlanTimestamptzCodecBinary{}
	case pgtype.TextFormatCode:
		return encodePlanTimestamptzCodecText{}
	}

	return nil
}

type encodePlanTimestamptzCodecBinary struct{}

func (encodePlanTimestamptzCodecBinary) Encode(value any, buf []byte) (newBuf []byte, err error) {
	ts, err := value.(pgtype.TimestamptzValuer).TimestamptzValue()
	if err != nil {
		return nil, err
	}

	if !ts.Valid {
		return nil, nil
	}

	var microsecSinceY2K int64
	switch ts.InfinityModifier {
	case pgtype.Finite:
		microsecSinceUnixEpoch := ts.Time.Unix()*1000000 + int64(ts.Time.Nanosecond())/1000
		microsecSinceY2K = microsecSinceUnixEpoch - microsecFromUnixEpochToY2K
	case pgtype.Infinity:
		microsecSinceY2K = infinityMicrosecondOffset
	case pgtype.NegativeInfinity:
		microsecSinceY2K = negativeInfinityMicrosecondOffset
	}

	buf = AppendInt64(buf, microsecSinceY2K)

	return buf, nil
}

type encodePlanTimestamptzCodecText struct{}

func (encodePlanTimestamptzCodecText) Encode(value any, buf []byte) (newBuf []byte, err error) {
	ts, err := value.(pgtype.TimestamptzValuer).TimestamptzValue()
	if err != nil {
		return nil, err
	}

	if !ts.Valid {
		return nil, nil
	}

	var s string

	switch ts.InfinityModifier {
	case pgtype.Finite:

		t := ts.Time.UTC().Truncate(time.Microsecond)

		// Year 0000 is 1 BC
		bc := false
		if year := t.Year(); year <= 0 {
			year = -year + 1
			t = time.Date(year, t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
			bc = true
		}

		s = t.Format(pgTimestamptzSecondFormat)

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

func (c *TimestamptzCodec) PlanScan(m *pgtype.Map, oid uint32, format int16, target any) pgtype.ScanPlan {

	switch format {
	case pgtype.BinaryFormatCode:
		name := getInterfaceName(target) //sanek
		switch name {
		case "*pgtype.timeWrapper":
			return &scanPlanBinaryTimestamptzToTimestamptzScanner{location: c.ScanLocation}
		}
	case pgtype.TextFormatCode:
		name := getInterfaceName(target) //sanek
		switch name {
		case "*pgtype.timeWrapper":
			return &scanPlanTextTimestamptzToTimestamptzScanner{location: c.ScanLocation}
		}
	}

	return nil
}

type scanPlanBinaryTimestamptzToTimestamptzScanner struct{ location *time.Location }

func (plan *scanPlanBinaryTimestamptzToTimestamptzScanner) Scan(src []byte, dst any) error {
	scanner := (dst).(pgtype.TimestamptzScanner)

	if src == nil {
		return scanner.ScanTimestamptz(pgtype.Timestamptz{Valid: true})
		//return scanner.ScanTimestamptz(pgtype.Timestamptz{})
	}

	if len(src) != 8 {
		return fmt.Errorf("invalid length for timestamptz: %v", len(src))
	}

	var tstz pgtype.Timestamptz
	microsecSinceY2K := int64(binary.BigEndian.Uint64(src))

	switch microsecSinceY2K {
	case infinityMicrosecondOffset:
		tstz = pgtype.Timestamptz{Valid: true, InfinityModifier: pgtype.Infinity}
	case negativeInfinityMicrosecondOffset:
		tstz = pgtype.Timestamptz{Valid: true, InfinityModifier: -pgtype.Infinity}
	default:
		tim := time.Unix(
			microsecFromUnixEpochToY2K/1000000+microsecSinceY2K/1000000,
			(microsecFromUnixEpochToY2K%1000000*1000)+(microsecSinceY2K%1000000*1000),
		)
		if plan.location != nil {
			tim = tim.In(plan.location)
		}
		tstz = pgtype.Timestamptz{Time: tim, Valid: true}
	}

	return scanner.ScanTimestamptz(tstz)
}

type scanPlanTextTimestamptzToTimestamptzScanner struct{ location *time.Location }

func (plan *scanPlanTextTimestamptzToTimestamptzScanner) Scan(src []byte, dst any) error {
	scanner := (dst).(pgtype.TimestamptzScanner)

	if src == nil {
		return scanner.ScanTimestamptz(pgtype.Timestamptz{})
	}

	var tstz pgtype.Timestamptz
	sbuf := string(src)
	switch sbuf {
	case "infinity":
		tstz = pgtype.Timestamptz{Valid: true, InfinityModifier: pgtype.Infinity}
	case "-infinity":
		tstz = pgtype.Timestamptz{Valid: true, InfinityModifier: -pgtype.Infinity}
	default:
		bc := false
		if strings.HasSuffix(sbuf, " BC") {
			sbuf = sbuf[:len(sbuf)-3]
			bc = true
		}

		var format string
		if len(sbuf) >= 9 && (sbuf[len(sbuf)-9] == '-' || sbuf[len(sbuf)-9] == '+') {
			format = pgTimestamptzSecondFormat
		} else if len(sbuf) >= 6 && (sbuf[len(sbuf)-6] == '-' || sbuf[len(sbuf)-6] == '+') {
			format = pgTimestamptzMinuteFormat
		} else {
			format = pgTimestamptzHourFormat
		}

		tim, err := time.Parse(format, sbuf)
		if err != nil {
			return err
		}

		if bc {
			year := -tim.Year() + 1
			tim = time.Date(year, tim.Month(), tim.Day(), tim.Hour(), tim.Minute(), tim.Second(), tim.Nanosecond(), tim.Location())
		}

		if plan.location != nil {
			tim = tim.In(plan.location)
		}

		tstz = pgtype.Timestamptz{Time: tim, Valid: true}
	}

	return scanner.ScanTimestamptz(tstz)
}

func (c *TimestamptzCodec) DecodeDatabaseSQLValue(m *pgtype.Map, oid uint32, format int16, src []byte) (driver.Value, error) {
	if src == nil {
		return nil, nil
	}

	var tstz Timestamptz
	err := codecScan(c, m, oid, format, src, &tstz)
	if err != nil {
		return nil, err
	}

	if tstz.InfinityModifier != pgtype.Finite {
		return tstz.InfinityModifier.String(), nil
	}

	return tstz.Time, nil
}

func (c *TimestamptzCodec) DecodeValue(m *pgtype.Map, oid uint32, format int16, src []byte) (any, error) {
	if src == nil {
		return nil, nil
	}

	var tstz Timestamptz
	err := codecScan(c, m, oid, format, src, &tstz)
	if err != nil {
		return nil, err
	}

	if tstz.InfinityModifier != pgtype.Finite {
		return tstz.InfinityModifier, nil
	}

	return tstz.Time, nil
}

func codecScan(codec pgtype.Codec, m *pgtype.Map, oid uint32, format int16, src []byte, dst any) error {
	scanPlan := codec.PlanScan(m, oid, format, dst)
	if scanPlan == nil {
		return fmt.Errorf("PlanScan did not find a plan")
	}
	return scanPlan.Scan(src, dst)
}
