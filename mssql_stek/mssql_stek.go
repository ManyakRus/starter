package mssql_stek

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ManyakRus/starter/logger"
	"time"

	"github.com/golang-module/carbon/v2"
	"gitlab.aescorp.ru/dsp_dev/claim/sync_service/pkg/object_model/entities/connections"

	//"github.com/ManyakRus/starter/common/pkg/v0/stopapp"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/mssql_gorm"
)

// log - глобальный логгер
var log = logger.GetLog()

// Loc - Локация (страна)
var Loc = time.Local

// FindDate_ClosedMonth - возвращает дату последнего закрытого периода в СТЕК
func FindDate_ClosedMonth(c connections.Connection) (time.Time, error) {
	var err error
	var Otvet time.Time

	Connection_id := c.ID

	ctxMain := contextmain.GetContext()
	ctx, ctxCancelFunc := context.WithTimeout(ctxMain, time.Second*600)
	defer ctxCancelFunc()

	db := mssql_gorm.GetConnection(Connection_id)
	db.WithContext(ctx)

	text_sql := `
		SELECT 
		    [Месяц] AT TIME ZONE 'Russian Standard Time' as Date
		FROM 
			[stack].[Закрытые месяцы]
		WHERE 
		    Задача = @Kod
	`

	kod := 252
	if c.IsLegal == true {
		kod = 11058
	}

	// Execute query
	//DateNow1 := time.Now()
	param1 := sql.Named("Kod", kod)
	DB, err := db.DB()
	if err != nil {
		sError := fmt.Sprint("db.DB() error: ", err)
		//log.Error(sError)
		//stopapp.StopAppAndWait()
		log.Panicln(sError)
		return Otvet, errors.New(sError)
	}

	rows, err := DB.QueryContext(ctx, text_sql, param1)
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error("rows.Close() error: ", err)
		}
	}()

	if err != nil {
		sError := fmt.Sprint("MSSQL query error: ", err, " connection_id: ", Connection_id)
		//log.Error(sError)
		//stopapp.StopAppAndWait()
		log.Panicln(sError)
		return Otvet, errors.New(sError)
	}

	//даты
	for rows.Next() {
		err := rows.Scan(&Otvet)
		if err != nil {
			log.Warn("Scan() error: ", err)
			continue
		}
	}

	Otvet = Otvet.Local()

	return Otvet, err
}

// FindDateFromTo - находит даты загрузки остатков и документов
func FindDateFromTo(Connection connections.Connection) (date1_balances, date2_balances, date1_doc, date2_doc time.Time, err error) {

	date_closed, err := FindDate_ClosedMonth(Connection)

	if err != nil {
		return date1_balances, date2_balances, date1_doc, date2_doc, err
	}

	if date_closed.Year() < 2022 {
		Text := "Error: Closed month < 2022"
		log.Panic(Text)
		err = errors.New(Text)
		return date1_balances, date2_balances, date1_doc, date2_doc, err
	}

	date1_balances, date2_balances, date1_doc, date2_doc = FindDates_from_DateClosed(date_closed)

	return date1_balances, date2_balances, date1_doc, date2_doc, nil
}

// FindDates_from_DateClosed - находит даты загрузки с учётом даты закрытого периода
func FindDates_from_DateClosed(date_closed time.Time) (date1_balances, date2_balances, date1_doc, date2_doc time.Time) {
	// дата сегодня
	carbon.SetLocation(Loc)
	DateNow1 := time.Now()
	DateNow2 := carbon.Time2Carbon(DateNow1).EndOfDay().Carbon2Time()

	// 2 месяца разделим по 1 месяцу
	date_closed_next1 := carbon.Time2Carbon(date_closed).AddMonth().StartOfMonth().Carbon2Time()

	date1_balances = date_closed_next1
	date2_balances = carbon.Time2Carbon(date1_balances).EndOfMonth().Carbon2Time()
	if date2_balances.After(DateNow2) {
		date2_balances = DateNow2
	}
	date1_doc = date_closed_next1
	date2_doc = DateNow2

	return date1_balances, date2_balances, date1_doc, date2_doc

}
