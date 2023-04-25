package txt

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/manyakrus/starter/common/v0/logger"
	"io/ioutil"
	"os"
)

// log - глобальный логгер
var log = logger.GetLog()

// CreateDocx1 - заполняет шаблон файла FilenameIn из MapReplace. Создаёт файл FilenameOut.
func CreateTxt(FilenameIn, FilenameOut string, MapReplace map[string]string) (err error) {
	//var err error

	// open output file FilenameOut
	FileOut, err := os.Create(FilenameOut)
	if err != nil {
		Text1 := fmt.Sprint("os.Create() FilenameOut: ", FilenameOut, " error: ", err)
		log.Error(Text1)
		return errors.New(Text1)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := FileOut.Close(); err != nil {
			Text1 := fmt.Sprint("fo.Close() FilenameOut: ", FilenameOut, " error: ", err)
			log.Error(Text1)
			err = errors.New(Text1)
		}
	}()

	// Read Write Mode FilenameIn
	FileIn, err := os.OpenFile(FilenameIn, os.O_RDWR, 0644)

	if err != nil {
		Text1 := fmt.Sprint("os.OpenFile() FilenameIn: ", FilenameIn, " error: ", err)
		log.Error(Text1)
		return errors.New(Text1)
	}
	defer func() {
		err = FileIn.Close()
		if err != nil {
			Text1 := fmt.Sprint("file.Close() FilenameIn: ", FilenameIn, " error: ", err)
			log.Error(Text1)
		}
	}()

	// заполнение
	MassB := make([]byte, 0)
	MassB, err = ioutil.ReadAll(FileIn)
	//count, err := FileIn.Read(MassB)
	if err != nil {
		Text1 := fmt.Sprint("file.Read() FilenameIn: ", FilenameIn, " error: ", err)
		log.Error(Text1)
		return errors.New(Text1)
	}
	if len(MassB) == 0 {
		Text1 := fmt.Sprint("file.Read() FilenameIn: ", FilenameIn, " count=0 ")
		log.Error(Text1)
		return errors.New(Text1)
	}

	for k, v := range MapReplace {
		sOld := []byte(k)
		sNew := []byte(v)
		MassB = bytes.ReplaceAll(MassB, sOld, sNew)
	}

	count, err := FileOut.Write(MassB)
	if err != nil {
		Text1 := fmt.Sprint("FileOut.Write() FilenameIn: ", FilenameIn, " error: ", err)
		log.Error(Text1)
		return errors.New(Text1)
	}
	if count == 0 {
		Text1 := fmt.Sprint("FileOut.Write() FilenameIn: ", FilenameIn, " count=0 ")
		log.Error(Text1)
		return errors.New(Text1)
	}

	return err
}
