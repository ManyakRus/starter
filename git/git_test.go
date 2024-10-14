package git

import (
	"testing"
)

func TestFind_LastTagVersion(t *testing.T) {

	Otvet, err := Find_LastCommitDescribe()
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

func TestFind_LastCommitTime(t *testing.T) {
	_, err := Find_LastCommitTime()
	if err != nil {
		t.Error("TestFind_LastCommitTime() error: ", err)
	}
}

func TestFind_LastCommitHash(t *testing.T) {
	_, err := Find_LastCommitHash()
	if err != nil {
		t.Error("Find_LastCommitHash() error: ", err)
	}
}

func TestFind_CommitDescribe(t *testing.T) {

	Hash, err := Find_LastCommitHash()
	if err != nil {
		t.Error("TestFind_CommitDescribe() error: ", err)
	}

	_, err = Find_CommitDescribe(Hash)
	if err != nil {
		t.Error("TestFind_CommitDescribe() error: ", err)
	}
}

func TestFind_CommitTime(t *testing.T) {

	Hash, err := Find_LastCommitHash()
	if err != nil {
		t.Error("TestFind_CommitTime() error: ", err)
	}

	_, err = Find_CommitTime(Hash)
	if err != nil {
		t.Error("TestFind_CommitTime() error: ", err)
	}
}

func TestFind_Find_LastCommitHashes(t *testing.T) {

	_, err := Find_LastCommitHashes(1)
	if err != nil {
		t.Error("Find_LastCommitHash() error: ", err)
	}

}
