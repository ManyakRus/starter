package humanize

import "testing"

func TestStringFromFloat64_underline(t *testing.T) {

	Otvet := StringFromFloat64_underline(1123.456)
	if Otvet != "1_123.46" {
		t.Error("Error")
	}
}

func TestStringFromFloat32_underline(t *testing.T) {

	Otvet := StringFromFloat32_underline(1123.456)
	if Otvet != "1_123.46" {
		t.Error("Error")
	}
}

func TestStringFromInt_underline(t *testing.T) {
	Otvet := StringFromInt_underline(1123)
	if Otvet != "1_123" {
		t.Error("Error")
	}
}

func TestStringFromInt64_underline(t *testing.T) {
	Otvet := StringFromInt64_underline(1123)
	if Otvet != "1_123" {
		t.Error("Error")
	}
}

func TestStringFromInt32_underline(t *testing.T) {
	Otvet := StringFromInt32_underline(1123)
	if Otvet != "1_123" {
		t.Error("Error")
	}
}
