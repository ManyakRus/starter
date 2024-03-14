package folders

import (
	"github.com/ManyakRus/starter/micro"
	"testing"
)

func TestFindFolders(t *testing.T) {
	//dir := `/home/user/GolandProjects/!sanek/image_packages/`
	dir := micro.ProgramDir()
	Otvet := FindFoldersTree(dir, true, false, false, []string{"vendor"})
	if Otvet == nil {
		t.Log("TestFindFolders() error: Otvet = nil")
	}
}
