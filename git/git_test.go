package git

import "testing"

func TestFind_LastTagVersion(t *testing.T) {

	Otvet, err := Find_LastTagVersion()
	if err != nil {
		t.Error(err)
	}

	if Otvet == "" {
		t.Error("TestFind_LastTagVersion() error: Otvet =''")
	}
}
