package carbon

import (
	"bytes"
	"fmt"
)

// MarshalJSON implements the interface Marshaler for Carbon struct.
// 实现 Marshaler 接口
func (c Carbon) MarshalJSON() ([]byte, error) {
	if c.Error != nil {
		return nil, c.Error
	}
	key, value := c.parseTag()
	if key == "layout" {
		return []byte(fmt.Sprintf(`"%s"`, c.Layout(value))), nil
	}
	if key == "format" {
		return []byte(fmt.Sprintf(`"%s"`, c.Format(value))), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, c.ToDateTimeString())), nil
}

// UnmarshalJSON implements the interface Unmarshaler for Carbon struct.
// 实现 Unmarshaler 接口
func (c *Carbon) UnmarshalJSON(b []byte) error {
	if c.Error != nil {
		return c.Error
	}
	key, value := c.parseTag()
	if key == "layout" {
		*c = ParseByLayout(string(bytes.Trim(b, `"`)), value, c.Location())
		c.tag = "layout:" + value
		return c.Error
	}
	if key == "format" {
		*c = ParseByFormat(string(bytes.Trim(b, `"`)), value, c.Location())
		c.tag = "format:" + value
		return c.Error
	}
	*c = Parse(string(bytes.Trim(b, `"`)), c.Location())
	return c.Error
}
