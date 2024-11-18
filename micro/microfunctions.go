// модуль с вспомогательными небольшими функциями

package micro

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
	"hash/fnv"
	"os/exec"
	"reflect"
	"runtime"
	"sort"
	"strconv"
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

// Pause_ctx - приостановка работы программы на нужное число миллисекунд, с учётом глобального контекста
func Pause_ctx(ctx context.Context, ms int) {

	Duration := time.Duration(ms) * time.Millisecond

	select {
	case <-ctx.Done():
	case <-time.After(Duration):
	}
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
	for f := len(r); f > 0; f-- {
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
	//filename := os.Args[0]
	filename, err := os.Executable()
	if err != nil {
		panic(err)
	}

	dir := filepath.Dir(filename)
	sdir := strings.ToLower(dir)

	substr := "/tmp/"
	pos1 := strings.Index(sdir, substr)
	if pos1 >= 0 {
		//linux
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
	} else {
		//Windows
		substr = "\\temp\\"
		pos1 = strings.Index(sdir, substr)
		if pos1 >= 0 {
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

// Max returns the largest of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smallest of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// Max returns the largest of x or y.
func MaxInt64(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

// Min returns the smallest of x or y.
func MinInt64(x, y int64) int64 {
	if x > y {
		return y
	}
	return x
}

// MaxDate returns the largest of x or y.
func MaxDate(x, y time.Time) time.Time {
	if x.Before(y) == true {
		return y
	}
	return x
}

// MinDate returns the smallest of x or y.
func MinDate(x, y time.Time) time.Time {
	if x.Before(y) == false {
		return y
	}
	return x
}

// GoGo - запускает функцию в отдельном потоке
func GoGo(ctx context.Context, fn func() error) error {
	var err error
	chanErr := make(chan error)

	go gogo_chan(fn, chanErr)

	select {
	case <-ctx.Done():
		Text1 := "error: TimeOut"
		err = errors.New(Text1)
		return err
	case err = <-chanErr:
		//print("err: ", err)
		break
	}

	return err
}

// gogo_chan - запускает функцию и возвращает ошибку в поток
// только совместно с GoGo()
func gogo_chan(fn func() error, chanErr chan error) {
	err := fn()
	chanErr <- err
}

// CheckInnKpp - проверяет правильность ИНН и КПП
func CheckInnKpp(Inn, Kpp string, is_individual bool) error {

	var err error

	if Inn == "" {
		Text1 := "ИНН не должен быть пустой"
		err = errors.New(Text1)
		return err
	}

	if is_individual == true {
		if len(Inn) != 12 {
			Text1 := "Длина ИНН должна быть 12 символов"
			err = errors.New(Text1)
			return err
		}
		if len(Kpp) != 0 {
			Text1 := "КПП должен быть пустой"
			err = errors.New(Text1)
			return err
		}
	} else {
		if len(Inn) != 10 {
			Text1 := "Длина ИНН должна быть 10 символов"
			err = errors.New(Text1)
			return err
		}
		if len(Kpp) != 9 {
			Text1 := "КПП должен быть 9 символов"
			err = errors.New(Text1)
			return err
		}

		err = CheckINNControlSum(Inn)
	}

	return err
}

// CheckINNControlSum - проверяет правильность ИНН по контрольной сумме
func CheckINNControlSum(Inn string) error {
	var err error

	if len(Inn) == 10 {
		err = CheckINNControlSum10(Inn)
	} else if len(Inn) == 12 {
		err = CheckINNControlSum12(Inn)
	} else {
		err = errors.New("ИНН должен быть 10 или 12 символов")
	}

	return err
}

// CheckINNControlSum10 - проверяет правильность 10-значного ИНН по контрольной сумме
func CheckINNControlSum10(Inn string) error {
	var err error

	MassKoef := [10]int{2, 4, 10, 3, 5, 9, 4, 6, 8, 0}

	var sum int
	var x int
	for i, _ := range Inn {
		s := Inn[i : i+1]
		var err1 error
		x, err1 = strconv.Atoi(s)
		if err1 != nil {
			err = errors.New("Неправильная цифра в ИНН: " + s)
			return err
		}

		sum = sum + x*MassKoef[i]
	}

	ControlSum := sum % 11
	ControlSum = ControlSum % 10
	if ControlSum != x {
		err = errors.New("Неправильная контрольная сумма ИНН")
		return err
	}

	return err
}

// CheckINNControlSum2 - проверяет правильность 12-значного ИНН по контрольной сумме
func CheckINNControlSum12(Inn string) error {
	var err error

	//контрольное чилос по 11 знакам
	MassKoef := [11]int{7, 2, 4, 10, 3, 5, 9, 4, 6, 8, 0}

	var sum int
	var x11 int
	for i := 0; i < 11; i++ {
		s := Inn[i : i+1]
		var err1 error
		x, err1 := strconv.Atoi(s)
		if err1 != nil {
			err = errors.New("Неправильная цифра в ИНН: " + s)
			return err
		}
		if i == 10 {
			x11 = x
		}

		sum = sum + x*MassKoef[i]
	}

	ControlSum := sum % 11
	ControlSum = ControlSum % 10

	if ControlSum != x11 {
		err = errors.New("Неправильная контрольная сумма ИНН")
		return err
	}

	//контрольное чилос по 12 знакам
	MassKoef2 := [12]int{3, 7, 2, 4, 10, 3, 5, 9, 4, 6, 8, 0}

	var sum2 int
	var x12 int
	for i := 0; i < 12; i++ {
		s := Inn[i : i+1]
		var err1 error
		x, err1 := strconv.Atoi(s)
		if err1 != nil {
			err = errors.New("Неправильная цифра в ИНН: " + s)
			return err
		}
		if i == 11 {
			x12 = x
		}

		sum2 = sum2 + x*MassKoef2[i]
	}

	ControlSum2 := sum2 % 11
	ControlSum2 = ControlSum2 % 10

	if ControlSum2 != x12 {
		err = errors.New("Неправильная контрольная сумма ИНН")
		return err
	}

	return err
}

// StringFromInt64 - возвращает строку из числа
func StringFromInt64(i int64) string {
	Otvet := ""

	Otvet = strconv.FormatInt(i, 10)

	return Otvet
}

// StringDate - возвращает строку дата без времени
func StringDate(t time.Time) string {
	Otvet := ""

	Otvet = t.Format("02.01.2006")

	return Otvet
}

// StringDateTime - возвращает строку дата и время, без миллисекунд
func StringDateTime(t time.Time) string {
	Otvet := ""

	Otvet = t.Format("02.01.2006 15:04:05")

	return Otvet
}

// ProgramDir_bin - возвращает каталог "bin" или каталог программы
func ProgramDir_bin() string {
	Otvet := ""

	dir := ProgramDir()
	FileName := dir + "bin" + SeparatorFile()

	ok, _ := FileExists(FileName)
	if ok == true {
		return FileName
	}

	Otvet = dir
	return Otvet
}

// SaveTempFile - записывает массив байт в файл
func SaveTempFile(bytes []byte) string {
	Otvet, err := SaveTempFile_err(bytes)
	if err != nil {
		TextError := fmt.Sprint("SaveTempFile() error: ", err)
		print(TextError)
		panic(TextError)
	}

	return Otvet
}

// SaveTempFile_err - записывает массив байт в файл, возвращает ошибку
func SaveTempFile_err(bytes []byte) (string, error) {
	Otvet := ""

	// create and open a temporary file
	f, err := os.CreateTemp("", "") // in Go version older than 1.17 you can use ioutil.TempFile
	if err != nil {
		return Otvet, err
	}

	// close and remove the temporary file at the end of the program
	defer f.Close()
	//defer os.Remove(f.Name())

	// write data to the temporary file
	if _, err := f.Write(bytes); err != nil {
		return Otvet, err
	}

	Otvet = f.Name()

	return Otvet, err
}

// Hash - возвращает число хэш из строки
func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// TextError - возвращает текст ошибки из error
func TextError(err error) string {
	Otvet := ""

	if err != nil {
		Otvet = err.Error()
	}

	return Otvet
}

// GetType - возвращает строку тип объекта
func GetType(myvar interface{}) string {
	return reflect.TypeOf(myvar).String()
}

// FindFileNameShort - возвращает имя файла(каталога) без пути
func FindFileNameShort(path string) string {
	Otvet := ""
	if path == "" {
		return Otvet
	}
	Otvet = filepath.Base(path)

	return Otvet
}

// CurrentDirectory - возвращает текущую директорию ОС
func CurrentDirectory() string {
	Otvet, err := os.Getwd()
	if err != nil {
		//log.Println(err)
	}

	return Otvet
}

// BoolFromInt64 - возвращает true если число <>0
func BoolFromInt64(i int64) bool {
	Otvet := false

	if i != 0 {
		Otvet = true
	}

	return Otvet
}

// BoolFromInt - возвращает true если число <>0
func BoolFromInt(i int) bool {
	Otvet := false

	if i != 0 {
		Otvet = true
	}

	return Otvet
}

// BoolFromString - возвращает true если строка = true, или =1
func BoolFromString(s string) bool {
	Otvet := false

	s = strings.TrimLeft(s, " ")
	s = strings.TrimRight(s, " ")
	s = strings.ToLower(s)

	if s == "true" || s == "1" {
		Otvet = true
	}

	return Otvet
}

// DeleteFileSeperator - убирает в конце / или \
func DeleteFileSeperator(dir string) string {
	Otvet := dir

	len1 := len(Otvet)
	if len1 == 0 {
		return Otvet
	}

	LastWord := Otvet[len1-1 : len1]
	if LastWord == SeparatorFile() {
		Otvet = Otvet[0 : len1-1]
	}

	return Otvet
}

// CreateFolder - создаёт папку на диске
func CreateFolder(FilenameFull string, FilePermissions uint32) error {
	var err error

	FileMode1 := os.FileMode(FilePermissions)
	if FilePermissions == 0 {
		FileMode1 = os.FileMode(0700)
	}

	if _, err := os.Stat(FilenameFull); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(FilenameFull, FileMode1)
		if err != nil {
			return err
		}
	}

	return err
}

// DeleteFolder - создаёт папку на диске
func DeleteFolder(FilenameFull string) error {
	var err error

	if _, err := os.Stat(FilenameFull); errors.Is(err, os.ErrNotExist) {
		return err
	}

	err = os.RemoveAll(FilenameFull)
	if err != nil {
		return err
	}

	return err
}

// ContextDone - возвращает true если контекст завершен
func ContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// StringFromUpperCase - возвращает строку, первая буква в верхнем регистре
func StringFromUpperCase(s string) string {
	Otvet := s
	if Otvet == "" {
		return Otvet
	}

	Otvet = strings.ToUpper(Otvet[:1]) + Otvet[1:]

	return Otvet
}

// StringFromLowerCase - возвращает строку, первая буква в нижнем регистре
func StringFromLowerCase(s string) string {
	Otvet := s
	if Otvet == "" {
		return Otvet
	}

	Otvet = strings.ToLower(Otvet[:1]) + Otvet[1:]

	return Otvet
}

// DeleteEndSlash - убирает в конце / или \
func DeleteEndSlash(Text string) string {
	Otvet := Text

	if Otvet == "" {
		return Otvet
	}

	LastSymbol := Otvet[len(Otvet)-1:]
	if LastSymbol == "/" || LastSymbol == `\` {
		Otvet = Otvet[0 : len(Otvet)-1]
	}

	return Otvet
}

// Int64FromString - возвращает int64 из строки
func Int64FromString(s string) (int64, error) {
	var Otvet int64
	var err error

	Otvet, err = strconv.ParseInt(s, 10, 64)

	return Otvet, err
}

// FindLastPos - возвращает позицию последнего вхождения
func FindLastPos(s, TextFind string) int {
	Otvet := strings.LastIndex(s, TextFind)
	return Otvet
}

// StringFloat64_Dimension2 - возвращает строку с 2 знака после запятой
func StringFloat64_Dimension2(f float64) string {
	Otvet := fmt.Sprintf("%.2f", f)

	return Otvet
}

// StringFloat32_Dimension2 - возвращает строку с 2 знака после запятой
func StringFloat32_Dimension2(f float32) string {
	Otvet := fmt.Sprintf("%.2f", f)

	return Otvet
}

// ShowTimePassed - показывает время прошедшее с момента старта
// запускать:
// defer micro.ShowTimePassed(time.Now())
func ShowTimePassed(StartAt time.Time) {
	fmt.Printf("Time passed: %s\n", time.Since(StartAt))
}

// ShowTimePassedSeconds - показывает время секунд прошедшее с момента старта
// запускать:
// defer micro.ShowTimePassedSeconds(time.Now())
func ShowTimePassedSeconds(StartAt time.Time) {
	fmt.Printf("Time passed: %s\n", time.Since(StartAt).Round(time.Second))
}

// ShowTimePassedMilliSeconds - показывает время миллисекунд прошедшее с момента старта
// запускать:
// defer micro.ShowTimePassedMilliSeconds(time.Now())
func ShowTimePassedMilliSeconds(StartAt time.Time) {
	fmt.Printf("Time passed: %s\n", time.Since(StartAt).Round(time.Millisecond))
}

// StructDeepCopy - копирует структуру из src в dist
// dist - обязательно ссылка &
func StructDeepCopy(src, dist interface{}) (err error) {
	buf := bytes.Buffer{}
	if err = gob.NewEncoder(&buf).Encode(src); err != nil {
		return
	}
	return gob.NewDecoder(&buf).Decode(dist)
}

// IsEmptyValue - возвращает true если значение по умолчанию (0, пустая строка, пустой слайс)
func IsEmptyValue(v any) bool {
	rv := reflect.ValueOf(v)
	Otvet := !rv.IsValid() || reflect.DeepEqual(rv.Interface(), reflect.Zero(rv.Type()).Interface())
	return Otvet
}

// StringIdentifierFromUUID - возвращает строку из UUID
func StringIdentifierFromUUID() string {
	Otvet := uuid.New().String()
	Otvet = strings.ReplaceAll(Otvet, "-", "")

	return Otvet
}

// IndexSubstringMin - возвращает индекс первого вхождения в строке
func IndexSubstringMin(s string, MassSubstr ...string) int {
	Otvet := -1

	for _, v := range MassSubstr {
		Otvet1 := -1
		if v != "" {
			Otvet1 = strings.Index(s, v)
		}
		if Otvet1 != -1 && (Otvet1 < Otvet || Otvet == -1) {
			Otvet = Otvet1
		}
	}

	return Otvet
}

// IndexSubstringMin2 - возвращает индекс первого вхождения в строке
func IndexSubstringMin2(s string, substr1, substr2 string) int {
	Otvet := -1

	Otvet1 := -1
	Otvet2 := -1
	if substr1 != "" {
		Otvet1 = strings.Index(s, substr1)
	}
	if substr2 != "" {
		Otvet2 = strings.Index(s, substr2)
	}

	if Otvet1 != -1 && (Otvet1 < Otvet2 || Otvet2 == -1) {
		Otvet = Otvet1
	} else {
		Otvet = Otvet2
	}

	return Otvet
}

// SortMapStringInt_Desc - сортирует map по значению, по убыванию
func SortMapStringInt_Desc(values map[string]int) []string {
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range values {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	ranked := make([]string, len(values))
	for i, kv := range ss {
		ranked[i] = kv.Key
	}
	return ranked
}

// IsNilInterface - проверка интерфейса на nil
func IsNilInterface(i any) bool {
	iv := reflect.ValueOf(i)
	if !iv.IsValid() {
		return true
	}

	switch iv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		return iv.IsNil()
	default:
		return false
	}
}

// StringFromMassInt64 - преобразование массива int64 в строку
func StringFromMassInt64(A []int64, delim string) string {

	var buffer bytes.Buffer
	for i := 0; i < len(A); i++ {
		s1 := StringFromInt64(A[i])
		buffer.WriteString(s1)
		if i != len(A)-1 {
			buffer.WriteString(delim)
		}
	}

	return buffer.String()
}

// IsInt - проверяет, является ли строка целым числом
func IsInt(s string) bool {
	Otvet := false
	if s == "" {
		return Otvet
	}

	for _, c := range s {
		if !unicode.IsDigit(c) {
			return Otvet
		}
	}

	Otvet = true
	return Otvet
}

// Int32FromString - возвращает int32 из строки
func Int32FromString(s string) (int32, error) {
	var Otvet int32
	var err error

	Otvet64, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return Otvet, err
	}

	Otvet = int32(Otvet64)

	return Otvet, err
}

// ExecuteShellCommand - выполняет команду в shell, и возвращает строку результата
func ExecuteShellCommand(TextCommand string, args ...string) (string, error) {
	Otvet := ""
	var err error

	MassByte, err := exec.Command(TextCommand, args...).CombinedOutput()
	Otvet = string(MassByte)
	if err != nil {
		return Otvet, err
	}

	return Otvet, err
}

// DeleteEndEndline - убирает в конце "\n"
func DeleteEndEndline(Text string) string {
	Otvet := Text

	if Otvet == "" {
		return Otvet
	}

	LastSymbol := Otvet[len(Otvet)-1:]
	if LastSymbol == "\n" {
		Otvet = Otvet[0 : len(Otvet)-1]
	}

	return Otvet
}

// Find_Directory_ModifiedTime - возвращает дату последнего изменения в папке internal
func Find_Directory_ModifiedTime(FolderName string) (time.Time, error) {
	var Otvet time.Time
	var err error

	dir := ProgramDir()
	dir = dir + FolderName

	ok, err := FileExists(dir)
	if err != nil {
		err = fmt.Errorf("Find_Directory_ModifiedTime() FileExists() error: %w", err)
		return Otvet, err
	}

	if ok == false {
		err = fmt.Errorf("Find_Directory_ModifiedTime() FileExists() error: file not exists: %s", dir)
		return Otvet, err
	}

	//найдём дату папки
	f, err := os.Open(dir)
	if err != nil {
		err = fmt.Errorf("Find_Directory_ModifiedTime() os.Open() error: %w", err)
		return Otvet, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		err = fmt.Errorf("Find_Directory_ModifiedTime() f.Stat() error: %w", err)
		return Otvet, err
	}

	Otvet = stat.ModTime()

	return Otvet, err
}

// Show_Repository_Code_ModifiedTime - выводит дату последнего изменения в папках cmd, internal, pkg, vendor
func Show_Repository_Code_ModifiedTime() {
	Date, err := Find_Repository_Code_ModifiedTime()
	if err != nil {
		println("Find_Repository_Code_ModifiedTime() error: ", err.Error())
		return
	}

	if Date.IsZero() {
		println("Last repository code modified time: not found")
		return
	}

	println("Last repository code modified time: ", Date.String())

}

// Find_Repository_Code_ModifiedTime - возвращает дату последнего изменения в папках cmd, internal, pkg, vendor
func Find_Repository_Code_ModifiedTime() (time.Time, error) {
	var Otvet time.Time
	var err error

	//cmd
	Time_cmd, err := Find_Directory_ModifiedTime("cmd")
	if err != nil {
		//return Otvet, err
	}

	//internal
	Time_internal, err := Find_Directory_ModifiedTime("internal")
	if err != nil {
		//return Otvet, err
	}

	//pkg
	Time_pkg, err := Find_Directory_ModifiedTime("pkg")
	if err != nil {
		//return Otvet, err
	}

	//vendor
	Time_vendor, err := Find_Directory_ModifiedTime("vendor")
	if err != nil {
		//return Otvet, err
	}

	//выбираем максимальную дату
	Otvet = TimeMax(Time_cmd, Time_internal, Time_pkg, Time_vendor)

	return Otvet, err
}

// TimeMax - возвращает максимальную дату
func TimeMax(x time.Time, y ...time.Time) time.Time {
	maxTime := x
	for _, val := range y {
		if val.After(maxTime) {
			maxTime = val
		}
	}
	return maxTime
}

// TimeMin - возвращает минимальную дату
func TimeMin(x time.Time, y ...time.Time) time.Time {
	minTime := x
	for _, val := range y {
		if val.Before(minTime) {
			minTime = val
		}
	}
	return minTime
}

// Show_Version - выводит версию сервиса на экран
func Show_Version(Version string) {
	println("Service version: ", Version)
}

// MassFrom_MapString - сортирует map по названию колонок и возвращает слайс
func MassFrom_MapString[V any](Map map[string]V) []V {
	Otvet := make([]V, 0)

	//сортировка по названию колонок
	keys := make([]string, 0, len(Map))
	for k := range Map {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	//
	for _, key1 := range keys {
		Value, ok := Map[key1]
		if ok == false {
			fmt.Printf("Map[%s] not found\n", key1)
		}
		Otvet = append(Otvet, Value)
	}

	return Otvet
}

// SortMass - сортирует слайс
func SortMass[T constraints.Ordered](s []T) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}

// MassFrom_Map - сортирует map по названию колонок и возвращает слайс
func MassFrom_Map[C constraints.Ordered, V any](Map map[C]V) []V {
	Otvet := make([]V, 0)

	//сортировка по названию колонок
	keys := make([]C, 0, len(Map))
	for k := range Map {
		keys = append(keys, k)
	}
	SortMass(keys)

	//
	for _, key1 := range keys {
		Value, ok := Map[key1]
		if ok == false {
			fmt.Printf("Map[%v] not found\n", key1)
		}
		Otvet = append(Otvet, Value)
	}

	return Otvet
}

// Substring - take at most last n characters, from start index
func Substring(input string, StartIndex int, length int) string {
	//asRunes := []rune(input)

	if StartIndex >= len(input) {
		return ""
	}

	if (StartIndex + length) >= len(input) {
		length = len(input) - StartIndex
	}

	//if StartIndex+length > len(asRunes) {
	//	length = len(asRunes) - StartIndex
	//}

	Otvet := string(input[StartIndex : StartIndex+length])
	return Otvet
}

// IntNot0 - возвращает первое ненулевое значение
func IntNot0(MassInt ...int) int {
	Otvet := 0

	for _, v := range MassInt {
		if v != 0 {
			Otvet = v
			break
		}
	}

	return Otvet
}

// InsertTextFrom - вставляет текст в середину строки
func InsertTextFrom(Text string, TextAdd string, IndexFrom int) string {
	var buffer bytes.Buffer

	//
	if IndexFrom >= len(Text) {
		return Text + TextAdd
	}

	//
	if IndexFrom < 0 {
		return TextAdd + Text
	}

	//
	s2 := SubstringLeft(Text, IndexFrom+1)
	buffer.WriteString(s2)
	buffer.WriteString(TextAdd)
	s3 := Substring(Text, IndexFrom, len(Text+TextAdd))
	buffer.WriteString(s3)

	Otvet := buffer.String()
	return Otvet
}
