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
	FileName string
	Name     string
	Files    []*File
	Folders  map[string]*Folder
}

func (f *Folder) String() string {
	j, _ := json.Marshal(f)
	return string(j)
}

// FindFoldersTree - возвращает дерево каталогов и файлов, начиная с директории dir
func FindFoldersTree(dir string, NeedFolders, NeedFiles, NeedDot bool, MassExclude []string) *Folder {
	dir = path.Clean(dir)
	var tree *Folder
	var nodes = map[string]interface{}{}
	var walkFun filepath.WalkFunc = func(p string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		// проверка кроме MassExclude
		if len(MassExclude) > 0 {
			for _, v := range MassExclude {
				if info.Name() == v {
					return nil
				}
				if p == v {
					return nil
				}
			}
		}

		if info.IsDir() {
			nodes[p] = &Folder{p, path.Base(p), []*File{}, map[string]*Folder{}}
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
			var ok bool
			parentFolder, ok = nodes[path.Dir(key)].(*Folder)
			if !ok {
				continue
			}
		}

		// найдём название Папки/Файла
		var Name string
		//var FolderName string
		switch value.(type) {
		case *File:
			Name = value.(*File).Name
		case *Folder:
			{
				Name = value.(*Folder).Name
				//FolderName = value.(*Folder).FileName
			}
		}

		// проверка скрытые файлы с точкой
		if NeedDot == false && len(Name) > 0 && Name[0:1] == "." {
			continue
		}

		//// проверка кроме MassExclude
		//if len(MassExclude) > 0 {
		//	for _, v := range MassExclude {
		//		if Name == v {
		//			continue
		//		}
		//		if FolderName == v {
		//			continue
		//		}
		//	}
		//}

		//
		switch v := value.(type) {
		case *File:
			if NeedFiles == false {
				break
			}
			parentFolder.Files = append(parentFolder.Files, v)
		case *Folder:
			if NeedFolders == false {
				break
			}
			parentFolder.Folders[v.Name] = v
		}
	}

	return tree
}
