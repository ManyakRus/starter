package table_connections

// Table_Connection - модель для таблицы connections: Подключения к БД СТЕКа.
type Table_Connection struct {
	BranchID int64  `json:"branch_id" gorm:"column:branch_id;default:0"`       //Филиал (ИД)
	DbName   string `json:"db_name" gorm:"column:db_name;default:\"\""`        //Имя таблицы
	DbScheme string `json:"db_scheme" gorm:"column:db_scheme;default:\"\""`    //Имя схемы
	ID       int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement:true"` //Уникальный технический идентификатор
	IsLegal  bool   `json:"is_legal" gorm:"column:is_legal"`                   //Это соединение для юридических лиц
	Login    string `json:"login" gorm:"column:login;default:\"\""`            //Логин
	Name     string `json:"name" gorm:"column:name;default:\"\""`              //Наименование
	Password string `json:"password" gorm:"column:password;default:\"\""`      //Пароль
	Port     string `json:"port" gorm:"column:port;default:\"\""`              //Номер порта
	Prefix   string `json:"prefix" gorm:"column:prefix;default:\"\""`          //Префикс для автоматизации
	Server   string `json:"server" gorm:"column:server;default:\"\""`          //Имя сервера, или ip-адрес

}
