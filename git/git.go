package git

import (
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/micro"
	"strings"
)

// Find_LastCommitVersion - возвращает последнюю версию в гит, образцы:
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
func Find_LastCommitVersion() (string, error) {
	Otvet := ""
	var err error

	//найдём список Хэшей коммитов
	cmd := "git"
	arg := make([]string, 0)
	arg = append(arg, "rev-list")
	arg = append(arg, "--all")
	arg = append(arg, "--max-count=1")

	output, err := micro.ExecuteShellCommand(cmd, arg...)
	if err != nil {
		return Otvet, err
	}

	if output == "" {
		return Otvet, err
	}
	MassHash0 := strings.Split(output, "\n")
	MassHash := make([]string, 0)

	for _, v := range MassHash0 {
		if v == "" {
			continue
		}
		MassHash = append(MassHash, v)
	}

	//найдём версии их Хэшей
	cmd = "git"
	arg = make([]string, 0)
	arg = append(arg, "describe")
	arg = append(arg, "--always")
	arg = append(arg, "--tags")
	arg = append(arg, MassHash...)

	Otvet, err = micro.ExecuteShellCommand(cmd, arg...)
	if err != nil {
		return Otvet, err
	}

	Otvet = micro.DeleteEndEndline(Otvet)

	return Otvet, err
}

// Show_LastCommitVersion - Выводит в консоль последнюю версию коммита в git
func Show_LastCommitVersion() {

	Text := "git commit version: "
	Otvet, err := Find_LastCommitVersion()
	if err != nil {
		Text = Text + err.Error()
	} else {
		Text = Text + Otvet
	}

	log.Info(Text)
}
