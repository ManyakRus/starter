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

func TestFindFiles_FromDirectory(t *testing.T) {

	_, err := FindFiles_FromDirectory(micro.ProgramDir(), "")
	if err != nil {
		t.Log("TestFindFiles_FromDirectory() error: ", err)
	}

}
