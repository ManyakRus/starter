package carbon

// SetTestNow sets a test Carbon instance (real or mock) to be returned when a "now" instance is created.
func SetTestNow(carbon Carbon) Carbon {
	return NewCarbon().SetTestNow(carbon)
}

// SetTestNow sets a test Carbon instance (real or mock) to be returned when a "now" instance is created.
func (c Carbon) SetTestNow(carbon Carbon) Carbon {
	c.testNow = carbon.TimestampNano()
	if carbon.loc != nil {
		c.loc = carbon.loc
	}
	return c
}

// ClearTestNow clears a test Carbon instance (real or mock) to be returned when a "now" instance is created.
func (c Carbon) ClearTestNow() Carbon {
	c.testNow = 0
	return c
}

// HasTestNow report whether there is testing time now.
func (c Carbon) HasTestNow() bool {
	return c.testNow > 0
}
