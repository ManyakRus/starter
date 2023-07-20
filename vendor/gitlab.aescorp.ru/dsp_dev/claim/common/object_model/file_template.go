package object_model

// FileTemplate Шаблоны файлов (справочник).
type FileTemplate struct {
	CommonStruct
	NameStruct
	FileID string `json:"file_id" gorm:"column:file_id;default:\"\""`
}
