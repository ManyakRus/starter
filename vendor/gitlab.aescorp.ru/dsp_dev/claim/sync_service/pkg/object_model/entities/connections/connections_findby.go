//Файл создан автоматически кодогенератором crud_generator
//Не изменяйте ничего здесь.

package connections

import (
	"gitlab.aescorp.ru/dsp_dev/claim/sync_service/pkg/db/db_constants"
)

// FindBy_BranchID_IsLegal - находит запись по BranchID+IsLegal
func (m *Connection) FindBy_BranchID_IsLegal() error {
	if Crud_manual_Connection == nil {
		return db_constants.ErrorCrudIsNotInit
	}

	err := Crud_Connection.FindBy_BranchID_IsLegal(m)

	return err
}
