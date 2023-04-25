// модуль для работы с базой данных

package db

import (
	//"fmt"

	"github.com/manyakrus/starter/contextmain"
	"github.com/manyakrus/starter/logger"

	_ "github.com/denisenkom/go-mssqldb"
	//"github.com/denisenkom/go-mssqldb"
	mssql "github.com/manyakrus/starter/mssql_connect"
)

// // log - глобальный логгер
var log = logger.GetLog()

type DBRecord struct {
	PersonalAccount string
	Summa           float64
}

func FindDebtorsList() ([]DBRecord, error) {
	Otvet := make([]DBRecord, 0)

	tsql := `
DECLARE @Mounth datetime = '20220701';

IF OBJECT_ID('tempdb..#ls0') IS NOT NULL DROP TABLE #ls0;
CREATE TABLE #ls0 ( [Счет-Параметры] int, [ДатКнц] datetime, [Значение] float, [Сумма] float);
create clustered index ux_ls0 on tempdb..#ls0 ([Счет-Параметры], [ДатКнц], [Значение]);

IF OBJECT_ID('tempdb..#ls2') IS NOT NULL DROP TABLE #ls2;
CREATE TABLE #ls2 ( [Счет-Параметры] int, [ДатКнц] datetime);
create clustered index ux_ls2 on tempdb..#ls2 ([Счет-Параметры], [ДатКнц]);

IF OBJECT_ID('tempdb..#ls') IS NOT NULL DROP TABLE #ls;
CREATE TABLE #ls ( id int);

IF OBJECT_ID('tempdb..#ns') IS NOT NULL DROP TABLE #ns;
CREATE TABLE #ns ( [Счет] int, [Сумма] float);


insert into #ns ([Счет], [Сумма])
select ns.Счет ,sum(Сумма) as Сумма
from stack.НСальдо ns 
where ns.[Номер услуги]/100 in (1,4,101) -- Это ЭЭ
and ns.[Месяц расчета] = dateadd(month,-1,@Mounth)
group by ns.Счет
HAVING sum(Сумма) >10000
;



insert into #ls0 ([Счет-Параметры], [ДатКнц], [Значение], [Сумма])
select #ns.[Счет], 
	sv.ДатКнц, 
	isnull(sv.Значение,0) as [Значение],
	#ns.[Сумма] as [Сумма]
from #ns

-- left JOIN 
--	stack.[Виды параметров] vp 
--ON 
--	vp.Название = 'СОСТОЯНИЕ' -- 113.871

 left join 
	stack.Свойства sv 
on 
	#ns.[Счет] = sv.[Счет-Параметры] 
--	and sv.[Виды-Параметры] = vp.row_id
	and	@Mounth BETWEEN sv.ДатНач AND sv.ДатКнц -- 5.589.755
	and sv.[Виды-Параметры] = 76 --'СОСТОЯНИЕ'

;



insert into #ls2 ([Счет-Параметры], [ДатКнц])
select sv.[Счет-Параметры],
	max(sv.ДатКнц) as ДатКнц  
from #ls0 as sv
group by
	sv.[Счет-Параметры]
;


select top 10 
	#ls0.[Счет-Параметры] as [Счет], 
	#ls0.[Сумма] as [Сальдо]
--	#ls0.ДатКнц as [ДатаСостояния], 
--	#ls0.Значение as [Состояние]
from #ls0 

join
	#ls2
on
	#ls2.[Счет-Параметры] = #ls0.[Счет-Параметры]
	and isnull(#ls2.[ДатКнц], 0) = isnull(#ls0.[ДатКнц], 0)
	
JOIN 
	stack.[Лицевые счета] l
on
	l.ROW_ID = #ls0.[Счет-Параметры]
	
where 1=1
	and l.тип = 5
	and #ls0.Значение <>2

	
order by 
	#ls0.[Сумма] DESC
 
`

	//tsql := `SELECT "" as aa, 0 as xx`
	//tsql := `select l.Номер, l.row_id
	//	from
	//	stack.[Лицевые счета] as l`

	ctx := contextmain.GetContext()

	// Execute query
	db := mssql.Conn
	//db.Exec("select 1")
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error("rows.Close() error: ", err)
		}
	}()

	//var count int

	// Iterate through the result set.
	for rows.Next() {
		var PersonalAccount string
		var Summa float64

		// Get values from row.
		err := rows.Scan(&PersonalAccount, &Summa)
		if err != nil {
			log.Warn("Scan() error: ", err)
			continue
		}

		Otvet1 := DBRecord{}
		Otvet1.PersonalAccount = PersonalAccount
		Otvet1.Summa = Summa
		Otvet = append(Otvet, Otvet1)
	}

	return Otvet, nil
}
