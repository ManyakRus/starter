package humanize

import "github.com/dustin/go-humanize"

// StringFromFloat64_underline - преобразование float64 в строку, с разделителями тысяч
func StringFromFloat64_underline(f float64) string {
	Otvet := ""

	Otvet = humanize.FormatFloat("#_###.##", f)

	return Otvet
}

// StringFromFloat32_underline - преобразование float64 в строку, с разделителями тысяч
func StringFromFloat32_underline(f float32) string {
	Otvet := ""

	Otvet = humanize.FormatFloat("#_###.##", float64(f))

	return Otvet
}

// StringFromInt_underline - преобразование int в строку, с разделителями тысяч
func StringFromInt_underline(i int) string {
	Otvet := ""

	Otvet = humanize.FormatInteger("#_###.", i)

	return Otvet
}

// StringFromInt32_underline - преобразование int32 в строку, с разделителями тысяч
func StringFromInt32_underline(i int32) string {
	Otvet := ""

	Otvet = humanize.FormatInteger("#_###.", int(i))

	return Otvet
}

// StringFromInt64_underline - преобразование int64 в строку, с разделителями тысяч
func StringFromInt64_underline(i int64) string {
	Otvet := ""

	Otvet = humanize.FormatInteger("#_###.", int(i))

	return Otvet
}
