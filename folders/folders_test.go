package folders

import (
	"github.com/ManyakRus/starter/micro"
	"testing"
)

func TestFindFolders(t *testing.T) {
	//dir := `/home/user/GolandProjects/!sanek/image_packages/`
	dir := micro.ProgramDir()
	FindFoldersTree(dir, true, false, false, "vendor")
}
