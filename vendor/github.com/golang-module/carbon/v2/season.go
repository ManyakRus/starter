package carbon

import (
	"strings"
)

// Season gets season name according to the meteorological division method like "Spring", i18n is supported.
// 获取当前季节(以气象划分)，支持i18n
func (c Carbon) Season() string {
	if c.IsInvalid() {
		return ""
	}
	if len(c.lang.resources) == 0 {
		c.lang.SetLocale(defaultLocale)
	}
	index := -1
	month := c.Month()
	switch {
	case month == 3 || month == 4 || month == 5:
		index = 0
	case month == 6 || month == 7 || month == 8:
		index = 1
	case month == 9 || month == 10 || month == 11:
		index = 2
	case month == 12 || month == 1 || month == 2:
		index = 3
	}
	if seasons, ok := c.lang.resources["seasons"]; ok {
		slice := strings.Split(seasons, "|")
		if len(slice) == 4 {
			return slice[index]
		}
	}
	return ""
}

// StartOfSeason returns a Carbon instance for start of the season.
// 本季节开始时间
func (c Carbon) StartOfSeason() Carbon {
	if c.IsInvalid() {
		return c
	}
	year, month, _ := c.Date()
	if month == 1 || month == 2 {
		return c.create(year-1, 12, 1, 0, 0, 0, 0)
	}
	return c.create(year, month/3*3, 1, 0, 0, 0, 0)
}

// EndOfSeason returns a Carbon instance for end of the season.
// 本季节结束时间
func (c Carbon) EndOfSeason() Carbon {
	if c.IsInvalid() {
		return c
	}
	year, month, _ := c.Date()
	if month == 1 || month == 2 {
		return c.create(year, 3, 0, 23, 59, 59, 999999999)
	}
	if month == 12 {
		return c.create(year+1, 3, 0, 23, 59, 59, 999999999)
	}
	return c.create(year, month/3*3+3, 0, 23, 59, 59, 999999999)
}

// IsSpring reports whether is spring.
// 是否是春季
func (c Carbon) IsSpring() bool {
	if c.IsInvalid() {
		return false
	}
	month := c.Month()
	if month == 3 || month == 4 || month == 5 {
		return true
	}
	return false
}

// IsSummer reports whether is summer.
// 是否是夏季
func (c Carbon) IsSummer() bool {
	if c.IsInvalid() {
		return false
	}
	month := c.Month()
	if month == 6 || month == 7 || month == 8 {
		return true
	}
	return false
}

// IsAutumn reports whether is autumn.
// 是否是秋季
func (c Carbon) IsAutumn() bool {
	if c.IsInvalid() {
		return false
	}
	month := c.Month()
	if month == 9 || month == 10 || month == 11 {
		return true
	}
	return false
}

// IsWinter reports whether is winter.
// 是否是冬季
func (c Carbon) IsWinter() bool {
	if c.IsInvalid() {
		return false
	}
	month := c.Month()
	if month == 12 || month == 1 || month == 2 {
		return true
	}
	return false
}
