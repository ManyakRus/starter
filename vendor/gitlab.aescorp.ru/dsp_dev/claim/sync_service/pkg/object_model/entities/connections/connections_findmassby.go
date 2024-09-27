//Файл создан автоматически кодогенератором crud_generator
//Не изменяйте ничего здесь.

package connections

import (
	"gitlab.aescorp.ru/dsp_dev/claim/sync_service/pkg/db/db_constants"
)

// FindMassBy_BranchID - находит запись по BranchID
func (m *Connection) FindMassBy_BranchID() ([]Connection, error) {
	Otvet := make([]Connection, 0)
	if Crud_manual_Connection == nil {
		return Otvet, db_constants.ErrorCrudIsNotInit
	}

	Otvet, err := Crud_Connection.FindMassBy_BranchID(m)

	return Otvet, err
}
