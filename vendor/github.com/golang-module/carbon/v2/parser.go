package carbon

import (
	"time"
)

// Parse parses a standard time string as a Carbon instance.
// 将标准格式时间字符串解析成 Carbon 实例
func (c Carbon) Parse(value string, timezone ...string) Carbon {
	if value == "" || value == "0" || value == "0000-00-00 00:00:00" || value == "0000-00-00" || value == "00:00:00" {
		return c
	}
	if len(timezone) > 0 {
		c.loc, c.Error = getLocationByTimezone(timezone[len(timezone)-1])
	}
	switch value {
	case "now":
		return c.Now(timezone...)
	case "yesterday":
		return c.Yesterday(timezone...)
	case "tomorrow":
		return c.Tomorrow(timezone...)
	}
	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, value, c.loc)
		if err == nil {
			c.time = t
			return c
		}
	}
	c.Error = invalidValueError(value)
	return c
}

// Parse parses a standard time string as a Carbon instance.
// 将标准时间字符串解析成 Carbon 实例
func Parse(value string, timezone ...string) Carbon {
	return NewCarbon().Parse(value, timezone...)
}

// ParseByFormat parses a time string as a Carbon instance by format.
// 通过格式模板将时间字符串解析成 Carbon 实例
func (c Carbon) ParseByFormat(value, format string, timezone ...string) Carbon {
	carbon := c.ParseByLayout(value, format2layout(format), timezone...)
	if carbon.Error != nil {
		carbon.Error = invalidFormatError(value, format)
	}
	return carbon
}

// ParseByFormat parses a time string as a Carbon instance by format.
// 通过格式模板将时间字符串解析成 Carbon 实例
func ParseByFormat(value, format string, timezone ...string) Carbon {
	return NewCarbon().ParseByFormat(value, format, timezone...)
}

// ParseByLayout parses a time string as a Carbon instance by layout.
// 通过布局模板将时间字符串解析成 Carbon 实例
func (c Carbon) ParseByLayout(value, layout string, timezone ...string) Carbon {
	if len(timezone) > 0 {
		c.loc, c.Error = getLocationByTimezone(timezone[len(timezone)-1])
	}
	if c.Error != nil {
		return c
	}
	if value == "" || value == "0" || value == "0000-00-00 00:00:00" || value == "0000-00-00" || value == "00:00:00" {
		return c
	}
	tt, err := time.ParseInLocation(layout, value, c.loc)
	if err != nil {
		c.Error = invalidLayoutError(value, layout)
		return c
	}
	c.time = tt
	return c
}

// ParseByLayout parses a time string as a Carbon instance by layout.
// 通过布局模板将时间字符串解析成 Carbon 实例
func ParseByLayout(value, layout string, timezone ...string) Carbon {
	return NewCarbon().ParseByLayout(value, layout, timezone...)
}
