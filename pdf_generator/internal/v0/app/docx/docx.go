package docx

import (
	"errors"
	"fmt"
	"github.com/nguyenthenguyen/docx"
	"github.com/ManyakRus/starter/common/v0/logger"
)

// log - глобальный логгер
var log = logger.GetLog()

// CreateDocx - заполняет шаблон файла FilenameIn из MapReplace. Создаёт файл FilenameOut.
func CreateDocx(FilenameIn, FilenameOut string, MapReplace map[string]string) error {
	var err error

	//dir1 := filepath.Dir(filename)

	// Read from docx file
	RDocx, err := docx.ReadDocxFile(FilenameIn)
	if err != nil {
		Text1 := fmt.Sprint("ReadDocxFile() FilenameIn: ", FilenameIn, " error: ", err)
		log.Error(Text1)
		return errors.New(Text1)
	}
	docx1 := RDocx.Editable()
	// Replace like https://golang.org/pkg/strings/#Replace

	for k, v := range MapReplace {
		err = docx1.Replace(k, v, -1)
		if err != nil {
			log.Debug("docx1.Replace() FilenameIn: ", FilenameIn, "error: ", err)
		}
	}

	//docx1.Replace("old_1_1", "new_1_1", -1)
	err = docx1.WriteToFile(FilenameOut)
	if err != nil {
		Text1 := fmt.Sprint("WriteToFile() FilenameOut: ", FilenameOut, " error: ", err)
		log.Error(Text1)
		return errors.New(Text1)
	}

	err = RDocx.Close()
	if err != nil {
		Text1 := fmt.Sprint("RDocx.Close() FilenameIn: ", FilenameIn, " error: ", err)
		log.Error(Text1)
		return errors.New(Text1)
	}

	return err
}
