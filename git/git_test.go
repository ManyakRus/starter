package git

import "testing"

func TestFind_LastTagVersion(t *testing.T) {

	Otvet, err := Find_LastCommitVersion()
	if err != nil {
		t.Error(err)
	}

	if Otvet == "" {
		t.Error("TestFind_LastTagVersion() error: Otvet =''")
	}
}

func TestShow_LastCommitVersion(t *testing.T) {
	Show_LastCommitVersion()
}
