// The functions and methods in this file will be removed in the future
// 该文件中的函数和方法将来会被移除

package carbon

import "time"

// Deprecated: It will be removed in the future, use CreateFromStdTime instead.
//
// FromStdTime converts standard time.Time to Carbon.
// 将标准 time.Time 转换成 Carbon，未来将移除，请使用 CreateFromStdTime 替代
func FromStdTime(tt time.Time) Carbon {
	return CreateFromStdTime(tt)
}

// Deprecated: It will be removed in the future, use CreateFromStdTime instead.
//
// Time2Carbon converts standard time.Time to Carbon.
// 将标准 time.Time 转换成 Carbon，未来将移除，请使用 CreateFromStdTime 替代
func Time2Carbon(tt time.Time) Carbon {
	return CreateFromStdTime(tt)
}

// Deprecated: It will be removed in the future, use ToStdTime instead.
//
// Carbon2Time converts Carbon to standard time.Time.
// 将 Carbon 转换成标准 time.Time，未来将移除，请使用 ToStdTime 替代
func (c Carbon) Carbon2Time() time.Time {
	return c.ToStdTime()
}
