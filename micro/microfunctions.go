// модуль с вспомогательными небольшими функциями

package micro

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"unicode"

	//"log"
	"os"
	"path/filepath"
	"time"
)

//var log = logger.GetLog()

// IsTestApp - возвращает true если это тестовая среда выполнения приложения
func IsTestApp() bool {
	Otvet := true

	stage, ok := os.LookupEnv("STAGE")
	if ok == false {
		panic(fmt.Errorf("Not found Env 'STAGE' !"))
	}

	switch stage {
	case "local", "dev", "test", "preprod":
		Otvet = true
	case "prod":
		Otvet = false
	default:
		panic(fmt.Errorf("Error, unknown stage(%v) !", stage))
	}

	return Otvet
}

// FileExists - возвращает true если файл существует
func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

// AddSeparator - добавляет в конец строки сеператор "/", если его там нет
func AddSeparator(dir string) string {
	otvet := dir

	if otvet == "" {
		return SeparatorFile()
	}

	if otvet[len(otvet)-1:] != SeparatorFile() {
		otvet = otvet + SeparatorFile()
	}

	return otvet
}

// SeparatorFile - возвращает символ сепаратора каталогов= / или \
func SeparatorFile() string {
	return string(filepath.Separator)
}

// Sleep - приостановка работы программы на нужное число миллисекунд
func Sleep(ms int) {
	duration := time.Duration(ms) * time.Millisecond
	time.Sleep(duration)
}

// Pause - приостановка работы программы на нужное число миллисекунд
func Pause(ms int) {
	Sleep(ms)
}

// FindDirUp - возвращает строку с именем каталога на уровень выше
func FindDirUp(dir string) string {
	otvet := dir
	if dir == "" {
		return otvet
	}

	if otvet[len(otvet)-1:] == SeparatorFile() {
		otvet = otvet[:len(otvet)-1]
	}

	pos1 := strings.LastIndex(otvet, SeparatorFile())
	if pos1 > 0 {
		otvet = otvet[0 : pos1+1]
	}

	return otvet
}

// ErrorJoin - возвращает ошибку из объединения текста двух ошибок
func ErrorJoin(err1, err2 error) error {
	var err error

	if err1 == nil && err2 == nil {

	} else if err1 == nil {
		err = err2
	} else if err2 == nil {
		err = err1
	} else {
		err = errors.New(err1.Error() + ", " + err2.Error())
	}

	return err
}

// SubstringLeft - возвращает левые символы строки
func SubstringLeft(str string, num int) string {
	if num <= 0 {
		return ``
	}
	if num > len(str) {
		num = len(str)
	}
	return str[:num]
}

// SubstringRight - возвращает правые символы строки
func SubstringRight(str string, num int) string {
	if num <= 0 {
		return ``
	}
	max := len(str)
	if num > max {
		num = max
	}
	num = max - num
	return str[num:]
}

// StringBetween - GetStringInBetween Returns empty string if no start string found
func StringBetween(str string, start string, end string) string {
	otvet := ""
	if str == "" {
		return otvet
	}

	pos1 := strings.Index(str, start)
	if pos1 == -1 {
		return otvet
	}
	pos1 += len(start)

	pos2 := strings.Index(str[pos1:], end)
	if pos2 == -1 {
		return otvet
	}
	pos2 = pos1 + pos2

	otvet = str[pos1:pos2]
	return otvet
}

// LastWord - возвращает последнее слово из строки
func LastWord(StringFrom string) string {
	Otvet := ""

	if StringFrom == "" {
		return Otvet
	}

	r := []rune(StringFrom)
	for f := len(r); f >= 0; f-- {
		r1 := r[f-1]
		if r1 == '_' {
		} else if unicode.IsLetter(r1) == false && unicode.IsDigit(r1) == false {
			break
		}

		Otvet = string(r1) + Otvet
	}

	return Otvet
}

// CurrentFilename - возвращает полное имя текущего исполняемого файла
func CurrentFilename() string {
	_, filename, _, _ := runtime.Caller(0)
	return filename
}

// ProgramDir - возвращает главный каталог программы, в конце "/"
func ProgramDir_Common() string {
	filename := os.Args[0]
	dir := filepath.Dir(filename)
	sdir := strings.ToLower(dir)

	if SubstringLeft(sdir, 5) == "/tmp/" {
		filename = CurrentFilename()
		dir = filepath.Dir(filename)

		substr := SeparatorFile() + "vendor" + SeparatorFile()
		pos_vendor := strings.Index(strings.ToLower(dir), substr)
		if pos_vendor >= 0 {
			dir = dir[0:pos_vendor]
		} else if dir[len(dir)-5:] == "micro" {
			dir = FindDirUp(dir)
			//dir = FindDirUp(dir)
			//dir = FindDirUp(dir)
		}
	}

	//dir, err := os.Getwd()
	//if err != nil {
	//	log.Fatalln(err)
	//	dir = ""
	//}

	dir = AddSeparator(dir)
	return dir
}

// ProgramDir - возвращает главный каталог программы, в конце "/"
func ProgramDir() string {
	Otvet := ProgramDir_Common()
	return Otvet
}

// FileNameWithoutExtension - возвращает имя файла без расширения
func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func BeginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

func EndOfMonth(date time.Time) time.Time {
	Otvet := date.AddDate(0, 1, -date.Day())
	return Otvet
	//return date.AddDate(0, 1, -date.Day())
}

//// GetPackageName - возвращает имя пакета
//func GetPackageName(temp interface{}) string {
//	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")
//	strs = strings.Split(strs[len(strs)-2], "/")
//	return strs[len(strs)-1]
//}

// StringAfter - возвращает строку, начиная после субстроки StringAfter
func StringAfter(StringFull, StringAfter string) string {
	Otvet := StringFull
	pos1 := strings.Index(StringFull, StringAfter)
	if pos1 == -1 {
		return Otvet
	}

	Otvet = Otvet[pos1+len(StringAfter):]

	return Otvet
}

// StringFrom - возвращает строку, начиная со субстроки StringAfter
func StringFrom(StringFull, StringAfter string) string {
	Otvet := StringFull
	pos1 := strings.Index(StringFull, StringAfter)
	if pos1 == -1 {
		return Otvet
	}

	Otvet = Otvet[pos1:]

	return Otvet
}

func Trim(s string) string {
	Otvet := ""

	Otvet = strings.Trim(s, " \n\r\t")

	return Otvet
}
