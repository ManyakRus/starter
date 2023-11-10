package object_model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/vmihailenco/msgpack/v5"

	"gitlab.aescorp.ru/dsp_dev/claim/common/object_model/front/front_format"
	// "gitlab.aescorp.ru/dsp_dev/claim/common/object_model/front/front_payment"
)

// BriefCase -- набор данных для конкретного портфеля
type BriefCase struct {
	Lawsuit            Lawsuit                // Дело
	ChangeItems        []ChangeItem           `json:"change_items"`        // 3. История изменений
	Comments           []Comment              `json:"comments"`            // 4. Комментарии
	DeletedPayments    []LawsuitPayment       `json:"deleted_payments"`    // 5. Удаленные платежи
	Files              []File                 `json:"files"`               // 6. Файлы
	Hashtags           []Hashtag              `json:"hashtags"`            // 7. Хештеги портфеля
	Invoices           []LawsuitInvoice       `json:"invoices"`            // 8. Счета фактуры
	Messages           []Message              `json:"messages"`            // 9. Сообщения
	Payments           []LawsuitPayment       `json:"payments"`            // 10. Платежи
	PenaltyCalculation PenaltyCalculationList `json:"penalty_calculation"` // 15. Таблица калькуляции пени
	StateDuties        []StateDuty            `json:"state_duties"`        // 12. Гос.пошлина
	StatusStates       []LawsuitStatusState   `json:"status_states"`       // 13. История статусов дела

	// TODO Переместить сюда StatusStates
	// TODO Добавить период претензии
}

// NewBriefCase Новый объект портфеля
func NewBriefCase() BriefCase {
	return BriefCase{}
}

// AsBriefCase -- попытка распаковать байты в объект
func AsBriefCase(b []byte) (BriefCase, error) {
	c := NewBriefCase()
	err := msgpack.Unmarshal(b, &c)
	if err != nil {
		return NewBriefCase(), err
	}
	return c, nil
}

// BriefCaseAsBytes -- упаковать объект в байты
func BriefCaseAsBytes(c *BriefCase) ([]byte, error) {
	b, err := msgpack.Marshal(c)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// MailTemplateView -- возвращает параметры шаблона для письма
func MailTemplateView(bc *BriefCase) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if bc == nil {
		return result, fmt.Errorf("MailTemplateView, BriefCase is nil")
	}

	// {"cmdType":"mailing","mailing":{"head":{"protVer":"1.0.1","sysId":"RAPIRA","sourceSystem":"BS012","created":"2023-01-23 13:03:59","minProcVer":"1.0.1","senderVer":"1.0.1"},"command":"init","mailingList":[{"mailingCode":"TEST5_CLAIMDEBT20230119T132742","mailingPhaseCode":"1/1","startAt":"","timezone":"0","endAt":""}]}}
	// {"cmdType":"customerTemplateMessage","customerTemplateMessage":{"head":{"protVer":"1.0.1","sysId":"RAPIRA","sourceSystem":"BS012","created":"2023-01-23 13:03:59","minProcVer":"1.0.1","senderVer":"1.0.1"},"mailingCode":"TEST5_CLAIMDEBT20230119T132742","mailingPhaseCode":"1/1","templateCode":"CLAIMPRETRIAL","channelCode":"1401","fieldList":[{"contactInfo":"nechaevaa@atomsbt.ru","userId":"","userAddress":"183039, г. Мурманск, ул. Новое Плато, д.5, кв.60","persAcc":"5140145126","isOrganisationAcc":"1","header":"Досудебная претензия","dbtDate":"2022-09-30","dbtSum":"12 372,12","infoPhone":"9021356077","claimDate":"25.12.2022","claimPretrialDate":"28.12.2022","contractDate":"15.11.2018","contractNumber":"5140145126","dbtSumPeriod":"12 372,12","dbtSumFull":"14 460,29","lkLink":"https://lkul-murmansk.atomsbt.ru/","organisation":"","attachment":""}]}}

	dbtSumFull := bc.Lawsuit.MainSum +
		bc.Lawsuit.Penny +
		bc.Lawsuit.Penalty +
		bc.Lawsuit.StateDuty +
		bc.Lawsuit.RestrictSum

	subject := "Уведомление о задолженности"
	template := "CLAIMDEBT"
	channel := "1401" // электронная почта
	mailingCode := bc.Lawsuit.NotifyClaimMailingCode
	inBlackList := false // TODO Добавить проверку на наличие в чёрном списке
	if bc.Lawsuit.Status.Code != "1" || bc.Lawsuit.Contract.Email == "" || inBlackList {
		subject = "Досудебная претензия"
		template = "CLAIMPRETRIAL"
		if bc.Lawsuit.Contract.Email == "" || !bc.Lawsuit.Contract.IsValidEmail {
			channel = "1406" // Почта России
		}
		mailingCode = bc.Lawsuit.NotifyPretrialMailingCode

		attachments := make([]interface{}, 0)
		for i := 0; i < len(bc.Files); i++ {
			if !strings.Contains(bc.Files[i].Name, "Претензия") {
				continue
			}
			tmp := map[string]interface{}{
				"bucketName":    "claim",
				"fileName":      bc.Files[i].Name, // Тут название без расширения
				"fileExtension": bc.Files[i].Extension,
				"fileSizeByte":  bc.Files[i].Size,
				"pathToFile":    bc.Files[i].FullName,
				"eTag":          bc.Files[i].FileID,
			}
			attachments = append(attachments, tmp)

			// Получаем только первый файл, поскольку их вообще по БП не должно быть больше
			break
		}
		result["attachments"] = attachments // Вложения
	}

	// Курск - https://lkul-kursk.atomsbt.ru/
	// Мурманск – https://lkul-murmansk.atomsbt.ru/
	// Тверь - https://lkul-tver.atomsbt.ru/
	// Смоленск - https://lkul-smolensk.atomsbt.ru/
	// Хакасия - https://lkul-khakasia.atomsbt.ru/
	lkLink := bc.Lawsuit.Branch.PersonalAreaLink

	result["StageCode"] = bc.Lawsuit.Stage.Code                                 // Этап
	result["StatusCode"] = bc.Lawsuit.Status.Code                               // Статус
	result["mailingCode"] = mailingCode                                         // TODO Код рассылка "CLAIM20221004T112001" Обязательно в таком формате
	result["templateCode"] = template                                           // TODO Код шаблона (имя-строка на стороне уведомлений)
	result["channelCode"] = channel                                             // TODO Канал доставки
	result["lkLink"] = lkLink                                                   // Ссылка на Личный кабинет
	result["claimDate"] = front_format.FrontDate(bc.Lawsuit.ClaimAt)            // Дата формирования претензии+3 к.д.
	result["claimPretrialDate"] = front_format.FrontDate(bc.Lawsuit.PretrialAt) // Дата формирования досудебной претензии+5 к.д.
	result["contactInfo"] = bc.Lawsuit.Contract.Email                           // Endpoint абонента
	result["contractDate"] = front_format.FrontDate(bc.Lawsuit.Contract.SignAt) // Дата договора
	result["contractNumber"] = string(bc.Lawsuit.Contract.Number)               // Номер договора
	result["dbtDate"] = bc.Lawsuit.DateFrom.Format("2006-01-02")                // Период
	result["dbtDateStr"] = bc.Lawsuit.ClaimPeriodStr                            // Период претензии
	result["dbtSum"] = front_format.Currency(bc.Lawsuit.DebtSum)                // TODO Сумма образовавшейся задолженности
	result["dbtSumFull"] = front_format.Currency(dbtSumFull)                    // Общая сумма долга
	result["dbtSumPeriod"] = front_format.Currency(bc.Lawsuit.DebtSum)          // Сумма долга за период
	result["mailingSubject"] = subject                                          // Заголовок письма
	result["infoPhone"] = bc.Lawsuit.Contract.Organization.Phone                // Телефон абонента
	result["isOrganisationAcc"] = "1"                                           // Если организация
	//
	if strings.Trim(bc.Lawsuit.Contract.Organization.FullName, " ") == "" {
		result["organisation"] = bc.Lawsuit.Contract.Organization.Name
	} else {
		result["organisation"] = bc.Lawsuit.Contract.Organization.FullName
	}
	result["persAcc"] = string(bc.Lawsuit.Contract.Number)               // Лицевой счёт / номер договора
	result["userAddress"] = bc.Lawsuit.Contract.Organization.PostAddress // Почтовый адрес абонента

	result["russianPostType"] = "REGISTERED"
	result["russianPostAttachmentID"] = "letter"
	result["documentID"] = fmt.Sprintf("%v", bc.Lawsuit.ID)
	result["documentRegistrationDate"] = bc.Lawsuit.CreatedAt.String() // 2016-05-31T16:34:49.000+03:00
	result["documentRegistrationNumber"] = fmt.Sprintf("%v", bc.Lawsuit.Number)
	result["documentTitle"] = "Досудебная претензия"

	// Акционерное общество
	// «АтомЭнергоСбыт»
	// (АО «АтомЭнергоСбыт»)
	// Филиал «АтомЭнергоСбыт» Мурманск
	// ул. Челюскинцев, д. 30
	// г. Мурманск, 183038
	// Телефон (8152) 69-23-59
	// Факс (8152) 69-24-21
	// www.atomsbt.ru
	// e-mail: info@murmansk.atomsbt.ru
	// ОКПО 16445017, ОГРН 1027700050278
	// ИНН 7704228075, КПП 519043001

	result["senderCode"] = "CODE_55"
	result["senderDepartmentCode"] = "51"
	result["senderFullName"] = "Акционерное общество\n«АтомЭнергоСбыт»"
	result["senderName"] = "АО «АтомЭнергоСбыт»"
	result["senderDepartmentName"] = "Филиал «АтомЭнергоСбыт» Мурманск"
	result["senderAddress"] = "ул. Челюскинцев, д. 30\nг. Мурманск, 183038"
	result["senderPhone"] = "Телефон: (8152) 69-23-59"
	result["senderFax"] = "Факс: (8152) 69-24-21"
	result["senderWWW"] = "www.atomsbt.ru"
	result["senderEMail"] = "e-mail: info@murmansk.atomsbt.ru"
	result["senderOKPO"] = "ОКПО 16445017, ОГРН 1027700050278"
	result["senderINN"] = "ИНН 7704228075, КПП 519043001"

	result["senderPersonName"] = bc.Lawsuit.Contract.CuratorClaim.SecondName

	result["recipientAddress"] = bc.Lawsuit.Contract.Organization.LegalAddress
	result["recipientOrganizationName"] = bc.Lawsuit.Contract.Organization.FullName
	result["recipientOrganizationINN"] = bc.Lawsuit.Contract.Organization.INN

	return result, nil
}

// BriefCaseView -- выборка данных из претензии для отображения в веб-морде
func BriefCaseView(bc *BriefCase, c *CommonRef, t *TypeRef, useFormat bool) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if bc == nil {
		return result, fmt.Errorf("MailTemplateView, BriefCase is nil")
	}
	if c == nil {
		return result, fmt.Errorf("MailTemplateView, CommonRef is nil")
	}
	if t == nil {
		return result, fmt.Errorf("MailTemplateView, TypeRef is nil")
	}

	lawSuit := bc.Lawsuit

	// ID
	result["Lawsuit_ID"] = lawSuit.ID

	// Дата претензии
	if useFormat {
		result["Lawsuit_CreatedAt"] = formatDate(lawSuit.CreatedAt)
		result["Lawsuit_ClaimAt"] = formatDate(lawSuit.ClaimAt)
		result["Lawsuit_PretrialAt"] = formatDate(lawSuit.PretrialAt)
	} else {
		result["Lawsuit_CreatedAt"] = lawSuit.CreatedAt
		result["Lawsuit_ClaimAt"] = lawSuit.ClaimAt
		result["Lawsuit_PretrialAt"] = lawSuit.PretrialAt
	}

	// Номер претензии
	result["Lawsuit_Number"] = string(lawSuit.Number)
	result["Lawsuit_NumberClaim"] = string(lawSuit.NumberClaim)
	result["Lawsuit_NumberTrial"] = lawSuit.NumberTrial

	// Строковый период претензии
	result["Lawsuit_Period"] = bc.Lawsuit.ClaimPeriodStr

	// TODO View LawsuitStageTypes проверить чтение из готового объекта, без перебора соответствий
	stage := LawsuitStageType{}
	for i := 0; i < len(t.LawsuitStageTypes); i++ {
		if t.LawsuitStageTypes[i].ID == lawSuit.StageID {
			stage = t.LawsuitStageTypes[i]
			break
		}
	}
	// Этап - для фильтрации
	result["Lawsuit_StageID"] = stage.ID // lawSuit.Stage.ID
	// Этап - для вариантов отображения
	result["Lawsuit_StageCode"] = stage.Code // lawSuit.Stage.Code
	// Этап - для вывода в таблицу
	result["Lawsuit_Stage"] = stage.Name // lawSuit.Stage.Name
	// Дата установки этапа
	if useFormat {
		result["Lawsuit_StageAt"] = formatDate(lawSuit.StageAt)
	} else {
		result["Lawsuit_StageAt"] = lawSuit.StageAt
	}

	// TODO View LawsuitStatusTypes проверить чтение из готового объекта, без перебора соответствий
	status := LawsuitStatusType{}
	for i := 0; i < len(t.LawsuitStatusTypes); i++ {
		if t.LawsuitStatusTypes[i].ID == lawSuit.StatusID {
			status = t.LawsuitStatusTypes[i]
			break
		}
	}
	// Статус - для фильтрации
	result["Lawsuit_StatusID"] = status.ID // lawSuit.Status.ID
	// Статус - для вариантов отображения
	result["Lawsuit_StatusCode"] = status.Code // lawSuit.Status.Code
	// Статус - для вывода в таблицу
	result["Lawsuit_Status"] = status.Name // lawSuit.Status.Name
	// Дата установки статуса
	if useFormat {
		result["Lawsuit_StatusAt"] = formatDate(lawSuit.StatusAt)
	} else {
		result["Lawsuit_StatusAt"] = lawSuit.StatusAt
	}

	// TODO View LawsuitReasonTypes проверить чтение из готового объекта, без перебора соответствий
	reason := LawsuitReasonType{}
	for i := 0; i < len(t.LawsuitReasonTypes); i++ {
		if t.LawsuitReasonTypes[i].ID == lawSuit.ReasonID {
			reason = t.LawsuitReasonTypes[i]
			break
		}
	}
	// Статус - для фильтрации
	result["Lawsuit_ReasonID"] = reason.ID // lawSuit.Reason.ID
	// Статус - для вывода в таблицу
	result["Lawsuit_Reason"] = reason.Name // lawSuit.Reason.Name

	// TODO View ClaimTypes проверить чтение из готового объекта, без перебора соответствий
	claimType := ClaimType{}
	for i := 0; i < len(t.ClaimTypes); i++ {
		if t.ClaimTypes[i].ID == lawSuit.ClaimTypeID {
			claimType = t.ClaimTypes[i]
			break
		}
	}
	// Статус - для фильтрации
	result["Lawsuit_ClaimTypeID"] = claimType.ID // lawSuit.ClaimType.ID
	// Статус - для вывода в таблицу
	result["Lawsuit_ClaimType"] = claimType.Name // lawSuit.ClaimType.Name

	// TODO View LawsuitTypes проверить чтение из готового объекта, без перебора соответствий
	lawsuitType := LawsuitType{}
	for i := 0; i < len(t.LawsuitTypes); i++ {
		if t.LawsuitTypes[i].ID == lawSuit.TypeID {
			lawsuitType = t.LawsuitTypes[i]
			break
		}
	}
	// Тип претензии - для фильтрации
	result["Lawsuit_LawsuitTypeID"] = lawsuitType.ID // lawSuit.Type.ID
	// Тип претензии - для вывода в таблицу
	result["Lawsuit_LawsuitType"] = lawsuitType.Name // lawSuit.Type.Name

	// TODO View Branches проверить чтение из готового объекта, без перебора соответствий
	branch := Branch{}
	for i := 0; i < len(c.Branches); i++ {
		if c.Branches[i].ID == lawSuit.BranchID {
			branch = c.Branches[i]
			break
		}
	}
	// Отделение - для фильтрации
	result["Lawsuit_BranchID"] = branch.ID // lawSuit.Branch.ID
	// Отделение - для вывода в таблицу
	result["Lawsuit_Branch"] = branch.Name // lawSuit.Branch.Name

	// "Lawsuit_Balance" = 319.9                DebtSum
	// "Lawsuit_Claim" = -                      DebtSumSentNotify
	// "Lawsuit_Main" = 444743.48999999993      MainSum
	// "Lawsuit_MainDebt" = 319.9               InvoiceSum
	// "Lawsuit_Penalty" = 0                    Penalty
	// "Lawsuit_Penny" = 1364.4099999999999     Penny
	// "Lawsuit_Percent_317" = 0                Percent317
	// "Lawsuit_Percent_395" = 0                Percent395
	// "Lawsuit_Pretrial" = 319.9               DebtSumSentClaim
	// "Lawsuit_ReceivedFunds" = 0              PaySum
	// "Lawsuit_Restrict" = 0                   RestrictSum
	// "Lawsuit_StateDuty" = 2000               StateDuty
	// "Lawsuit_TotalDebt" = 448107.8999999999  dbtSumFull
	// "Lawsuit_UnknownPayments" = 0

	dbtSumFull := 0.0
	switch lawsuitType.Code {
	case "1":
		dbtSumFull = lawSuit.DebtSum + lawSuit.Penny + lawSuit.RestrictSum + lawSuit.StateDuty
	case "2":
		dbtSumFull = lawSuit.MainSum + lawSuit.DebtSum + lawSuit.RestrictSum + lawSuit.StateDuty
	case "3":
		dbtSumFull = lawSuit.MainSum + lawSuit.Penny + lawSuit.DebtSum + lawSuit.StateDuty
	}

	// Сумма по основному виду деятельности (руб.)
	if lawsuitType.Code == "1" {
		if useFormat {
			result["Lawsuit_Main"] = Currency(lawSuit.DebtSum)
		} else {
			result["Lawsuit_Main"] = lawSuit.DebtSum
		}
	} else {
		if useFormat {
			result["Lawsuit_Main"] = Currency(lawSuit.MainSum)
		} else {
			result["Lawsuit_Main"] = lawSuit.MainSum
		}
	}

	// Пени по день фактической оплаты долга (руб.)
	if lawsuitType.Code == "2" {
		if useFormat {
			result["Lawsuit_Penny"] = Currency(lawSuit.DebtSum)
		} else {
			result["Lawsuit_Penny"] = lawSuit.DebtSum
		}
	} else {
		if useFormat {
			result["Lawsuit_Penny"] = Currency(lawSuit.Penny)
		} else {
			result["Lawsuit_Penny"] = lawSuit.Penny
		}
	}

	// Сумма ограничений (руб.)
	if lawsuitType.Code == "3" {
		if useFormat {
			result["Lawsuit_Restrict"] = Currency(lawSuit.DebtSum)
		} else {
			result["Lawsuit_Restrict"] = lawSuit.DebtSum
		}
	} else {
		if useFormat {
			result["Lawsuit_Restrict"] = Currency(lawSuit.RestrictSum)
		} else {
			result["Lawsuit_Restrict"] = lawSuit.RestrictSum
		}
	}

	// Сумма процентов по 395 ГК РФ (руб.)
	if useFormat {
		result["Lawsuit_Percent_395"] = Currency(bc.Lawsuit.Percent395)
	} else {
		result["Lawsuit_Percent_395"] = bc.Lawsuit.Percent395
	}
	// Сумма процентов по 317.1 ГК РФ (руб.)
	if useFormat {
		result["Lawsuit_Percent_317"] = Currency(bc.Lawsuit.Percent317)
	} else {
		result["Lawsuit_Percent_317"] = bc.Lawsuit.Percent317
	}
	// Сумма договорной/законной неустойки (руб.)
	if useFormat {
		result["Lawsuit_Penalty"] = Currency(bc.Lawsuit.Penalty)
	} else {
		result["Lawsuit_Penalty"] = bc.Lawsuit.Penalty
	}
	// Сумма госпошлины (руб.)
	if useFormat {
		result["Lawsuit_StateDuty"] = Currency(bc.Lawsuit.StateDuty)
	} else {
		result["Lawsuit_StateDuty"] = bc.Lawsuit.StateDuty
	}
	// Поступило денежных средств
	if useFormat {
		result["Lawsuit_ReceivedFunds"] = Currency(lawSuit.PaySum)
	} else {
		result["Lawsuit_ReceivedFunds"] = lawSuit.PaySum
	}
	// TODO Поле "Общий долг": Полная сумма долга
	if useFormat {
		result["Lawsuit_TotalDebt"] = Currency(dbtSumFull)
	} else {
		result["Lawsuit_TotalDebt"] = dbtSumFull
	}
	// TODO Поле "Основной долг": По счёт фактурам
	if useFormat {
		result["Lawsuit_MainDebt"] = Currency(lawSuit.InvoiceSum)
	} else {
		result["Lawsuit_MainDebt"] = lawSuit.InvoiceSum
	}
	// TODO Поле "Остаток долга": "Основной долг" - Поступило денежных средств
	if useFormat {
		result["Lawsuit_Balance"] = Currency(lawSuit.DebtSum)
	} else {
		result["Lawsuit_Balance"] = lawSuit.DebtSum
	}

	if lawSuit.UnknownPayments {
		result["Lawsuit_UnknownPayments"] = 1
	} else {
		result["Lawsuit_UnknownPayments"] = 0
	}

	changes := make([]interface{}, 0)
	for i := 0; i < len(bc.ChangeItems); i++ {
		// TODO Костыль, нужно по-быстрому решить, чтобы не парсить в веб
		actionCode := 0
		newValueCode := 0
		prevValueCode := 0
		if bc.ChangeItems[i].Key == "Обновление статуса" {
			actionCode = 1
			// "Сформировано уведомление (2)"
			value := regexp.MustCompile(`\d`).FindStringSubmatch(bc.ChangeItems[i].Value)
			if len(value) == 1 {
				newValueCode, _ = strconv.Atoi(value[0])
			}
		}
		tmp := map[string]interface{}{
			"ID":            bc.ChangeItems[i].ID,
			"CreatedAt":     formatTime(bc.ChangeItems[i].CreatedAt),
			"Action":        bc.ChangeItems[i].Key,
			"NewValue":      bc.ChangeItems[i].Value,
			"PrevValue":     bc.ChangeItems[i].Prev,
			"ActionCode":    actionCode,
			"NewValueCode":  newValueCode,
			"PrevValueCode": prevValueCode,
		}
		changes = append(changes, tmp)
	}
	result["Lawsuit_Changes"] = changes

	contract := &lawSuit.Contract
	// ID
	result["Contract_ID"] = contract.ID
	// № Договор
	result["Contract_Number"] = string(contract.Number)
	// Дата договора
	if useFormat {
		result["Contract_SignAt"] = formatDate(contract.SignAt)
	} else {
		result["Contract_SignAt"] = contract.SignAt
	}
	// Категория договора
	category := "Неизвестно"
	for i := 0; i < len(t.ContractCategoryTypes); i++ {
		if t.ContractCategoryTypes[i].ID == contract.CategoryID {
			category = t.ContractCategoryTypes[i].Name
			break
		}
	}
	result["Contract_Category"] = category
	result["Contract_CategoryID"] = contract.CategoryID
	// Статус
	if contract.Status == "" {
		result["Contract_Status"] = "Активен"
	} else {
		result["Contract_Status"] = contract.Status
	}
	// Почтовый адрес
	if contract.PostAddress == "" {
		result["Contract_PostAddress"] = "Не указан"
	} else {
		result["Contract_PostAddress"] = contract.PostAddress
	}
	// E-mail
	if contract.Email == "" {
		result["Contract_Email"] = "Не указан"
	} else {
		result["Contract_Email"] = contract.Email
	}
	// Валидность E-mail
	if contract.IsValidEmail {
		result["Contract_EmailValid"] = 1
	} else {
		result["Contract_EmailValid"] = 0
	}

	// Судебный отдел
	result["Contract_CuratorLegal"] = fmt.Sprintf("%v %v", contract.CuratorLegal.SecondName, contract.CuratorLegal.Name)
	result["Contract_CuratorLegal_Email"] = contract.CuratorLegal.Email
	result["Contract_CuratorLegal_Phone"] = contract.CuratorLegal.Phone
	// Расчётный отдел
	result["Contract_CuratorPayment"] = fmt.Sprintf("%v %v", contract.CuratorPayment.SecondName, contract.CuratorPayment.Name)
	result["Contract_CuratorPayment_Email"] = contract.CuratorPayment.Email
	result["Contract_CuratorPayment_Phone"] = contract.CuratorPayment.Phone
	// Договорной отдел
	result["Contract_CuratorContract"] = fmt.Sprintf("%v %v", contract.CuratorContract.SecondName, contract.CuratorContract.Name)
	result["Contract_CuratorContract_Email"] = contract.CuratorContract.Email
	result["Contract_CuratorContract_Phone"] = contract.CuratorContract.Phone
	// Отдел тех. аудита
	result["Contract_CuratorTechAudit"] = fmt.Sprintf("%v %v", contract.CuratorTechAudit.SecondName, contract.CuratorTechAudit.Name)
	result["Contract_CuratorTechAudit_Email"] = contract.CuratorTechAudit.Email
	result["Contract_CuratorTechAudit_Phone"] = contract.CuratorTechAudit.Phone
	// Куратор претензии
	result["Contract_CuratorClaim"] = fmt.Sprintf("%v %v", contract.CuratorClaim.SecondName, contract.CuratorClaim.Name)
	result["Contract_CuratorClaim_Email"] = contract.CuratorClaim.Email
	result["Contract_CuratorClaim_Phone"] = contract.CuratorClaim.Phone
	result["Contract_CuratorClaim_FullName"] = fmt.Sprintf("%v %v %v", contract.CuratorClaim.SecondName, contract.CuratorClaim.Name, contract.CuratorClaim.ParentName)

	result["Contract_DaysToResolveClaim"] = contract.DaysToResolveClaim

	result["Contract_PaymentDay"] = 18
	for i := 0; i < len(contract.PaymentDays); i++ {
		if time.Now().After(contract.PaymentDays[i].DateFrom) &&
			time.Now().Before(contract.PaymentDays[i].DateTo) {
			result["Contract_PaymentDay"] = contract.PaymentDays[i].Day
			break
		}
	}

	paymentSchedules := make([]interface{}, 0)
	for i := 0; i < len(contract.PaymentSchedules); i++ {
		if contract.PaymentSchedules[i].ContractID == contract.ID {
			if useFormat {
				tmp := map[string]interface{}{
					"Day":     fmt.Sprintf("%d число", contract.PaymentSchedules[i].Day),
					"Percent": fmt.Sprintf("%d %%", contract.PaymentSchedules[i].Percent),
				}
				paymentSchedules = append(paymentSchedules, tmp)
			} else {
				tmp := map[string]interface{}{
					"Day":     fmt.Sprintf("%d", contract.PaymentSchedules[i].Day),
					"Percent": fmt.Sprintf("%d", contract.PaymentSchedules[i].Percent),
				}
				paymentSchedules = append(paymentSchedules, tmp)
			}
		}
	}
	if len(paymentSchedules) == 0 {
		tmp := map[string]interface{}{
			"Day":     "18",
			"Percent": "100 %",
		}
		paymentSchedules = append(paymentSchedules, tmp)
	}
	// Срок оплаты по договору
	result["Contract_PaymentSchedules"] = paymentSchedules

	invoices := make([]interface{}, 0)
	totalSum := 0.0
	totalCorrectionSum := 0.0
	totalDebtSum := 0.0
	totalPayment := 0.0
	totalBalance := 0.0
	for i := 0; i < len(bc.Invoices); i++ {
		// TODO Пока аналитики вообще скрываю
		// if strings.Trim(bc.Invoices[i].Document.Analytics, " ") != "" {
		// 	continue
		// }

		paymentSum := 0.0
		// correctionType1Sum := 0.0
		// correctionType2Sum := 0.0
		correctionSum := 0.0
		for j := 0; j < len(bc.Payments); j++ {
			if bc.Invoices[i].ID == int64(bc.Payments[j].InvoiceID) {
				if bc.Payments[j].IsCorrective {
					correctionSum += bc.Payments[j].Sum * -1
					// correctionType1Sum += bc.Payments[j].Sum
				} else if !bc.Payments[j].IsCorrective && bc.Payments[j].Document.DocumentTypeID == 35 {
					// correctionType2Sum += bc.Payments[j].Sum
					correctionSum += bc.Payments[j].Sum
				} else {
					paymentSum += bc.Payments[j].Sum
				}
			}
		}

		note := bc.Invoices[i].Document.Note
		if note == "" {
			note = "Не задано"
		}

		prefix := "СФ"
		// Пенни
		if bc.Lawsuit.Type.Code == "2" {
			prefix = "C"
		}
		number := fmt.Sprintf("%v:%v", prefix, bc.Invoices[i].Document.Number)
		numberFull := bc.Invoices[i].Document.NumberFull
		sum := Currency(bc.Invoices[i].Sum)
		if bc.Invoices[i].IsCorrective {
			number = "К" + number
			// sum = Currency(0.0)
		}

		// currentAfterCorrectionSum := bc.Invoices[i].Sum - correctionType1Sum + correctionType2Sum
		currentAfterCorrectionSum := bc.Invoices[i].Sum + correctionSum
		currentDebtSum := currentAfterCorrectionSum - paymentSum
		// correctionBothTypesSum := correctionType1Sum + correctionType2Sum

		tmp := map[string]interface{}{
			"ID":               bc.Invoices[i].ID,
			"ClaimNumber":      lawSuit.NumberClaim,                            // Поле "Претензия"
			"Date":             formatDate(bc.Invoices[i].Document.DocumentAt), // Поле "Дата С/Ф"
			"DistributionDate": formatTime(bc.Invoices[i].Document.CreatedAt),  // Поле "Дата разнесения"
			"Number":           number,                                         // Поле "Номер С/Ф"
			"NumberFull":       numberFull,                                     // Поле "Номер С/Ф" полный
			"Type":             bc.Invoices[i].Document.Analytics,              // Поле "Тип начисления"
			"Count":            bc.Invoices[i].Count,                           // Кол-во кВт
			"Sum":              sum,                                            // Поле "Начислено"
			"Correction":       Currency(correctionSum),                        // Поле "Корректировка"
			"DebtSum":          Currency(currentDebtSum),                       // Поле "Долг в претензии"
			"Payment":          Currency(paymentSum),                           // Поле "Оплачено"
			"Balance":          Currency(currentAfterCorrectionSum),            // Поле "Остаток"
			"Note":             note,                                           // Поле "Примечание"
		}

		// TODO Аналитики нужно считать отдельно
		// if strings.Trim(bc.Invoices[i].Document.Analytics, " ") == "" {
		totalSum += bc.Invoices[i].Sum
		totalCorrectionSum += correctionSum
		totalPayment += paymentSum

		totalDebtSum += currentDebtSum
		totalBalance += currentAfterCorrectionSum
		// }

		invoices = append(invoices, tmp)
	}
	totalDebtSumState := fmt.Sprintf("%.2f", totalDebtSum) >= fmt.Sprintf("%.2f", 0.0)
	// Счета фактуры по данному договору
	result["Contract_Invoices"] = invoices
	// Суммы счетов фактур по данному договору
	result["Contract_TotalInvoices"] = map[string]interface{}{
		"Sum":          Currency(totalSum),           // Поле "Начислено"
		"CorSum":       Currency(totalCorrectionSum), // Поле "Корректировка"
		"DebtSum":      Currency(totalDebtSum),       // Поле "Долг в претензии"
		"DebtSumState": totalDebtSumState,            // Поля признак для регулирования состояния Долг/Переплата на фронте
		"Payment":      Currency(totalPayment),       // Поле "Оплачено"
		"Balance":      Currency(totalBalance),       // Поле "Остаток"
	}

	payments := make([]interface{}, 0)
	totalSum = 0.0
	totalCorrectionSum = 0.0
	totalDebtSum = 0.0
	totalPayment = 0.0
	totalBalance = 0.0
	totalUnknownPayment := 0.0
	// totalPaymentsBeforeClaim := 0.0         // сумма всех платежей с момента выставления с/ф
	// totalPaymentsAllClaim := 0.0            // сумма всех платежей с момента формирования претензии
	for i := 0; i < len(bc.Payments); i++ { // Перебираем платежи в конкретной претензии
		// TODO Пока аналитики вообще скрываю
		// if strings.Trim(bc.Payments[i].Document.Analytics, " ") != "" {
		// 	continue
		// }

		note := bc.Payments[i].Document.Note
		if note == "" {
			note = "Не задано"
		}

		number := "ПП:" + bc.Payments[i].Document.Number
		isInvoiceOut := false // если текущий документ с/ф
		correction := ""
		payment := Currency(bc.Payments[i].Sum)
		if bc.Payments[i].Document.DocumentTypeID == 35 {
			prefix := "СФ"
			// Пенни
			if bc.Lawsuit.Type.Code == "2" {
				prefix = "C"
			}
			number := fmt.Sprintf("%v:%v", prefix, bc.Payments[i].Document.Number)
			correction = Currency(bc.Payments[i].Sum)
			payment = ""
			isInvoiceOut = true
			if !bc.Payments[i].IsCorrective {
				number = "И" + number
			}
		}
		if bc.Payments[i].IsCorrective {
			number = "К" + number
			correction = Currency(bc.Payments[i].Sum * -1)
		}

		// paymentDoc := bc.Payments[i].Sum

		isPaymentAfterCreated := true

		invoiceSum := 0.0
		for j := 0; j < len(bc.Invoices); j++ {
			if int64(bc.Payments[i].InvoiceID) == bc.Invoices[j].ID {
				invoiceSum += bc.Invoices[j].Sum
				if !isInvoiceOut {
					controlDate := lawSuit.CreatedAt
					strControlDate := controlDate.Local().Format("2006-01-02 15:04:05")
					docDate := bc.Payments[i].Document.CreatedAt
					strDocDate := docDate.Local().Format("2006-01-02 15:04:05")
					if strControlDate >= strDocDate {
						// totalPaymentsBeforeClaim += paymentDoc // Считаем сумму платежей до формирования претензии
						isPaymentAfterCreated = false
					}
					// totalPaymentsAllClaim += paymentDoc // Считаем сумму всех платежей
				}
				break
			}
		}

		tmp := map[string]interface{}{
			"ID":                    bc.Payments[i].ID,
			"InvoiceID":             bc.Payments[i].InvoiceID,                       // Ссылка на С/Ф
			"ClaimNumber":           lawSuit.NumberClaim,                            // Поле "Претензия"
			"Date":                  formatDate(bc.Payments[i].Document.DocumentAt), // Поле "Дата С/Ф"
			"DistributionDate":      formatTime(bc.Payments[i].Document.CreatedAt),  // Поле "Дата разнесения"
			"isPaymentAfterCreated": isPaymentAfterCreated,                          // Поле для подкраски в голубой
			"Number":                number,                                         // Поле "Номер С/Ф"
			"Type":                  bc.Payments[i].Document.Analytics,              // Поле "Тип начисления"
			"Sum":                   "",                                             // Поле "Начислено"
			"Correction":            correction,                                     // Поле "Корректировка"
			"DebtSum":               "",                                             // Поле "Долг в претензии"
			"Payment":               payment,                                        // Поле "Оплачено"
			"Balance":               "",                                             // Поле "Остаток"
			"Note":                  note,                                           // Поле "Примечание"
		}

		// TODO Аналитики нужно считать отдельно
		// if strings.Trim(bc.Payments[i].Document.Analytics, " ") == "" {
		//  Платежи без фактур
		if bc.Payments[i].InvoiceID == 0 {
			totalUnknownPayment += bc.Payments[i].Sum
		} else {
			totalDebtSum += invoiceSum - bc.Payments[i].Sum
			totalPayment += bc.Payments[i].Sum
			totalBalance += invoiceSum - bc.Payments[i].Sum

			if bc.Payments[i].IsCorrective {
				totalCorrectionSum += bc.Payments[i].Sum
			} else {
				totalSum += invoiceSum
			}
		}
		// }

		payments = append(payments, tmp)
	}

	deletedPayments := make([]interface{}, 0)
	for i := 0; i < len(bc.DeletedPayments); i++ { // Перебираем платежи в конкретной претензии
		// TODO Пока аналитики вообще скрываю
		// if strings.Trim(bc.DeletedPayments[i].Document.Analytics, " ") != "" {
		// 	continue
		// }

		note := bc.DeletedPayments[i].Document.Note
		if note == "" {
			note = "Не задано"
		}

		number := "ПП:" + bc.DeletedPayments[i].Document.Number
		isInvoiceOut := false // если текущий документ с/ф
		correction := ""
		payment := Currency(bc.DeletedPayments[i].Sum)
		if bc.DeletedPayments[i].Document.DocumentTypeID == 35 {
			prefix := "СФ"
			// Пенни
			if bc.Lawsuit.Type.Code == "2" {
				prefix = "C"
			}
			number := fmt.Sprintf("%v:%v", prefix, bc.DeletedPayments[i].Document.Number)
			correction = Currency(bc.DeletedPayments[i].Sum)
			payment = ""
			isInvoiceOut = true
			if !bc.DeletedPayments[i].IsCorrective {
				number = "И" + number
			}
		}
		if bc.DeletedPayments[i].IsCorrective {
			number = "К" + number
			correction = Currency(bc.DeletedPayments[i].Sum * -1)
		}

		isPaymentAfterCreated := true

		invoiceSum := 0.0
		for j := 0; j < len(bc.Invoices); j++ {
			if int64(bc.DeletedPayments[i].InvoiceID) == bc.Invoices[j].ID {
				invoiceSum += bc.Invoices[j].Sum
				if !isInvoiceOut {
					controlDate := lawSuit.CreatedAt
					strControlDate := controlDate.Local().Format("2006-01-02 15:04:05")
					docDate := bc.DeletedPayments[i].Document.CreatedAt
					strDocDate := docDate.Local().Format("2006-01-02 15:04:05")
					if strControlDate >= strDocDate {
						isPaymentAfterCreated = false
					}
				}
				break
			}
		}

		tmp := map[string]interface{}{
			"ID":                    bc.DeletedPayments[i].ID,
			"InvoiceID":             bc.DeletedPayments[i].InvoiceID,                       // Ссылка на С/Ф
			"ClaimNumber":           lawSuit.NumberClaim,                                   // Поле "Претензия"
			"Date":                  formatDate(bc.DeletedPayments[i].Document.DocumentAt), // Поле "Дата С/Ф"
			"DistributionDate":      formatTime(bc.DeletedPayments[i].Document.CreatedAt),  // Поле "Дата разнесения"
			"isPaymentAfterCreated": isPaymentAfterCreated,                                 // Поле для подкраски в голубой
			"Number":                number,                                                // Поле "Номер С/Ф"
			"Type":                  bc.DeletedPayments[i].Document.Analytics,              // Поле "Тип начисления"
			"Sum":                   "",                                                    // Поле "Начислено"
			"Correction":            correction,                                            // Поле "Корректировка"
			"DebtSum":               "",                                                    // Поле "Долг в претензии"
			"Payment":               payment,                                               // Поле "Оплачено"
			"Balance":               "",                                                    // Поле "Остаток"
			"Note":                  note,                                                  // Поле "Примечание"
			"DeletedDate":           formatTime(bc.DeletedPayments[i].DeletedAt),           // Поле "Дата удаления"
		}

		deletedPayments = append(deletedPayments, tmp)
	}

	DebtSumOnCreate := lawSuit.GetStatusState(1).TotalDebt
	DebtSumSentNotify := lawSuit.GetStatusState(2).TotalDebt
	DebtSumSentClaim := lawSuit.GetStatusState(4).TotalDebt

	if useFormat {
		result["Lawsuit_DebtSumOnCreate"] = Currency(DebtSumOnCreate)
	} else {
		result["Lawsuit_DebtSumOnCreate"] = DebtSumOnCreate
	}
	if useFormat {
		result["Lawsuit_Claim"] = Currency(DebtSumSentNotify)
	} else {
		result["Lawsuit_Claim"] = DebtSumSentNotify
	}
	if !lawSuit.NotifyClaimDone {
		result["Lawsuit_Claim"] = "-"
	}
	if useFormat {
		result["Lawsuit_Pretrial"] = Currency(DebtSumSentClaim)
	} else {
		result["Lawsuit_Pretrial"] = DebtSumSentClaim
	}
	if !lawSuit.NotifyPretrialDone {
		result["Lawsuit_Pretrial"] = "-"
	}

	// TODO Костыль! Для PDF сервиса
	// if lawSuit.Status.Code == "1" || lawSuit.Status.Code == "3" {
	// 	// Формирование претензии вовремя
	// 	result["sum_total"] = result["Lawsuit_TotalDebt"]   // общий размер задолженности
	// 	result["sum_electricity"] = result["Lawsuit_Main"]  // за потребленную электрическую энергию
	// 	result["sum_penalty"] = result["Lawsuit_Penny"]     // неустойка/пени
	// 	result["sum_cost"] = result["Lawsuit_StateDuty"]    // государственная пошлина
	// 	result["sum_restrict"] = result["Lawsuit_Restrict"] // ограничения
	// } else {
	// 	// Формирование претензии в произвольное время
	// 	ss := lawSuit.GetStatusState(2)
	// 	if useFormat {
	// 		result["sum_total"] = Currency(ss.TotalDebt)      // общий размер задолженности
	// 		result["sum_electricity"] = Currency(ss.MainSum)  // за потребленную электрическую энергию
	// 		result["sum_penalty"] = Currency(ss.PennySum)     // неустойка/пени
	// 		result["sum_cost"] = Currency(lawSuit.StateDuty)  // государственная пошлина
	// 		result["sum_restrict"] = Currency(ss.RestrictSum) // ограничения
	// 	} else {
	// 		result["sum_total"] = ss.TotalDebt      // общий размер задолженности
	// 		result["sum_electricity"] = ss.MainSum  // за потребленную электрическую энергию
	// 		result["sum_penalty"] = ss.PennySum     // неустойка/пени
	// 		result["sum_cost"] = lawSuit.StateDuty  // государственная пошлина
	// 		result["sum_restrict"] = ss.RestrictSum // ограничения
	// 	}
	// }

	// Платёжные документы по данному договору
	result["Contract_Payments"] = payments
	// Удаленные платёжные документы по данному договору
	result["Contract_DeletedPayments"] = deletedPayments
	// Суммы платёжных документов по данному договору
	result["Contract_TotalPayments"] = map[string]interface{}{
		"Sum":     Currency(totalSum),           // Поле "Начислено"
		"CorSum":  Currency(totalCorrectionSum), // Поле "Корректировка"
		"DebtSum": Currency(totalDebtSum),       // Поле "Долг в претензии"
		"Payment": Currency(totalPayment),       // Поле "Оплачено"
		"Balance": Currency(totalBalance),       // Поле "Остаток"
	}
	// Суммы нераспознанных платёжных документов по данному договору
	result["Contract_TotalUnknownPayments"] = map[string]interface{}{
		"Payment": Currency(totalUnknownPayment), // Поле "Сумма"
	}

	organization := &contract.Organization
	// ID
	result["Organization_ID"] = organization.ID
	// Email ЮЛ
	if organization.Email == "" {
		result["Organization_Email"] = "Не указан"
	} else {
		result["Organization_Email"] = organization.Email
	}
	// Наименование ЮЛ
	result["Organization_Name"] = organization.Name
	result["Organization_FullName"] = organization.FullName
	// ИНН ЮЛ
	result["Organization_INN"] = organization.INN
	// КПП ЮЛ
	result["Organization_KPP"] = organization.KPP
	// Юридический адрес ЮЛ
	result["Organization_LegalAddress"] = organization.LegalAddress
	// Категория организации
	category = "Неизвестно"
	for i := 0; i < len(t.OrganizationCategoryTypes); i++ {
		if t.OrganizationCategoryTypes[i].ID == organization.CategoryID {
			category = t.OrganizationCategoryTypes[i].Name
			break
		}
	}
	result["Organization_CategoryID"] = organization.CategoryID
	result["Organization_Category"] = category
	// Состояние организации
	state := "Действующее"
	code := "1"
	color := "green"
	for i := 0; i < len(t.OrganizationStateTypes); i++ {
		if t.OrganizationStateTypes[i].ID == organization.StateID {
			state = t.OrganizationStateTypes[i].Name
			code = t.OrganizationStateTypes[i].Code
			color = t.OrganizationStateTypes[i].Color
			break
		}
	}
	result["Organization_State"] = state
	result["Organization_StateCode"] = code
	result["Organization_StateColor"] = color
	result["Organization_StateID"] = organization.StateID

	// Ликвидность организации - deprecated
	// if organization.IsLiquidated {
	// 	result["Organization_Liquidity"] = "Ликвидирован"
	// } else {
	// 	result["Organization_Liquidity"] = "Действующий"
	// }
	// Банкротство организации - deprecated
	// if organization.IsBankrupt {
	// 	result["Organization_Bankrupt"] = "Банкрот"
	// } else {
	// 	result["Organization_Bankrupt"] = "Действующий"
	// }

	FileMail := ""
	FileMailName := ""
	FileClaim := ""
	FileClaimName := ""
	FileClaimDetail := ""
	FileClaimDetailName := ""
	for i := 0; i < len(bc.Files); i++ {
		if strings.Contains(bc.Files[i].Name, "Письмо") {
			FileMail = bc.Files[i].FileID
			FileMailName = bc.Files[i].FullName
		}

		if strings.Contains(bc.Files[i].Name, "Претензия") {
			FileClaim = bc.Files[i].FileID
			FileClaimName = bc.Files[i].FullName
		}

		if strings.Contains(bc.Files[i].Name, "Реестр") {
			FileClaimDetail = bc.Files[i].FileID
			FileClaimDetailName = bc.Files[i].FullName
		}
	}
	result["File_Mail"] = FileMail
	result["File_MailName"] = FileMailName
	result["File_Claim"] = FileClaim
	result["File_ClaimName"] = FileClaimName
	result["File_ClaimDetail"] = FileClaimDetail
	result["File_ClaimDetailName"] = FileClaimDetailName

	// TODO Переделать под нормальные статусы
	if lawSuit.NotifyClaimDone {
		result["Notify_ClaimStatus"] = fmt.Sprintf("%v", "Доставлено успешно")
	} else if lawSuit.NotifyClaimAt.Before(time.Now().AddDate(-10, 1, 1)) {
		result["Notify_ClaimStatus"] = fmt.Sprintf("%v", "Ожидание")
	} else if !contract.IsValidEmail {
		result["Notify_ClaimStatus"] = fmt.Sprintf("%v", "Не доставлено (недоступен канал)")
	} else if contract.Email == "" {
		result["Notify_ClaimStatus"] = fmt.Sprintf("%v", "Не доставлено (отсутствует канал)")
	} else {
		result["Notify_ClaimStatus"] = fmt.Sprintf("%v", contract.Email)
	}
	if lawSuit.NotifyPretrialDone {
		result["Notify_PretrialStatus"] = fmt.Sprintf("%v", "Доставлено успешно")
	} else if lawSuit.NotifyPretrialAt.Before(time.Now().AddDate(-10, 1, 1)) {
		result["Notify_PretrialStatus"] = fmt.Sprintf("%v", "Ожидание")
	} else if !contract.IsValidEmail {
		result["Notify_StatusPretrial"] = fmt.Sprintf("%v", "Не доставлено (недоступен канал)")
	} else if contract.Email == "" {
		result["Notify_PretrialStatus"] = fmt.Sprintf("%v", "Не доставлено (отсутствует канал)")
	} else {
		result["Notify_PretrialStatus"] = fmt.Sprintf("%v", contract.Email)
	}
	result["Notify_ClaimChannel"] = lawSuit.NotifyClaimChannel
	result["Notify_ClaimCode"] = lawSuit.NotifyClaimCode
	result["Notify_ClaimDone"] = lawSuit.NotifyClaimDone
	result["Notify_ClaimMailingCode"] = lawSuit.NotifyClaimMailingCode
	result["Notify_PretrialChannel"] = lawSuit.NotifyPretrialChannel
	result["Notify_PretrialCode"] = lawSuit.NotifyPretrialCode
	result["Notify_PretrialDone"] = lawSuit.NotifyPretrialDone
	result["Notify_PretrialMailingCode"] = lawSuit.NotifyPretrialMailingCode
	if useFormat {
		result["Notify_ClaimAt"] = formatDate(lawSuit.NotifyClaimAt)
		result["Notify_PretrialAt"] = formatDate(lawSuit.NotifyPretrialAt)
	} else {
		result["Notify_ClaimAt"] = lawSuit.NotifyClaimAt
		result["Notify_PretrialAt"] = lawSuit.NotifyPretrialAt
	}

	hashtags := make([]interface{}, 0)
	for i := 0; i < len(bc.Hashtags); i++ {
		tmp := map[string]interface{}{
			"ID":          bc.Hashtags[i].ID,
			"Name":        bc.Hashtags[i].Name,
			"Description": bc.Hashtags[i].Description,
		}
		hashtags = append(hashtags, tmp)
	}
	result["Lawsuit_Hashtags"] = hashtags

	return result, nil
}
