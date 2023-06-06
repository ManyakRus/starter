package folders

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"path/filepath"
)

type File struct {
	Name string
}

type Folder struct {
	Name    string
	Files   []*File
	Folders map[string]*Folder
}

func (f *Folder) String() string {
	j, _ := json.Marshal(f)
	return string(j)
}

// FindFoldersTree - возвращает дерево каталогов и файлов, начиная с директории dir
func FindFoldersTree(dir string, NeedFolders, NeedFiles, NeedDot bool, exclude string) *Folder {
	dir = path.Clean(dir)
	var tree *Folder
	var nodes = map[string]interface{}{}
	var walkFun filepath.WalkFunc = func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			nodes[p] = &Folder{path.Base(p), []*File{}, map[string]*Folder{}}
		} else {
			nodes[p] = &File{path.Base(p)}
		}
		return nil
	}
	err := filepath.Walk(dir, walkFun)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range nodes {
		var parentFolder *Folder
		if key == dir {
			tree = value.(*Folder)
			continue
		} else {
			parentFolder = nodes[path.Dir(key)].(*Folder)
		}

		switch v := value.(type) {
		case *File:
			if NeedFiles == false {
				break
			}
			if NeedDot == false && len(v.Name) > 0 && v.Name[0:1] == "." {
				break
			}
			if exclude != "" && len(v.Name) >= len(exclude) && v.Name[0:len(exclude)] == exclude {
				break
			}
			parentFolder.Files = append(parentFolder.Files, v)
		case *Folder:
			if NeedFolders == false {
				break
			}
			if NeedDot == false && len(v.Name) > 0 && v.Name[0:1] == "." {
				break
			}
			if exclude != "" && len(v.Name) >= len(exclude) && v.Name[0:len(exclude)] == exclude {
				break
			}
			parentFolder.Folders[v.Name] = v
		}
	}

	return tree
}
