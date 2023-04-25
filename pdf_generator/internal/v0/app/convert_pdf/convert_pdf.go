package convert_pdf

import (
	"bytes"
	"fmt"
	"github.com/manyakrus/starter/common/v0/logger"
	"os/exec"
	"path/filepath"
)

// log - глобальный логгер
var log = logger.GetLog()

// ConvertToPdf - конвертирует в формат .pdf любой файл понимаемый LibreOffice
// в том числе .xlsx, .docx, .txt, .fodt, .fods
func ConvertToPdf(FilenameXLSX string) error {
	var err error

	sDir := filepath.Dir(FilenameXLSX)
	cmd := exec.Command("soffice", "--headless", "--convert-to", "pdf", "--outdir", sDir, FilenameXLSX)
	//cmd := exec.Command("soffice", "--headless", "--convert-to", "pdf", FilenameXLSX)
	//cmd := exec.Command(TextCommand)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	err = cmd.Start()
	if err != nil {
		sError := fmt.Sprintf("error: %v: ", err)
		log.Errorf(sError)
		return err
	}
	err = cmd.Wait()

	//log.Debugf("Command finished with error: %v: ", err)
	log.Debugf("Convert to pdf finished%v: ", buf.String())

	return err
}
