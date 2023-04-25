package docx

import (
	"github.com/manyakrus/starter/pdf_generator/internal/v0/app/programdir"
	"testing"
)

func TestCreateDocx1(t *testing.T) {

	ProgramDir := programdir.ProgramDir()
	filename := ProgramDir + "templates/test.docx"

	FilenameOut := ProgramDir + "templates/test_ready.docx"

	map1 := make(map[string]string)
	map1["{{name}}"] = "Никитин А.В."

	err := CreateDocx(filename, FilenameOut, map1)
	if err != nil {
		t.Error("docx_test.TestCreateDocx1() error: ", err)
	}

}

func TestCreateClaim(t *testing.T) {

	ProgramDir := programdir.ProgramDir()
	filename := ProgramDir + "templates/Претензия_Шаблон.docx"

	FilenameOut := ProgramDir + "templates/Претензия.docx"

	map1 := make(map[string]string)
	map1["{{filial_name}}"] = "Центральный офис"
	map1["{{filial_address}}"] = "115432, город Москва, Проектируемый 4062-й пр-д, д. 6 стр. 25"
	map1["{{filial_phone}}"] = "+7 (495) 363-13-26"
	map1["{{filial_fax}}"] = "+7 (495) 784-77-01"
	map1["{{filial_okpo}}"] = "57082325"
	map1["{{filial_ogrn}}"] = "1027700050278"
	map1["{{filial_inn}}"] = "7704228075"
	map1["{{filial_kpp}}"] = "772501001"
	map1["{{partner_name}}"] = `ООО "Ромашка"`
	map1["{{partner_address}}"] = "г.Москва ул.Ленина д.1"
	map1["{{claim_date}}"] = "01.08.2022"
	map1["{{claim_number}}"] = "111"
	map1["{{partner_fullname}}"] = `Общество с ограниченной ответственностью "Ромашка"`
	map1["{{contract_number}}"] = "222"
	map1["{{contract_date}}"] = "01.01.2022"
	map1["{{summa_dolg}}"] = "10000"
	map1["{{period}}"] = "январь 2022 г."
	map1["{{date_limit}}"] = "31.12.2022"
	map1["{{summa_electricity}}"] = "9000"
	map1["{{summa_penalty}}"] = "600"
	map1["{{summa_costs}}"] = "400"

	err := CreateDocx(filename, FilenameOut, map1)
	if err != nil {
		t.Error("docx_test.TestCreateClaim() error: ", err)
	}

}
