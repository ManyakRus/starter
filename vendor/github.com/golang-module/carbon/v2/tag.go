package carbon

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	// invalid pointer error
	// 无效的指针错误
	invalidPtrError = func() error {
		return fmt.Errorf("invalid pointer, please make sure the pointer is valid")
	}

	// invalid tag error
	// 无效的标签错误
	invalidTagError = func() error {
		return fmt.Errorf("invalid tag, please make sure the tag is valid")
	}
)

// Tag gets tag.
// 获取标签
func (c Carbon) Tag() string {
	return c.tag
}

// SetTag sets tag.
// 设置标签
func (c Carbon) SetTag(tag string) Carbon {
	if c.Error != nil {
		return c
	}
	c.tag = tag
	return c
}

// parseTag parses tag.
// 解析标签
func (c Carbon) parseTag() (key, value string) {
	if c.tag == "" || len(c.tag) <= 7 {
		return "", ""
	}
	return strings.TrimSpace(c.tag[:6]), strings.TrimSpace(c.tag[7:])
}

// LoadTag loads tag.
// 加载标签
func LoadTag(v interface{}) error {
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if typ.Kind() != reflect.Ptr {
		return invalidPtrError()
	}
	params := make([]reflect.Value, 1)
	for i := 0; i < val.Elem().NumField(); i++ {
		field := typ.Elem().Field(i)
		value := val.Elem().Field(i)
		if reflect.TypeOf(Carbon{}) != value.Type() {
			continue
		}
		tag := field.Tag.Get("carbon")
		if tag == "" {
			tag = "layout:" + DateTimeLayout
		}
		if !strings.Contains(tag, "layout:") && !strings.Contains(tag, "format:") {
			return invalidTagError()
		}
		params[0] = reflect.ValueOf(tag)
		value.Set(value.MethodByName("SetTag").Call(params)[0])
	}
	return nil
}
