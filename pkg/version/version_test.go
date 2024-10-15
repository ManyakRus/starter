package version

import (
	"github.com/ManyakRus/starter/micro"
	"testing"
)

func TestShow_Version(t *testing.T) {

	micro.Show_Version(Version)

}
