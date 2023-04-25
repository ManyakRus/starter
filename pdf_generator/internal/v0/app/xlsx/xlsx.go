package xlsx

import (
	// "github.com/manyakrus/starter/common/v0/contextmain"
	// stopapp "github.com/manyakrus/starter/common/v0/stopapp"

	"errors"
	"fmt"
	xlst "github.com/ivahaev/go-xlsx-templater"
	logger "github.com/manyakrus/starter/common/v0/logger"
)

// log - глобальный логгер
var log = logger.GetLog()

// CreateXLSX - заполняет шаблон файла FilenameIn из MapReplace. Создаёт файл FilenameOut.
func CreateXLSX(FilenameIn, FilenameOut string, map1 map[string]interface{}) error {
	var err error

	doc := xlst.New()
	err = doc.ReadTemplate(FilenameIn)
	if err != nil {
		Text1 := fmt.Sprint("doc.ReadTemplate() FilenameIn: ", FilenameIn, " error: ", err)
		log.Error(Text1)
		return errors.New(Text1)
	}

	err = doc.Render(map1)
	if err != nil {
		Text1 := fmt.Sprint("doc.Render() FilenameIn: ", FilenameIn, " error: ", err)
		log.Error(Text1)
		return errors.New(Text1)
	}

	err = doc.Save(FilenameOut)
	if err != nil {
		Text1 := fmt.Sprint("doc.Save() FilenameOut: ", FilenameOut, " error: ", err)
		log.Error(Text1)
		return errors.New(Text1)
	}

	return err
}
