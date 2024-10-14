package git

import (
	_ "embed"
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/micro"
	"strconv"
	"strings"
	"time"
)

// Find_LastCommitHash - возвращает хэш последней версии в гит
func Find_LastCommitHash() (string, error) {
	Otvet := ""
	var err error

	//найдём версии их Хэшей
	cmd := "git"
	arg := make([]string, 0)
	arg = append(arg, "rev-parse")
	arg = append(arg, "HEAD")

	Otvet, err = micro.ExecuteShellCommand(cmd, arg...)
	if err != nil {
		return Otvet, err
	}

	Otvet = micro.DeleteEndEndline(Otvet)

	return Otvet, err
}

// Find_LastCommitDescribe - возвращает последнюю версию в гит, образцы:
func Find_LastCommitDescribe() (string, error) {
	Otvet := ""
	var err error

	//найдём версии их Хэшей
	cmd := "git"
	arg := make([]string, 0)
	arg = append(arg, "describe")
	arg = append(arg, "--long")
	arg = append(arg, "--tags")

	Otvet, err = micro.ExecuteShellCommand(cmd, arg...)
	if err != nil {
		return Otvet, err
	}

	Otvet = micro.DeleteEndEndline(Otvet)

	return Otvet, err
}

// Find_CommitDescribe - возвращает последнюю версию в гит, образцы:
// v1.2.159-15-ga4b0c32b
// v1.2.159-14-gafa2f9b5
// v1.2.159-13-g27f8c242
// v1.2.159-12-gc716a327
// v1.2.159-11-g1c7efca0
// v1.2.159-10-g1f369547
// v1.2.159-9-gdc9f7202
// v1.2.159-8-g6c40f58b
// v1.2.159-7-g1052bb20
// v1.2.159-6-ged68de47
// v1.2.159-5-gaf92f802
// v1.2.159-4-gb49931d5
// v1.2.159-3-g0ff81ea4
// v1.2.159-2-gc29d509e
// v1.2.159-1-g3d8ae0fd
// v1.2.159
// v1.2.158-20-g06257859
// v1.2.158-15-g1cbe3bc2
// v1.2.158-14-gcf2bce22
// v1.2.158-13-gc45d16a8
func Find_CommitDescribe(Hash string) (string, error) {
	Otvet := ""
	var err error

	//найдём версии их Хэшей
	cmd := "git"
	arg := make([]string, 0)
	arg = append(arg, "describe")
	arg = append(arg, "--long")
	arg = append(arg, "--tags")
	arg = append(arg, Hash)

	Otvet, err = micro.ExecuteShellCommand(cmd, arg...)
	if err != nil {
		return Otvet, err
	}

	Otvet = micro.DeleteEndEndline(Otvet)

	return Otvet, err
}

// Find_CommitTime - возвращает время последнего коммита
func Find_CommitTime(Hash string) (time.Time, error) {
	var Otvet time.Time
	var err error

	//найдём версии их Хэшей
	cmd := "git"
	arg := make([]string, 0)
	arg = append(arg, "show")
	arg = append(arg, "-s")
	arg = append(arg, `--format="%ct"`)
	arg = append(arg, Hash)

	//время в формате строка число unix
	sTime, err := micro.ExecuteShellCommand(cmd, arg...)
	if err != nil {
		return Otvet, err
	}
	sTime = strings.ReplaceAll(sTime, "\n", "")
	sTime = strings.ReplaceAll(sTime, `"`, "")

	//время в формате число unix
	iTime, err := micro.Int64FromString(sTime)
	if err != nil {
		return Otvet, err
	}

	//время в формате time
	Otvet = time.Unix(iTime, 0)

	return Otvet, err
}

// Find_LastCommitTime - возвращает время последнего коммита
func Find_LastCommitTime() (time.Time, error) {
	var Otvet time.Time
	var err error

	//найдём версии их Хэшей
	cmd := "git"
	arg := make([]string, 0)
	arg = append(arg, "show")
	arg = append(arg, "-s")
	arg = append(arg, `--format="%ct"`)

	//время в формате строка число unix
	sTime, err := micro.ExecuteShellCommand(cmd, arg...)
	if err != nil {
		return Otvet, err
	}
	sTime = strings.ReplaceAll(sTime, "\n", "")
	sTime = strings.ReplaceAll(sTime, `"`, "")

	//время в формате число unix
	iTime, err := micro.Int64FromString(sTime)
	if err != nil {
		return Otvet, err
	}

	//время в формате time
	Otvet = time.Unix(iTime, 0)

	return Otvet, err
}

// Find_LastCommitHashes - возвращает массив последних коммитов в гит
func Find_LastCommitHashes(count int) ([]string, error) {
	MassOtvet := make([]string, 0)
	var err error

	//найдём список Хэшей коммитов
	cmd := "git"
	arg := make([]string, 0)
	arg = append(arg, "rev-list")
	arg = append(arg, "--all")
	arg = append(arg, "--max-count="+strconv.Itoa(count))

	output, err := micro.ExecuteShellCommand(cmd, arg...)
	if err != nil {
		return MassOtvet, err
	}

	if output == "" {
		return MassOtvet, err
	}
	MassHash0 := strings.Split(output, "\n")

	//уберём пустые строки
	for _, v := range MassHash0 {
		if v == "" {
			continue
		}
		MassOtvet = append(MassOtvet, v)
	}

	return MassOtvet, err
}

// Show_LastCommitVersion - Выводит в консоль последнюю версию коммита в git
func Show_LastCommitVersion() {

	Text := "git commit version: "
	Otvet, err := Find_LastCommitDescribe()
	if err != nil {
		Text = Text + err.Error()
	} else {
		Text = Text + Otvet
	}

	log.Info(Text)
}
