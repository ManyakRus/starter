//File generated automatic with crud_generator app
//Do not change anything here.

package connections

import ()

// Connection - модель для таблицы connections: Подключения к БД СТЕКа.
type Connection struct {
	BranchID int64  `json:"branch_id" gorm:"column:branch_id;default:0" db:"branch_id"`          //Филиал (ИД)
	DbName   string `json:"db_name" gorm:"column:db_name;default:\"\"" db:"db_name"`             //Имя таблицы
	DbScheme string `json:"db_scheme" gorm:"column:db_scheme;default:\"\"" db:"db_scheme"`       //Имя схемы
	ID       int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement:true;default:0" db:"id"` //Уникальный технический идентификатор
	IsLegal  bool   `json:"is_legal" gorm:"column:is_legal" db:"is_legal"`                       //Это соединение для юридических лиц
	Login    string `json:"login" gorm:"column:login;default:\"\"" db:"login"`                   //Логин
	Name     string `json:"name" gorm:"column:name;default:\"\"" db:"name"`                      //Наименование
	Password string `json:"password" gorm:"column:password;default:\"\"" db:"password"`          //Пароль
	Port     string `json:"port" gorm:"column:port;default:\"\"" db:"port"`                      //Номер порта
	Server   string `json:"server" gorm:"column:server;default:\"\"" db:"server"`                //Имя сервера, или ip-адрес

}
