package contextmain

import (
	"testing"
)

func TestGetContext(t *testing.T) {
	ctx := GetContext()
	if ctx == nil {
		t.Error("contextmain_test.TestGetContext() Wrong GetContext() !")
	}
}

func TestGetNewContext(t *testing.T) {
	ctx := GetNewContext()
	if ctx == nil {
		t.Error("contextmain_test.TestGetNewContext() Wrong GetContext() !")
	}
}
