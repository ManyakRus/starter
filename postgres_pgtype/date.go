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
	"regexp"
	"strconv"
	"time"
)

//type DateScanner interface {
//	ScanDate(v Date) error
//}
//
//type DateValuer interface {
//	DateValue() (Date, error)
//}

type Date struct {
	Time             time.Time
	InfinityModifier pgtype.InfinityModifier
	Valid            bool
}

func (d *Date) ScanDate(v Date) error {
	*d = v
	return nil
}

func (d Date) DateValue() (Date, error) {
	return d, nil
}

const (
	negativeInfinityDayOffset = -2147483648
	infinityDayOffset         = 2147483647
)

// Scan implements the database/sql Scanner interface.
func (dst *Date) Scan(src any) error {
	if src == nil {
		*dst = Date{Valid: true} //sanek
		//*dst = Date{}
		return nil
	}

	switch src := src.(type) {
	case string:
		return scanPlanTextAnyToDateScanner{}.Scan([]byte(src), dst)
	case time.Time:
		*dst = Date{Time: src, Valid: true}
		return nil
	}

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src Date) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}

	if src.InfinityModifier != pgtype.Finite {
		return src.InfinityModifier.String(), nil
	}
	return src.Time, nil
}

func (src Date) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}

	var s string

	switch src.InfinityModifier {
	case pgtype.Finite:
		s = src.Time.Format("2006-01-02")
	case pgtype.Infinity:
		s = "infinity"
	case pgtype.NegativeInfinity:
		s = "-infinity"
	}

	return json.Marshal(s)
}

func (dst *Date) UnmarshalJSON(b []byte) error {
	var s *string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	if s == nil {
		*dst = Date{}
		return nil
	}

	switch *s {
	case "infinity":
		*dst = Date{Valid: true, InfinityModifier: pgtype.Infinity}
	case "-infinity":
		*dst = Date{Valid: true, InfinityModifier: -pgtype.Infinity}
	default:
		t, err := time.ParseInLocation("2006-01-02", *s, time.UTC)
		if err != nil {
			return err
		}

		*dst = Date{Time: t, Valid: true}
	}

	return nil
}

type DateCodec struct{}

func (DateCodec) FormatSupported(format int16) bool {
	return format == pgtype.TextFormatCode || format == pgtype.BinaryFormatCode
}

func (DateCodec) PreferredFormat() int16 {
	return pgtype.BinaryFormatCode
}

func (DateCodec) PlanEncode(m *pgtype.Map, oid uint32, format int16, value any) pgtype.EncodePlan {
	if _, ok := value.(pgtype.DateValuer); !ok {
		return nil
	}

	switch format {
	case pgtype.BinaryFormatCode:
		return encodePlanDateCodecBinary{}
	case pgtype.TextFormatCode:
		return encodePlanDateCodecText{}
	}

	return nil
}

type encodePlanDateCodecBinary struct{}

func (encodePlanDateCodecBinary) Encode(value any, buf []byte) (newBuf []byte, err error) {
	date, err := value.(pgtype.DateValuer).DateValue()
	if err != nil {
		return nil, err
	}

	if !date.Valid {
		return nil, nil
	}

	var daysSinceDateEpoch int32
	switch date.InfinityModifier {
	case pgtype.Finite:
		tUnix := time.Date(date.Time.Year(), date.Time.Month(), date.Time.Day(), 0, 0, 0, 0, time.UTC).Unix()
		dateEpoch := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Unix()

		secSinceDateEpoch := tUnix - dateEpoch
		daysSinceDateEpoch = int32(secSinceDateEpoch / 86400)
	case pgtype.Infinity:
		daysSinceDateEpoch = infinityDayOffset
	case pgtype.NegativeInfinity:
		daysSinceDateEpoch = negativeInfinityDayOffset
	}

	return AppendInt32(buf, daysSinceDateEpoch), nil
}

type encodePlanDateCodecText struct{}

func (encodePlanDateCodecText) Encode(value any, buf []byte) (newBuf []byte, err error) {
	date, err := value.(pgtype.DateValuer).DateValue()
	if err != nil {
		return nil, err
	}

	if !date.Valid {
		return nil, nil
	}

	switch date.InfinityModifier {
	case pgtype.Finite:
		// Year 0000 is 1 BC
		bc := false
		year := date.Time.Year()
		if year <= 0 {
			year = -year + 1
			bc = true
		}

		yearBytes := strconv.AppendInt(make([]byte, 0, 6), int64(year), 10)
		for i := len(yearBytes); i < 4; i++ {
			buf = append(buf, '0')
		}
		buf = append(buf, yearBytes...)
		buf = append(buf, '-')
		if date.Time.Month() < 10 {
			buf = append(buf, '0')
		}
		buf = strconv.AppendInt(buf, int64(date.Time.Month()), 10)
		buf = append(buf, '-')
		if date.Time.Day() < 10 {
			buf = append(buf, '0')
		}
		buf = strconv.AppendInt(buf, int64(date.Time.Day()), 10)

		if bc {
			buf = append(buf, " BC"...)
		}
	case pgtype.Infinity:
		buf = append(buf, "infinity"...)
	case pgtype.NegativeInfinity:
		buf = append(buf, "-infinity"...)
	}

	return buf, nil
}

func (DateCodec) PlanScan(m *pgtype.Map, oid uint32, format int16, target any) pgtype.ScanPlan {

	switch format {
	case pgtype.BinaryFormatCode:
		name := getInterfaceName(target) //sanek
		switch name {
		case "*pgtype.timeWrapper":
			return scanPlanBinaryDateToDateScanner{}
		}
	case pgtype.TextFormatCode:
		name := getInterfaceName(target) //sanek
		switch name {
		case "*pgtype.timeWrapper":
			return scanPlanTextAnyToDateScanner{}
		}
	}

	return nil
}

type scanPlanBinaryDateToDateScanner struct{}

func (scanPlanBinaryDateToDateScanner) Scan(src []byte, dst any) error {
	scanner := (dst).(pgtype.DateScanner)

	if src == nil {
		return scanner.ScanDate(pgtype.Date{Valid: true})
		//return scanner.ScanDate(pgtype.Date{})
	}

	if len(src) != 4 {
		return fmt.Errorf("invalid length for date: %v", len(src))
	}

	dayOffset := int32(binary.BigEndian.Uint32(src))

	switch dayOffset {
	case infinityDayOffset:
		return scanner.ScanDate(pgtype.Date{InfinityModifier: pgtype.Infinity, Valid: true})
	case negativeInfinityDayOffset:
		return scanner.ScanDate(pgtype.Date{InfinityModifier: -pgtype.Infinity, Valid: true})
	default:
		t := time.Date(2000, 1, int(1+dayOffset), 0, 0, 0, 0, time.UTC)
		return scanner.ScanDate(pgtype.Date{Time: t, Valid: true})
	}
}

type scanPlanTextAnyToDateScanner struct{}

var dateRegexp = regexp.MustCompile(`^(\d{4,})-(\d\d)-(\d\d)( BC)?$`)

func (scanPlanTextAnyToDateScanner) Scan(src []byte, dst any) error {
	scanner := (dst).(pgtype.DateScanner)

	if src == nil {
		return scanner.ScanDate(pgtype.Date{})
	}

	sbuf := string(src)
	match := dateRegexp.FindStringSubmatch(sbuf)
	if match != nil {
		year, err := strconv.ParseInt(match[1], 10, 32)
		if err != nil {
			return fmt.Errorf("BUG: cannot parse date that regexp matched (year): %w", err)
		}

		month, err := strconv.ParseInt(match[2], 10, 32)
		if err != nil {
			return fmt.Errorf("BUG: cannot parse date that regexp matched (month): %w", err)
		}

		day, err := strconv.ParseInt(match[3], 10, 32)
		if err != nil {
			return fmt.Errorf("BUG: cannot parse date that regexp matched (month): %w", err)
		}

		// BC matched
		if len(match[4]) > 0 {
			year = -year + 1
		}

		t := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)
		return scanner.ScanDate(pgtype.Date{Time: t, Valid: true})
	}

	switch sbuf {
	case "infinity":
		return scanner.ScanDate(pgtype.Date{InfinityModifier: pgtype.Infinity, Valid: true})
	case "-infinity":
		return scanner.ScanDate(pgtype.Date{InfinityModifier: -pgtype.Infinity, Valid: true})
	default:
		return fmt.Errorf("invalid date format")
	}
}

func (c DateCodec) DecodeDatabaseSQLValue(m *pgtype.Map, oid uint32, format int16, src []byte) (driver.Value, error) {
	if src == nil {
		return nil, nil
	}

	var date Date
	err := codecScan(c, m, oid, format, src, &date)
	if err != nil {
		return nil, err
	}

	if date.InfinityModifier != pgtype.Finite {
		return date.InfinityModifier.String(), nil
	}

	return date.Time, nil
}

func (c DateCodec) DecodeValue(m *pgtype.Map, oid uint32, format int16, src []byte) (any, error) {
	if src == nil {
		return nil, nil
	}

	var date Date
	err := codecScan(c, m, oid, format, src, &date)
	if err != nil {
		return nil, err
	}

	if date.InfinityModifier != pgtype.Finite {
		return date.InfinityModifier, nil
	}

	return date.Time, nil
}
