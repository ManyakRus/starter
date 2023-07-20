package object_model

type Connection struct {
	ID       int64  `json:"id"        gorm:"column:id;primaryKey;autoIncrement:true"`
	Name     string `json:"name"      gorm:"column:name;default:\"\""`
	IsLegal  bool   `json:"is_legal"  gorm:"column:is_legal;default:false"`
	BranchId int64  `json:"branch_id" gorm:"column:branch_id;default:0"`
	Server   string `json:"server"    gorm:"column:server;default:\"\""`
	Port     string `json:"port"      gorm:"column:port;default:\"\""`
	DbName   string `json:"db_name"   gorm:"column:db_name;default:\"\""`
	DbScheme string `json:"db_scheme" gorm:"column:db_scheme;default:\"\""`
	Login    string `json:"login"     gorm:"column:login;default:\"\""`
	Password string `json:"password"  gorm:"column:password;default:\"\""`
}
