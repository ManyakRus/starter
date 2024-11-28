package microl

import (
	"os"
	"testing"
)

func TestSet_FieldFromEnv(t *testing.T) {
	type Struct1 struct {
		TestSet_FieldFromEnv string
	}

	Name := "TestSet_FieldFromEnv"
	os.Setenv(Name, Name)

	Struct := Struct1{}
	Set_FieldFromEnv_String(&Struct, Name, true)

	if Struct.TestSet_FieldFromEnv != Name {
		t.Error("Set_FieldFromEnv_String() error")
	}

}
