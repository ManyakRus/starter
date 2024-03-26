package connections

// Crud_manual_Connection - объект контроллер crud операций
var Crud_manual_Connection ICrud_manual_Connection

// интерфейс CRUD операций сделанных вручную, для использования в DB или GRPC или NRPC
type ICrud_manual_Connection interface {
}

// SetCrudManualInterface - заполняет интерфейс crud: DB, GRPC, NRPC
func (m Connection) SetCrudManualInterface(crud ICrud_manual_Connection) {
	Crud_manual_Connection = crud

	return
}
